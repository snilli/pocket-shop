# Generate ent code from schemas
ent-gen:
	go generate ./ent

# Create new ent schema (usage: make ent-new name=Book)
ent-new:
	@if [ -z "$(name)" ]; then echo "Usage: make ent-new name=SchemaName"; exit 1; fi
	go run -mod=mod entgo.io/ent/cmd/ent new $(name)

clean:
	rm -rf ent/hook ent/migrate ent/predicate ent/runtime ent/enttest ent/user ent/product ent/order
	cd ent && rm -f client.go ent.go gql_*.go mutation.go order*.go product*.go runtime.go tx.go user*.go ent.graphql
	rm -f ent/graphql/generated.go

fmt:
	go fmt ./...

tidy:
	go mod tidy


# Generate new migration (usage: make migrate-diff name=add_feature)
migrate-diff:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-diff name=migration_name"; exit 1; fi
	go generate ./ent
	@export $$(cat ../go-long/.env | xargs) && atlas migrate diff $(name) --env postgres

# Apply migrations to database (loads .env from go-long)
migrate-apply:
	@export $$(cat ../go-long/.env | xargs) && atlas migrate apply --env postgres --allow-dirty

# Check migration status
migrate-status:
	@export $$(cat ../go-long/.env | xargs) && atlas migrate status --env postgres

# Rehash migration files (after manual edit)
migrate-hash:
	@export $$(cat ../go-long/.env | xargs) && atlas migrate hash --env postgres

migrate: migrate-diff migrate-apply
