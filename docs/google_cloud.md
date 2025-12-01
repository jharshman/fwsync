# Google Cloud

Protect your VM with a Google Cloud Firewall Rule. Create and associate a
Firewall Policy with a new or existing Google Cloud Instance and manage its
allowed IPv4 Addresses with fwsync.

## Authentication
The recommended method of authentication is to run the following command:

```bash
$ gcloud auth application-default login
```

## Quick Start

Create an instance or use an existing instance:
```bash
gcloud compute instances create my-dev-vm \
  --zone <ZONE> \
  --project <PROJECT> \
  --machine-type <MACHINE_TYPE> \
  --image-family <IMAGE_FAMILY> \
  --tags my-dev-vm
```

Create and associate to the instance:
```bash
$ gcloud compute firewall-rules create allow-dev-vm \
  --allow TCP \
  --direction INGRESS \
  --network <NETWORK> \
  --source-ranges <YOUR PUBLIC IP> \
  --target-tags <my-dev-vm> \
  --project <PROJECT>
```

> **Important:** Instance association is done via Instance Tag and Target Tag.
The Target Tag on the Firewall must match one of the defined Tags
on the Instance.

