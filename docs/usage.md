# Usage — Terraform Provider for CloudLab

> This page explains **how to use** the CloudLab Terraform Provider to create, manage, and destroy experiments using Terraform.
> It assumes you already have a CloudLab account, a valid PEM certificate, and Terraform installed.

---

## 1. Installation

Add the provider to your `versions.tf`:

```hcl
terraform {
  required_version = ">= 1.5.0"
  required_providers {
    cloudlab = {
      source  = "CSC478-WCU/cloudlab"
      version = "~> 0.1"
    }
  }
}
```

Then initialize Terraform:

```bash
terraform init
```

---

## 2. Provider Configuration

In your `provider.tf`, configure the connection to the CloudLab portal:

```hcl
provider "cloudlab" {
  project  = "cloud-edu"            # Default project name
  pem_path = "~/cloudlab.pem"       # Path to your PEM certificate
  server   = "boss.emulab.net"      # CloudLab portal (default)
  port     = 3069
  path     = "/usr/testbed"
  timeout  = "15m"                  # Default timeout for portal calls
}
```

**Notes:**

* The PEM file authenticates to the XML-RPC API (not SSH).
* You can override any of these values at the resource level.

---

## 3. Basic Example

Create a simple experiment with one node:

```hcl
# main.tf
resource "cloudlab_portal_experiment" "demo" {
  name            = "tf-demo"
  project         = "cloud-edu"
  wait_for_status = "provisioned" # or "ready"

  rawpc = [{
    name          = "node0"
    hardware_type = "d430"
    exclusive     = true
    disk_image    = "urn:publicid:IDN+emulab.net+image+Ubuntu22-64-STD"
    aggregate     = "urn:publicid:IDN+apt.emulab.net+authority+cm"
    routable_ip   = true
  }]
}

output "experiment_url" {
  value = cloudlab_portal_experiment.demo.url
}

output "node_ips" {
  value = cloudlab_portal_experiment.demo.nodes
}
```

Then apply:

```bash
terraform apply
```

Terraform will:

1. Encode your experiment definition as JSON (`spec_json`).
2. Send it to CloudLab via the `portal.startExperiment` XML-RPC call.
3. Wait until the experiment reaches the requested status.
4. Return outputs such as node IPs and the experiment URL.

---

## 4. Multi-Node Example

```hcl
resource "cloudlab_portal_experiment" "cluster" {
  name  = "multi-node"
  project = "cloud-edu"

  rawpc = [
    {
      name          = "n1"
      hardware_type = "d430"
      aggregate     = "urn:publicid:IDN+apt.emulab.net+authority+cm"
      exclusive     = true
    },
    {
      name          = "n2"
      hardware_type = "d430"
      aggregate     = "urn:publicid:IDN+apt.emulab.net+authority+cm"
      exclusive     = true
    }
  ]

  link = [{
    name      = "l01"
    interface = [
      { node = "n1", ifname = "eth1" },
      { node = "n2", ifname = "eth1" }
    ]
  }]
}
```

---

## 5. Virtual Machines (XenVM)

Xen VMs can be instantiated on a host node using `instantiate_on`:

```hcl
rawpc = [{
  name          = "host1"
  hardware_type = "d430"
  exclusive     = true
}]

xenvm = [
  {
    name           = "vm01"
    cores          = 4
    ram_mb         = 8192
    disk_gb        = 50
    instantiate_on = "host1"
  }
]
```

---

## 6. Optional Resource Blocks

### Blockstores

Attach persistent storage to a node:

```hcl
blockstore = [
  { name = "bs0", mount = "/var/lib/db", size_gb = 100 }
]
```

### Bridged Links (QoS)

Add bandwidth/latency constraints:

```hcl
bridged_link = [{
  name           = "wan"
  bandwidth_mbps = 1000
  latency_ms     = 15
  plr            = 0.001
  interface = [
    { node = "n1" },
    { node = "n2" }
  ]
}]
```

---

## 7. Outputs

After `terraform apply`, outputs include:

* `uuid` — unique CloudLab experiment identifier
* `url` — portal URL
* `status` — current lifecycle phase
* `expires` — scheduled expiration time
* `nodes` — map of node name → IPv4 address

Example:

```text
uuid    = "5c7af02e-7c84-4e1d-bb92-d2d6938d31c8"
url     = "https://www.cloudlab.us/showexp.php?pid=cloud-edu&eid=tf-demo"
status  = "ready"
expires = "2025-10-25T12:00:00Z"
nodes = {
  "node0" = "155.98.123.45"
}
```

---

## 8. Destroying an Experiment

When finished:

```bash
terraform destroy
```

This will call:

```
portal.terminateExperiment
```

and free all allocated resources.

---

## 9. Parameters Reference

| Key                                             | Description                  | Default                  |
| ----------------------------------------------- | ---------------------------- | ------------------------ |
| `name`                                          | Experiment name              | *(required)*             |
| `project`                                       | CloudLab project             | Provider-level `project` |
| `pem_path`                                      | Path to PEM certificate      | `~/cloudlab.pem`         |
| `wait_for_status`                               | `"provisioned"` or `"ready"` | `"provisioned"`          |
| `timeout`                                       | Time limit for calls         | `"10m"`                  |
| `rawpc`, `xenvm`, `link`, `lan`, `bridged_link` | Experiment topology          | —                        |

---

## 10. Notes

* `wait_for_status = "ready"` ensures all aggregates and nodes are up and have IPv4 addresses before returning.
* All node kinds and link types directly map to the schema accepted by the **terraform-profile** repository.
* The provider is idempotent: if you rerun `terraform apply`, it reads the current state and updates only if necessary.

---

*2025 Tyler Geiger — CloudLab Terraform Provider Project*
*Designed for educational and research use; not affiliated with the official CloudLab team.*
