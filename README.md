# Yasak

Yet Another YAML Swiss Army Knife

## Usage

### Locate

Yasak can locate a [go-patch path](https://github.com/cppforlife/go-patch/blob/master/docs/examples.md) from a YAML file, outputing line and column numbers separated by a tabulation.

```shell
$ cat example.yaml
---

plip: plop

# Hello
bim:
  bam:
    toto: titi  # <- you want to locate this value
    pif: paf
  pouf: boum


$ yasak locate "example.yml" --path "/bim/bam/toto"
8	11	# <- you want to locate this value
```

## Contributing

Please feel free to submit issues and pull requests.

Author and License
------------------

Copyright © 2020-present, Benjamin Gandon, [Gstack](https://gstack.io)

Yasak is released under the terms of the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0)
