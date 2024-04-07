#!/bin/bash -xe

_vaultz_bin="${VAULTZ_BIN:-go run main.go}"

function cmd_help() {
    $_vaultz_bin "--help"
}

function cmd_basics() {
    true
}

function cmd_all() {
    cmd_basics
    cmd_help
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

