.PHONY: gen_geovision

# --- Path Setup ---
PROTO_DIR   := ./protos/

PROTOC_INCLUDES := -I$(PROTO_DIR)

# --- Build Targets ---
gen_geovision: DESTDIR := ./gen/go
gen_geovision: PROTO_OUT := \
    --go_out=$(DESTDIR) --go_opt=paths=source_relative \
    --go-grpc_out=$(DESTDIR) --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=$(DESTDIR) --grpc-gateway_opt=paths=source_relative
gen_geovision: PROTO_FILES := $(wildcard $(PROTO_DIR)/geovision/*.proto)
gen_geovision: PROTO_FILES += $(wildcard $(PROTO_DIR)/model/*.proto)
gen_geovision: proto_gen

gen_geovision_openapi: DESTDIR := ./gen/http
gen_geovision_openapi: PROTO_OUT := --openapiv2_out=$(DESTDIR) --openapiv2_opt=allow_merge=true
gen_geovision_openapi: PROTO_FILES := $(wildcard $(PROTO_DIR)/geovision/*.proto)
gen_geovision_openapi: PROTO_FILES += $(wildcard $(PROTO_DIR)/model/*.proto)
gen_geovision_openapi: proto_gen

proto_gen:
	@echo "--- Cleaning: $(DESTDIR) ---"
	rm -rf $(DESTDIR)
	@echo "--- Creating: $(DESTDIR) ---"
	mkdir -p "$(DESTDIR)"

	@echo "--- Generating protobufs ---"
	protoc $(PROTOC_INCLUDES) ${PROTO_OUT} $(PROTO_FILES)
