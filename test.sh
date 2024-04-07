#!/bin/bash -xe

_vaultz_bin="${VAULTZ_BIN:-go run main.go}"

export VAULTZ_BASE="target/test/vaultz"

function cmd_prepare() {
    rm -rf "target/test" || true
    mkdir -p "target/test/sample"

    echo aaavalue > "target/test/sample/a.secret.txt"
    echo bbbvalue > "target/test/sample/b.secret.txt"
}

function cmd_test_help() {
    $_vaultz_bin "--help"
}

function cmd_test_crypt() {
    $_vaultz_bin keygen --name kreader1
    $_vaultz_bin keygen --name kreader2
    $_vaultz_bin keygen --name kwriter

    $_vaultz_bin encrypt --file target/test/sample/a.secret.txt
    $_vaultz_bin decrypt --file target/test/sample/b.secret.txt
}

function cmd_all() {
    cmd_test_crypt
    cmd_test_help
}

cmd_prepare

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

