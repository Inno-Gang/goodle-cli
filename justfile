#!/usr/bin/env just --justfile

go-mod := `go list`
flags := '-trimpath -ldflags="-s -w"'

# Run goodle
run:
    go run .

# Build `goodle-cli`
build:
    go build {{flags}}

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