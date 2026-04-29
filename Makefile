.PHONY: codegen codegen-clean proto proto-lint proto-clean oapi oapi-clean

codegen: proto oapi

codegen-clean: proto-clean oapi-clean

proto:
	buf generate

proto-lint:
	buf lint

proto-clean:
	rm -rf internal/generated/auth

oapi:
	mkdir -p internal/generated/api
	oapi-codegen -config api/cfg/server.yaml api/schema.yaml

oapi-clean:
	rm -rf internal/generated/api