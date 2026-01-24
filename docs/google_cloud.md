# Google Cloud

Protect your VM with a Google Cloud Firewall Rule. Create and associate a
Firewall Policy with a new or existing Google Cloud Instance and manage its
allowed IPv4 Addresses with fwsync.

## Prerequisites
1. GCP account
1. VM Instance
1. Firewall Rule associated with running Instance

## Authentication
The recommended method of authentication is to run the following command:

```bash
$ gcloud auth application-default login
```

## Quick Start

```
$ fwsync init --provider google --project YOUR_PROJECT
```

Whenever your ISP leases you a new IP, you can run `fwsync update` to seemlessly update your managed firewall rule.

