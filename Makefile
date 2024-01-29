build: agent server

agent:
	cd cmd/agent/ && go build -o agent *.go

server:
	cd cmd/server && go build -o server *.go

t:
	go test ./...
