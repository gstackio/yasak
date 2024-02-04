Yasak
=====

Yet Another YAML Swiss Army Knife.

I've written this tool in order to help me update component version in my YAML
manifests. Indeed, I treat my YAML manifest files with a lot of atention and
respect, because these potentially are elements of proper communication with
teammates. So, they are tailored pieces of writing, where each space is there
for a reason. As such, I wanted my automations update tools not to mess up
with the human-friendly layout of my YAML documents.

More specifically:

- Keep the ordering of YAML keys in dictionaries.
- Keep the comments, especially at the end of modified lines.
- Keep the choosen spacing between sections, possibly made with multiple new
  lines between main sections.

All in all, I wanted to keep my YAML files made for humans, and don't let a
robot rewrite my manifests for any robots. Because I'm not a robot, I'm a
human being, so are my teammates and this absolutely _does_ matter.


Usage
-----

At this stage, Yasak implements one single verb `locate`, accepting one single
argument introduced by `--path`, referring to the `--path` argument accepted
by [`bosh interpolate`][bosh_int].

[bosh_int]: https://bosh.io/docs/cli-v2/#interpolate

### Locate

Yasak can basically locate the start of a YAML node in a serialized YAML text
file. The YAML node is specified as a [go-patch path][go_patch_path]. Yasak
output includes, the line at which the value appears, the column number at
which the value starts, and any end-of-line comment at that can thus stay
where it was.

These 3 fields are output on one single line, and separated by horizontal
tabulation characters, i.e. `HT` in [ASCII RFC 20][rfc_20_s2].

[go_patch_path]: https://github.com/cppforlife/go-patch/blob/master/docs/examples.md
[rfc_20_s2]: https://datatracker.ietf.org/doc/html/rfc20#section-2

```shell
$ cat -n "example.yml"
     1  ---
     2
     3  plip: plop
     4
     5  # Hello
     6  bim:
     7    bam:
     8      toto: titi  # <- you should update this value
     9      pif: paf
    10    pouf: boum

$ yasak locate "example.yml" --path "/bim/bam/toto"
8	11	# <- you should update this value
```

The horizontal tabulation `HT` separator is specifically choosen in order to
match the default separator of the `cut` utility, as showcased in the Bash
snippet below.

```bash
coords=$(yasak locate "${file}" --path "${path}")
line=$(cut -f 1 <<< "${coords}")
column=$(cut -f 2 <<< "${coords}")
line_comment=$(cut -f 3- <<< "${coords}")
```


Contributing
------------

Please feel free to submit issues and pull requests.


Author and License
------------------

Copyright © 2020-present, Benjamin Gandon, [Gstack][gstack_io]

Yasak is released under the terms of the [Apache 2.0 license][apache2_license].

[gstack_io]: https://gstack.io
[apache2_license]: http://www.apache.org/licenses/LICENSE-2.0
