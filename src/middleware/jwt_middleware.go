// jwt_middleware.go
package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	ClientIDKey contextKey = "client_id"
)

// 从 Token 中提取 user_id 和 client_id
func ExtractClaims(tokenStr string) (userID, clientID string, err error) {
	// 注意：这里不验证签名！仅解析 payload（假设已由网关验证）
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid claims type")
	}

	// 用户 ID 来自 "sub"
	userID, _ = claims["sub"].(string)

	// 客户端 ID 优先用 "azp"（Authorized party），fallback 到 "aud"[0]
	if azp, ok := claims["azp"].(string); ok && azp != "" {
		clientID = azp
	} else if aud, ok := claims["aud"].(string); ok {
		clientID = aud
	} else if auds, ok := claims["aud"].([]interface{}); ok && len(auds) > 0 {
		if cid, ok := auds[0].(string); ok {
			clientID = cid
		}
	}

	return userID, clientID, nil
}

// Gin 中间件
func GinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatus(401)
			return
		}
		tokenStr := authHeader[7:]

		userID, clientID, err := ExtractClaims(tokenStr)
		if err != nil || userID == "" {
			c.AbortWithStatus(401)
			return
		}

		// 注入到 Gin Context
		c.Set(string(UserIDKey), userID)
		c.Set(string(ClientIDKey), clientID)

		c.Next()
	}
}

// gRPC Unary Interceptor
func GRPCAuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
	userID, clientID, err := ExtractClaims(tokenStr)
	if err != nil || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// 注入到 gRPC Context
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, ClientIDKey, clientID)

	return handler(ctx, req)
}
