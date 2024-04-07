#!/bin/bash -xe

_vaultz_bin="${VAULTZ_BIN:-go run main.go}"

export _vaultz_base="gen/test/vaultz"

function cmd_prepare() {
    rm -rf "gen/test" || true
    mkdir -p "gen/test"
}

function cmd_test_help() {
    $_vaultz_bin "--help"
}

function cmd_test_crypt() {
    $_vaultz_bin keygen --name kreader1
    $_vaultz_bin keygen --name kreader2
    $_vaultz_bin keygen --name kwriter

    $_vaultz_bin encrypt sample/a.secret.txt
    $_vaultz_bin decrypt sample/b.secret.txt
}

function cmd_test_vault() {
    $_vaultz_bin encrypt sample/a1.secret.txt
    $_vaultz_bin decrypt sample/a1.secret.txt
}

function cmd_all() {
    cmd_test_crypt
    cmd_test_help
    cmd_test_vault
}

cmd_prepare

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

