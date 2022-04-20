#!/bin/bash
echo "operation count"
read count
cd ./internal/client
go test -bench=. -benchtime="${count}x"
