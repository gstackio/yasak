#!/usr/bin/env bash

./yasak locate "example.yml" --path "/bim/bam/toto"

./yasak locate "example.yml" --path "/bim/bam/toti"

./yasak locate "example.yml" --path "/bim/bam/toti?"

./yasak locate "example.yml" --path "/bim/boum/toto"

./yasak locate "example.yml" --path "/bim/boum?/toto"

bash example-edit-value.sh "example.yml" "/bim/bam/toto" "plop"
