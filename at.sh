#!/bin/bash
metricstest -test.v -test.run=^TestIteration1$ -binary-path=cmd/server/server

metricstest -test.v -test.run=^TestIteration2[AB]*$ \
    -source-path=. \
    -agent-binary-path=cmd/agent/agent

metricstest -test.v -test.run=^TestIteration3[AB]*$ \
    -source-path=. \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration4$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -server-port=$SERVER_PORT \
    -source-path=.

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration5$ \
  -agent-binary-path=cmd/agent/agent \
  -binary-path=cmd/server/server \
  -server-port=$SERVER_PORT \
  -source-path=.

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration6$ \
  -agent-binary-path=cmd/agent/agent \
  -binary-path=cmd/server/server \
  -server-port=$SERVER_PORT \
  -source-path=.

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration7$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -server-port=$SERVER_PORT \
    -source-path=.

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
    metricstest -test.v -test.run=^TestIteration8$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -server-port=$SERVER_PORT \
    -source-path=.

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration9$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -file-storage-path=$TEMP_FILE \
    -server-port=$SERVER_PORT \
    -source-path=.