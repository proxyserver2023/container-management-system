build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix 'static' -o consignment-cli cmd/consignment-cli/main.go
	docker build -t consignment-cli .

run:
	docker run -e MICRO_REGISTRY=mdns consignment-cli
