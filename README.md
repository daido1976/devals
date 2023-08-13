# devals

`devals` is a tool built on top of [helmfile/vals](https://github.com/helmfile/vals).

While `helmfile/vals` only accepts JSON/YAML format as input, `devals` supports the dotenv format. The name "devals" is short for "dotenv vals".

## Installation

### Using Homebrew

Install `devals` via Homebrew with the following command:

```sh
brew install daido1976/tap/devals
```

### From Binary

You can also download the binary from the [releases page](https://github.com/daido1976/devals/releases).

## [WIP] Usage

```sh
$ devals -h
Usage of devals:
  -i string
        Input .env file (required)
  -keep-comments
        Keep comments and empty lines in the output
  -o string
        Output file. If not specified, writes to stdout.
```
