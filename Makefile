SWAGGER = ezcards-api-v2-openapi.yaml
MOUNT = /local
LANG = go
GENERATED = ./ez-client-go

up:
	docker compose up -d

generate:
	@mkdir -p "$(GENERATED)" 
	docker compose --profile codegen run --rm swagger-codegen generate \
		-i $(MOUNT)/$(SWAGGER) \
		-l $(LANG) \
		-o $(MOUNT)/$(GENERATED)

generate-go:
	$(MAKE) generate LANG=go

langs:
	docker compose --profile codegen run --rm swagger-codegen langs
