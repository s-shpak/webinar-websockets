.PHONY: generate-certs
generate-certs:
	openssl ecparam -genkey -name secp384r1 -out server.key
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

.PHONY: build
build: build-client build-server

.PHONY: build-client
build-client:
	go build -o client ./cmd/client

.PHONY: build-server
build-server:
	go build -o server ./cmd/server