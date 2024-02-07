#!/bin/bash

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration11$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -database-dsn='postgres://mc_user:Dthcbz@localhost:5432/mc_db?sslmode=disable' \
    -server-port=$SERVER_PORT \
    -source-path=.