# What Is CloudLab?

CloudLab is an **academic cloud testbed** that gives researchers full control over **bare-metal** and **virtualized** infrastructure. Unlike commercial clouds, CloudLab provides **root access** to hosts, programmable networks, and the freedom to run custom hypervisors, kernels, and images.

---

## Core Concepts

### Aggregates (Sites)
CloudLab resources are grouped into independent **aggregates** (sites). Common ones include:

- **Utah (APT/Emulab)** 
- **Wisconsin (Wisc)** 
- **Clemson** 
- **UMass** 

Each aggregate exposes compute, storage, and networking controlled through the **CloudLab Portal** and **XML-RPC API**. Aggregates can be targeted explicitly (via URNs) or selected by CloudLab’s scheduler when unspecified.

### Projects, Profiles, Experiments
- **Project** — your CloudLab organization (membership + permissions)  
- **Profile** — a **Python** blueprint that generates an **RSpec** (resource spec XML) using `geni.portal` and `geni.rspec.*` libraries  
- **Experiment** — a *live* instantiation of a profile, deployed to one or more aggregates and producing a **manifest** (hostnames, IPs, metadata)

### RSpec & Manifest
- **RSpec** — declarative XML describing nodes, links/LANs, images, blockstores, and constraints  
- **Manifest** — the realized allocation: nodes, control IPs, interfaces, per-aggregate status

---

## Resource Types & Topologies

### Nodes
- **RawPC (bare-metal)**: choose `hardware_type` (e.g., `d430`), set `exclusive`, select `disk_image` URN, and optional `routable_ip` for public control  
- **XenVM (virtual)**: define `cores`, `ram`, `disk`, and `instantiate_on` (host RawPC)

### Storage
- **Blockstores** (per node): name, mount point, size; normalized to strings like `200GB` in the profile

### Networking
- **Links**: point-to-point interfaces between nodes/VMs  
- **LANs**: broadcast segments for \>2 participants  
- **Bridged Links**: apply QoS — `bandwidth` (Mbps), `latency` (ms), `plr` (packet loss rate)

### Control Network
- **Control Interface** is separate from your experiment data links. Setting `routable_ip = true` on a node provides a globally routable control plane address (useful for SSH/Ansible).

---

## Lifecycle & Scheduling

1. **Instantiate** a profile with parameters (via Portal or API)  
2. **Allocate** resources — scheduler places nodes on aggregates  
3. **Provision** images, links, blockstores, VMs  
4. **Ready** — nodes are booted and reachable; a manifest is produced  
5. **Expire/Renew** — experiments have time limits; you can **extend** within policy

Status typically progresses:  
`provisioning → provisioned → creating → created → booting → booted → ready`

---

## Images & URNs

Disk images are referenced via **URNs**, e.g.:  
- `urn:publicid:IDN+emulab.net+image+Ubuntu22-64-STD`  
Custom images may be site-specific; verify availability on your target aggregate.

---

## Profiles in Practice

A minimal profile defining a RawPC:

```python
import geni.portal as portal
import geni.rspec.pg as RSpec

pc = portal.Context()
request = RSpec.Request()

n = request.RawPC("n1")
n.hardware_type = "d430"
# n.disk_image = "urn:publicid:IDN+emulab.net+image+Ubuntu22-64-STD"

pc.printRequestRSpec(request)
```

Profiles accept user parameters (node count, types, images), perform semantic checks, and emit a valid RSpec. CloudLab clones and runs the profile, then allocates accordingly.

---

## API & Security

- **Portal XML-RPC**: primary control plane; methods like `startExperiment`, `experimentStatus`, `terminateExperiment`  
- **Mutual TLS (mTLS)**: authenticate with your **PEM** containing private key + certificate  
- **Manifests/Status**: retrieved from the Portal and used by tooling (like this Terraform provider) to flatten **outputs** (URL, nodes map, status, expires)

---

## Why It Matters

CloudLab enables **system-level** experimentation (kernels, SDN, storage stacks) that commercial clouds restrict. Pairing CloudLab with **Terraform** brings professional reproducibility and CI/CD into academic workflows.
