# nocloud-cli
Official CLI for [NoCloud Platform](https://github.com/slntopp/nocloud)

## Table of Contents

* [Usage](#usage)
* [Installation](#installation)
* [Notes](#notes)

## Usage

Check out `nocloud -h` for list of available commands.
You can also use `nocloud help [command]` to see help notes for particular command.

You must authorize first, to use other commands, see `nocloud help login`.

## Installation

### Download Pre-Built

1. Navigate to Releases pages
2. Pick CLI Release Tag matching NoCloud version you running
3. Download archived binary and checksum for you platform
4. Unpack archive: `tar -xvzf nocloud-<version>-<os>-<arch>.tar.gz`
5. Move `nocloud` binary file to `/usr/local/bin`
6. You're all set!

> macOS users might'd need to allow it in Privacy Settings

### Building from Source

> You must have golang (version `1.17`) environment set

1. Clone this repo
2. Run `go build -o nocloud`
3. Put freshly built `nocloud` binary file to `/usr/local/bin`
4. You're all set!

> macOS users might'd need to allow it in Privacy Settings

### Login

> CLI currently supports only `standard` type of Authorization

### Auto-Completion

> CLI is based on [spf13/cobra](https://github.com/spf13/cobra) and supports auto-completion,
> you can see register scripts by running `nocloud completion [your shell]`
> To check how to register completions, add `-h` flag
