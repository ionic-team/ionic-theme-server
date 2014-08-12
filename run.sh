#!/bin/bash
CGO_CFLAGS="-I$PWD/lib/libsass" CGO_LDFLAGS="-L$PWD/lib/libsass" go run main.go
