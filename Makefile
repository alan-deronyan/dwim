# DWIM makefile

build: # Build DWIM
	@echo "Building DWIM..."
	@go build -o bin/schema_gen cmd/schema_gen.go ingest/main.go

# Generate structs from schemas, depends on build target above
generate: build
	./bin/schema_gen

fluree-docker:
	@echo "Initializing Fluree Docker..."
	@docker volume create fluree_server
	@docker create --name "fluree_server" -p 58090:8090 -v fluree_server:/opt/fluree-server/data fluree/server

fluree-start:
	@echo "Starting Fluree..."
	@docker start fluree_server

fluree-stop:
	@echo "Stopping Fluree..."
	@docker stop fluree_server
