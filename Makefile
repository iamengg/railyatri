proto:
	protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative bookingStub/*.proto

buildClient:
	go build -o client/bin/  client/cmd/main.go 

buildServer:	
	go build -o server/bin/ server/cmd/main.go

build:	
	make buildClient
	make buildServer

clear:
	rm -rf client/bin/*	
	rm -rf server/bin/*

