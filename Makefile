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
i4:
	metricstest -test.v -test.run=^TestIteration4$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=8081 \
            -source-path=.
i5:
	metricstest -test.v -test.run=^TestIteration5$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=8081 \
            -source-path=.

sprint1: build i1 i2 i3 i4 i5