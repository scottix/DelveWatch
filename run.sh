#!/bin/bash

go get ./...

./DelveWatch -delve=./test/test1.go -verbose
