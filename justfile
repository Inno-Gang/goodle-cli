#!/usr/bin/env just --justfile

go-mod := `go list`
flags := '-ldflags="-s -w"'

# Build `goodle-cli`
build:
    go build {{flags}}

# Run goodle
run:
    go run .

# Install `goodle-cli` to the GOBIN
install:
    go install {{flags}}

# Update the dependencies
update:
    go get -u
    go mod tidy -v

# Rename go.mod name
rename new-go-mod:
    find . -type f -not -path './.git/*' -exec sed -i '' -e "s|{{go-mod}}|{{new-go-mod}}|g" {} \;