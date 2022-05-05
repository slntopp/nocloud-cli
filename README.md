# NoCloud CLI

Official CLI for [NoCloud Platform](https://github.com/slntopp/nocloud)

## Table of Contents

* [Usage](#usage)
* [Installation](#installation)
* [Notes](#notes)
* [Example Configs](#example-configs)

## Usage

Check out `nocloud -h` for list of available commands.
You can also use `nocloud help [command]` to see help notes for particular command.

You must authorize first, to use other commands, see `nocloud help login`.

## Installation

[See CLI reference at main repo page](https://github.com/slntopp/nocloud#nocloud-cli)

## Notes

### Login

> CLI currently supports only `standard` type of Authorization

### Auto-Completion

> CLI is based on [spf13/cobra](https://github.com/spf13/cobra) and supports auto-completion,
> you can see register scripts by running `nocloud completion [your shell]`
> To check how to register completions, add `-h` flag

## Example Configs

### Setting

Use this config for `nocloud settings apply`

```yaml
key: setting-from-yaml
value: not a json
description: sample setting
public: false
```

### DNS

Use this config for `nocloud dns apply`

```yaml
name: cluster.nocloud.
locations:
  tunnel:
    cname:
    - host: tunnelserver
      ttl: 360
    txt:
    - text: I'm a DNS record
      ttl: 360
```
