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

NoCloud CLI currently doesn't offer any official way to install CLI(published binaries or installers).
You can clone this repo, run `go build -o nocloud` and put `nocloud` binary to `/usr/bin`. You're all set!

> macOS users might'd need to allow it in Privacy Settings

### Login

> CLI currently supports only `standard` type of Authorization

### Auto-Completion

> CLI is based on [spf13/cobra](https://github.com/spf13/cobra) and supports auto-completion,
> you can see register scripts by running `nocloud completion [your shell]`
> To check how to register completions, add `-h` flag
