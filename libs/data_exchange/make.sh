#!/bin/sh

rm *.pb.go
protoc --go_out=. --go_opt=module=github.com/flinkcoin/flink/libs/data *.proto
