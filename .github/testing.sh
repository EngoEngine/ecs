#!/usr/bin/env bash

go test -coverprofile=coverage.out -bench=. .
go tool cover -func=coverage.out -o=coverage.out

out = $(grep -Po '(?<=\(statements\)\s+)\d+.\d+%' coverage.out)

echo "$out"
