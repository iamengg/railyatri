build:
	go build -o client/bin/  client/cmd/main.go 
	go build -o server/bin/ server/cmd/main.go

clear:
	rm -rf client/bin/*	
	rm -rf server/bin/*

