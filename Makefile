build: agent server

agent:
	cd cmd/agent/ && go build -o agent *.go

server:
	cd cmd/server && go build -o server *.go

agent_ext:
	go build -o agent -ldflags " -X 'main.buildDate=$(shell date +'%Y.%m.%d %H:%M:%S')' -X main.buildVersion=$(VER) -X main.buildCommit=$(COMMIT)" cmd/agent/main.go

server_ext:
	go build -o server -ldflags " -X 'main.buildDate=$(shell date +'%Y.%m.%d %H:%M:%S')' -X main.buildVersion=$(VER) -X main.buildCommit=$(COMMIT)" cmd/server/main.go

build_multichecker:
	go build cmd/staticlint/multichecker.go

checks: vet check

check:
	./multichecker ./...

vet:
	go vet ./...

t:
	go test ./...
b:
	go test -v --bench . --benchmem
cov:
	go clean -testcache && \
	go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./... && go tool cover -func profile.cov && \
	cat profile.cov.tmp | grep -v "mock_mem_store.go" > profile.cov && \
	go tool cover -func profile.cov

pheap:
	go tool pprof -http=":9090" -seconds=60 http://localhost:8080/debug/pprof/heap
pheapsave:
	curl -s -v http://localhost:8080/debug/pprof/heap?seconds=60 > tmp.out 
	go tool pprof -http=":9090" -seconds=60 tmp.out
doclocal:
	godoc -http=:8081

swaggen:
	swag init --output ./swagger/
