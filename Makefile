build: agent server
agent:
	cd cmd/agent/ && go build -o agent *.go
server:
	cd cmd/server && go build -o server *.go
t:
	go test ./...
i1:
	metricstest -test.v -test.run=^TestIteration1$$ -binary-path=./cmd/server/server
i2:
	metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent
i3:
	metricstest -test.v -test.run=^TestIteration3[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server
