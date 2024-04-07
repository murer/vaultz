#!/bin/bash -xe

_vaultz_bin="${VAULTZ_BIN:-go run main.go}"

function cmd_help() {
    $_vaultz_bin "--help"
}

function cmd_crypt() {
    $_vaultz_bin encrypt sample/a1.secret.txt
    $_vaultz_bin decrypt sample/a1.secret.txt
}

function cmd_all() {
    cmd_crypt
    cmd_help
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

