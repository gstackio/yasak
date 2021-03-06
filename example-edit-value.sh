#!/usr/bin/env bash

function main() {
	local file=${1:-"example.yml"}
	local path=${2:-"/bim/bam/toto"}
	local new_value=$3

	local coords line column line_comment
	coords=$(./yasak locate example.yml --path "${path}")
	line=$(cut -f 1 <<< "${coords}")
	column=$(cut -f 2 <<< "${coords}")
	line_comment=$(cut -f 3- <<< "${coords}")

	local replacement
	replacement="${new_value}"
	if [[ -n ${line_comment} ]]; then
		replacement+=" ${line_comment}"
	fi
	sed -e "${line}s/^\\(.\\{$((${column} - 1))\\}\\).*$/\\1${replacement}/" "${file}"
}

main "$@"
