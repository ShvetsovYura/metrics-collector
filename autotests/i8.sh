#!/bin/bash

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
    metricstest -test.v -test.run=^TestIteration8$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -server-port=$SERVER_PORT \
    -source-path=.