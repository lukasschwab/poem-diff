# poem-diff

A demo tool for comparing translations of poems. See ["Four Translations of *Todesfuge*"](https://lukasschwab.me/blog/gen/celan-translations.html#translations).

> [!WARNING]
> There's some LLM slop in `diff.go`.

## Installation

Run `go install .` in this directory to install the `poem-diff` binary.

## Usage

```console
$ poem-diff --left felstiner.txt --right weimar.txt --format ansi
```

+ Files `--left` and `--right` must have the same number of lines.
+ `--format` must be either `ansi` or `html`.
    + `ansi` mode has pretty colors.
    + `html` mode renders a table; see e.g. [felstiner-hamburger.diff.html](https://lukasschwab.me/blog/img/celan-translations/diffs/felstiner-hamburger.diff.html).

![termshot](https://github.com/user-attachments/assets/06bc4aa2-c307-4555-94cd-f4cacc62126a)
