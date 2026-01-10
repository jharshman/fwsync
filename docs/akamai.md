# Akamai (Formerly Linode)

Protect your VM on Akamai. Create and associate a Firewall Policy
with a new or existing Linode Instance and manage its allowed 
IPv4 Addresses with fwsync.

## Prerequisites
1. Linode account
1. API Key
1. Linode Instance
1. Firewall Rule associated with running Instance

## Authentication
To authenticate with Linode, login to your account and create and copy
a new API Key. Set the `LINODE_TOKEN` environment variable for your shell.

## Quick Start

```
# Keep this variable exported in your shells's rc file.
$ export LINODE_TOKEN="YOUR_LINODE_API_TOKEN"
$ fwsync init --provider linode
```

Whenever your ISP leases you a new IP, you can run `fwsync update` to seemlessly update your managed firewall rule.

