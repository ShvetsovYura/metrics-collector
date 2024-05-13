build: agent server

agent:
	cd cmd/agent/ && go build -o agent *.go

server:
	cd cmd/server && go build -o server *.go

t:
	go test ./...

cov:
	go clean -testcache
	go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./... && go tool cover -func profile.cov
	cat profile.cov.tmp | grep -v "mock_mem_store.go" > profile.cov
	go tool cover -func profile.cov

