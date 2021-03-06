#!/usr/bin/env bash

function main() {
    local file=${1:-"example.yml"}
    local path=${2:-"/bim/bam/toto"}
    local new_value=$3

    local coords line column line_comment
    coords=$(yasak locate "${file}" --path "${path}")
    line=$(cut -f 1 <<< "${coords}")
    column=$(cut -f 2 <<< "${coords}")
    line_comment=$(cut -f 3- <<< "${coords}")

    local replacement repeat
    replacement="${new_value}"
    if [[ -n ${line_comment} ]]; then
        replacement+=" ${line_comment}"
    fi
    repeat=$((${column} - 1))
    sed -e "${line}s/^\\(.\\{${repeat},${repeat}\\}\\).*$/\\1${replacement}/" "${file}"
}

main "$@"
