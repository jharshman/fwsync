# fwsync

[![CI](https://github.com/jharshman/fwsync/actions/workflows/ci.yaml/badge.svg)](https://github.com/jharshman/fwsync/actions/workflows/ci.yaml)

Provides CLI interface to update your personal Firewall Rules
associated with your Cloud Development VM.

## Installation

TODO how to install
Install by running the following in your terminal:
```bash
$ curl https://github.com/jharshman/fwsync/... | sh
```

## Usage

### Init
After installing, you can invoke the CLI by typing `fwsync` in your terminal.
This by default will display some usage information.

To initialize fwsync type `fwsync init`. This will walk you through steps in
selecting the correct firewall to manage and will write out fwsync's config file
which will be located at `$HOME/.bitly_firewall`.

### Update
If your IP updates and you notice  you've lost access to your CloudVM,
you can invoke `fwsync update` to automatically detect your new IP address
and update your Firewall Rule.