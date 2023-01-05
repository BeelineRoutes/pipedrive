
# go params
GOCMD=go

# normal entry points
	
update:
	clear 
	@echo "updating dependencies..."
	@go get -u -t ./...
	@go mod tidy 

build:
	clear 
	@echo "building..."
	@$(GOCMD) build .
	
test:
	clear
	@echo "testing Pipedrive..."
	@$(GOCMD) test ./...

