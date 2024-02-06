#!/bin/bash

SERVER_PORT=$(random unused-port)
    ADDRESS="localhost:${SERVER_PORT}"
    TEMP_FILE=$(random tempfile)
    metricstest -test.v -test.run=^TestIteration10[AB]$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
    -server-port=$SERVER_PORT \
    -source-path=.
