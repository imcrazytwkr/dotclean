[private]
default:
    @just --list --unsorted

build:
    make build

man:
    make man

test: test-go

[parallel]
format: format-go format-flake format-just

[group("format")]
format-go:
    go fmt ./...
    go mod tidy

[group("format")]
format-flake:
    nix fmt flake.nix

[group("format")]
format-just:
    just --fmt --unstable

[group("test")]
test-go: format
    go test ./... | sed '/^[[:space:]]*?/d'

[group("test")]
test-e2e:
    ./e2e/run.sh

alias test-logged := test-go-logged

[group("ci")]
test-go-logged: format
    go test ./... -v 2>&1 | tee test.log

[group("ci")]
test-e2e-logged:
    ./e2e/run.sh 2>&1 | tee test-e2e.log

vendor-hash:
    go mod vendor
    nix hash path vendor
