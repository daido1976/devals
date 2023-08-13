# devals

`devals` is a CLI tool that safely manages environment variables in the development environment.

`devals` is built on [helmfile/vals](https://github.com/helmfile/vals).

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
        Input dotenv format file (required)
  -keep-comments
        Keep comments and empty lines in the output
  -o string
        Output file. If not specified, writes to stdout.
```

## Supported Backends

See [`helmfile/vals` documentation](https://github.com/helmfile/vals#supported-backends) for supported backends.
