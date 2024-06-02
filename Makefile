build: agent server

agent:
	cd cmd/agent/ && go build -o agent *.go

server:
	cd cmd/server && go build -o server *.go

t:
	go test ./...
b:
	go test -v --bench . --benchmem
cov:
	go clean -testcache
	go test -v -coverpkg=./...m=all -coverprofile=profile.cov.tmp ./... && go tool cover -func profile.cov
	cat profile.cov.tmp | grep -v "mock_mem_store.go" > profile.cov
	go tool cover -func profile.cov

pheap:http://localhost:8081/pkg/?m=all
	go tool pprof -http=":9090" -seconds=60 http://localhost:8080/debug/pprof/heap
pheapsave:
	curl -s -v http://localhost:8080/debug/pprof/heap?seconds=60 > tmp.out 
	go tool pprof -http=":9090" -seconds=60 tmp.out
doclocal:
	godoc -http=:8081
.PHONY
swaggen:
	swag init --output ./swagger/
