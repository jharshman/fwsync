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

## How fwsync manages the Firewall
fwsync manages a single inbound rule labeled `fwsync` that allows TCP traffic
from your IPs. If the rule doesn't exist it is created. Every other rule, and the
Firewall's inbound and outbound policies, are left untouched. You can edit the
`fwsync` rule's ports or protocol in Cloud Manager and fwsync will keep them,
only replacing its IPv4 addresses.

Make sure the Firewall's inbound policy is set to **DROP**. With an inbound
policy of ACCEPT all traffic is allowed and the `fwsync` rule has no effect.
fwsync prints a warning when it detects this.

> Note: fwsync v0.0.3 and earlier replaced all of the Firewall's rules and set
> the inbound policy to ACCEPT. If you used one of those versions, check your
> Firewall's inbound policy and re-add any rules that were removed.

Whenever your ISP leases you a new IP, you can run `fwsync update` to seemlessly update your managed firewall rule.

