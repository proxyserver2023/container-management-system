build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/alamin-mahamud/container-management-system/consignment-service \
		proto/consignment.proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -installsuffix 'static' -o consignment-service cmd/consignment-service/main.go
	docker build -t consignment-service .

run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns consignment-service
