#!/bin/bash -xe

function cmd_basics() {
    true
}

function cmd_all() {
    cmd_basics
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"

