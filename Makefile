# DWIM makefile

build: # Build DWIM
	@echo "Building DWIM..."
	@go build -o bin/schema_gen cmd/schema_gen.go ingest/main.go

# Generate structs from schemas, depends on build target above
generate: build
	@alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
	./bin/schema_gen

deps:
	@echo "Installing dependencies..."
	@docker pull quay.io/goswagger/swagger

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
