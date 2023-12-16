#!/usr/bin/env bash

go test -coverprofile=coverage.out -bench=. .
go tool cover -func=coverage.out -o=coverage.out

outputFromCoverage=$(grep -Po '(?<=\(statements\)\s\s)\d+.\d+%' coverage.out)

echo "$outputFromCoverage"
