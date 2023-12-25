agent:
	cd $$HOME/metrics-collector/cmd/agent/ && go build -o agent *.go
server:
	cd $$HOME/metrics-collector/cmd/server && go build -o server *.go
run:
	sudo metricstest -test.v -test.run=^TestIteration1$ -binary-path=$$HOME/metrics-collector/cmd/server/server
