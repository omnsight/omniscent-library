// jwt_middleware.go
package middleware

import (
	"context"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "user_id"
	UserRolesKey ContextKey = "user_roles"
)

// IdentityInterceptor parses claims without verifying signature (Gateway trusted)
func GrpcGatewayIdentityInterceptor(clientID string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 1. Extract Token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing auth header")
		}
		tokenString := strings.TrimPrefix(values[0], "Bearer ")

		// 2. Parse Claims (Unverified because Gateway already verified it)
		parser := jwt.NewParser()
		token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to parse token: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Internal, "invalid claims structure")
		}

		// 3. Extract User ID
		userID, _ := claims["sub"].(string)

		// 4. Extract Roles for THIS specific Client
		var roles []string
		if resAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			if clientMap, ok := resAccess[clientID].(map[string]interface{}); ok {
				if roleList, ok := clientMap["roles"].([]interface{}); ok {
					for _, r := range roleList {
						if rStr, ok := r.(string); ok {
							roles = append(roles, rStr)
						}
					}
				}
			}
		}

		// 5. Inject into Context
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRolesKey, roles)

		return handler(ctx, req)
	}
}

func AuthMiddleware(clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 1. Parse Unverified (Trusting the Gateway has already verified signature)
		parser := jwt.NewParser()
		token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid token claims"})
			return
		}

		// 2. Extract User ID (Subject)
		userID, _ := claims["sub"].(string)

		// 3. Extract Client Roles for THIS specific client
		var roles []string
		if resAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			if clientMap, ok := resAccess[clientID].(map[string]interface{}); ok {
				if roleList, ok := clientMap["roles"].([]interface{}); ok {
					for _, r := range roleList {
						if rStr, ok := r.(string); ok {
							roles = append(roles, rStr)
						}
					}
				}
			}
		}

		// 4. Inject into Context for downstream handlers
		c.Set("userID", userID)
		c.Set("userRoles", roles)

		c.Next()
	}
}
