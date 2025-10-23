
# CI/CD with GitHub Actions (Terraform → CloudLab → Ansible) 

This tutorial shows a **generic** CI/CD setup that:

1) Provisions a CloudLab cluster via **Terraform** (using this provider),

2) Exports node IPs from Terraform outputs,

3) Invokes a **reusable Ansible** workflow that does a minimal remote action on the nodes (touches a file).

!!! example "Reference example & motivation"
    Example repository: **[github.com/csc478-wcu/cc2-cluster](https://github.com/csc478-wcu/cc2-cluster)**  
    That repo hosts a campus project integrating Terraform + CloudLab + CI/CD + Ansible.  
    This tutorial abstracts those ideas into a **minimal, reusable** template for any project.

---

## Secrets Required

Store these in your GitHub repository **Secrets and variables → Actions → Repository secrets**:

- **`CLOUDLAB_PEM_B64`** — Base64-encoded CloudLab **decrypted** PEM (private key + certificate) used for **mTLS** to the Portal and **SSH** to nodes.

!!! warning
    Do **not** commit PEM files. Always pass keys via GitHub Secrets.

---

## Repository Layout

You can drop the following files into your repo:

```

.
├── infra/
│   ├── main.tf
│   ├── outputs.tf
│   ├── provider.tf
│   ├── variables.tf
│   └── versions.tf
├── ansible/
│   ├── ansible.cfg
│   ├── inventory.ini
│   └── playbooks/
│       └── bootstrap.yml
└── .github/
    └── workflows/
        ├── terraform.yml
        └── ansible-deploy.yml

```

---

## 1) Terraform configuration (`infra/`)

### `versions.tf`

```hcl
terraform {
  required_version = ">= 1.5.0"
  required_providers {
    cloudlab = {
      source  = "csc478-wcu/cloudlab"
      version = ">= 1.0.4"
    }
  }
}
```

### `provider.tf`

```hcl
provider "cloudlab" {
  project  = "cloud-edu"
  server   = "boss.emulab.net"
  port     = 3069
  path     = "/usr/testbed"
  pem_path = "./cloudlab_decrypted.pem"  # written by the workflow
  timeout  = "30m"
}
```

### `variables.tf`

```hcl
variable "aggregate" {
  description = "Aggregate short form (matches GH Actions choice)"
  type        = string
}

variable "hardware_type" {
  description = "Hardware type for rawpc nodes"
  type        = string
}

variable "num_workers" {
  description = "Number of worker nodes to create"
  type        = number
  default     = 2
  validation {
    condition     = var.num_workers >= 1 && var.num_workers <= 20
    error_message = "num_workers must be between 1 and 20."
  }
}

variable "aggregate_map" {
  description = "Map aggregate short names to URNs"
  type        = map(string)
  default = {
    "emulab.net"          = "urn:publicid:IDN+emulab.net+authority+cm"
    "utah.cloudlab.us"    = "urn:publicid:IDN+utah.cloudlab.us+authority+cm"
    "clemson.cloudlab.us" = "urn:publicid:IDN+clemson.cloudlab.us+authority+cm"
    "wisc.cloudlab.us"    = "urn:publicid:IDN+wisc.cloudlab.us+authority+cm"
    "apt.emulab.net"      = "urn:publicid:IDN+apt.emulab.net+authority+cm"
  }
}
```

### `main.tf`

This creates **1 admin node** (`kubeadm`) and **N workers**, places all nodes on the chosen aggregate, and builds a single **LAN** connecting them.  
(Choose a hardware type that exists at the selected site.)

```hcl
resource "cloudlab_portal_experiment" "experiment_name" {
  name            = "your_exp_name"
  wait_for_status = "ready"

  # kubeadm + worker1..workerN
  dynamic "rawpc" {
    for_each = toset(concat(
      ["kubeadm"],
      [for i in range(var.num_workers) : "worker${i + 1}"]
    ))
    content {
      name          = rawpc.key
      hardware_type = var.hardware_type
      aggregate     = var.aggregate_map[var.aggregate]
      exclusive     = true
      routable_ip   = true
    }
  }

  # LAN with all nodes
  lan {
    name = "lan0"

    dynamic "interface" {
      for_each = toset(concat(
        ["kubeadm"],
        [for i in range(var.num_workers) : "worker${i + 1}"]
      ))
      content {
        node = interface.key
      }
    }
  }
}
```

### `outputs.tf`

```hcl
output "url" {
  value = cloudlab_portal_experiment.experiment_name.url
}

output "nodes" {
  value = cloudlab_portal_experiment.experiment_name.nodes
}
```

---

## 2) Minimal Ansible (`ansible/`)

We’ll keep Ansible **generic**: it reads the Terraform `nodes` output (host → IP map), builds an inventory dynamically, and **touches a file** on each node to show connectivity.

### `ansible.cfg`

```ini
[defaults]
host_key_checking = False
forks = 20
timeout = 600
interpreter_python = auto_silent
retry_files_enabled = False
inventory = ./ansible/inventory.ini

[ssh_connection]
pipelining = True
ssh_args = -o ControlMaster=auto -o ControlPersist=600s \
           -o ServerAliveInterval=30 -o ServerAliveCountMax=10
```

### `inventory.ini`

```ini
[local]
localhost ansible_connection=local
```

### `playbooks/bootstrap.yml`

```yaml
---
- hosts: localhost
  connection: local
  gather_facts: false
  vars:
    nodes: "{{ nodes_json | from_json }}"
    ssh_user: "{{ ssh_user | default('your_username') }}" # Change default to your cloudlab username
  tasks:
    - name: Register CloudLab nodes from Terraform outputs
      add_host:
        name: "{{ item.key }}"
        ansible_host: "{{ item.value }}"
        ansible_user: "{{ ssh_user }}"
        groups: cloudlab_nodes
      loop: "{{ nodes | dict2items }}"

- hosts: cloudlab_nodes
  become: yes
  tasks:
    - name: Ensure /local/repository exists
      ansible.builtin.file:
        path: /local/repository
        state: directory
        mode: "0755"

    - name: Touch a marker file to prove end-to-end CI
      ansible.builtin.file:
        path: /local/repository/touched
        state: touch
        mode: "0644"
```


---

## 3) GitHub Actions (`.github/workflows/`)

### `terraform.yml` — main workflow

Runs on manual dispatch, provisions CloudLab via Terraform, and passes the node map to the Ansible workflow.

```yaml
name: Terraform Apply Cloudlab

on:
  workflow_dispatch:
    inputs:
      branch:
        description: "Branch to deploy"
        required: true
        default: "main"
      ssh_user:
        description: "SSH user on nodes (CloudLab username)"
        required: true
        default: "geni"
      num_workers:
        description: "Number of worker nodes"
        required: true
        type: number
        default: 2
      aggregate:
        description: "Select CloudLab aggregate"
        required: true
        type: choice
        options:
          - "emulab.net"
          - "utah.cloudlab.us"
          - "clemson.cloudlab.us"
          - "wisc.cloudlab.us"
          - "apt.emulab.net"
        default: "emulab.net"
      hardware_type:
        description: "Select hardware type (must exist at chosen aggregate)"
        required: true
        default: "d430"

jobs:
  terraform:
    runs-on: ubuntu-latest
    outputs:
      nodes_json: ${{ steps.collect.outputs.nodes_json }}
      url: ${{ steps.collect.outputs.url }}

    steps:
      - name: Checkout repo
        uses: actions/checkout@v4
        with:
          ref: ${{ inputs.branch }}

      - name: Create CloudLab PEM from secret
        working-directory: ./infra
        run: |
          set -euo pipefail
          echo "${CLOUDLAB_PEM_B64}" | base64 -d > cloudlab_decrypted.pem
          chmod 600 cloudlab_decrypted.pem
          head -n1 cloudlab_decrypted.pem | grep -q -- 'BEGIN' || { echo "PEM decode failed"; exit 1; }
        env:
          CLOUDLAB_PEM_B64: ${{ secrets.CLOUDLAB_PEM_B64 }}

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: 1.7.2
          terraform_wrapper: false

      - name: Terraform Init
        working-directory: ./infra
        run: terraform init

      - name: Terraform Apply
        working-directory: ./infra
        run: |
          terraform apply -auto-approve \
            -var="aggregate=${{ inputs.aggregate }}" \
            -var="hardware_type=${{ inputs.hardware_type }}" \
            -var="num_workers=${{ inputs.num_workers }}"

      - name: Show outputs (human)
        working-directory: ./infra
        run: terraform output

      - name: Install jq
        run: |
          sudo apt-get update
          sudo apt-get install -y jq

      - name: Collect outputs for Ansible
        id: collect
        working-directory: ./infra
        shell: bash
        run: |
          set -euo pipefail
          tf_json="$(terraform output -json)"
          echo "$tf_json" | jq -e '.' >/dev/null
          nodes_json="$(echo "$tf_json" | jq -c '.nodes.value')"
          url="$(echo "$tf_json" | jq -r '.url.value')"
          {
            echo "nodes_json<<EOF"
            echo "$nodes_json"
            echo "EOF"
          } >> "$GITHUB_OUTPUT"
          echo "url=$url" >> "$GITHUB_OUTPUT"

  ansible:
    needs: terraform
    uses: ./.github/workflows/ansible-deploy.yml
    with:
      nodes_json: ${{ needs.terraform.outputs.nodes_json }}
      ssh_user: ${{ inputs.ssh_user }}
      branch: ${{ inputs.branch }}
    secrets:
      CLOUDLAB_PEM_B64: ${{ secrets.CLOUDLAB_PEM_B64 }}
```

### `ansible-deploy.yml` — reusable workflow

```yaml
name: Ansible Deploy (Reusable)

on:
  workflow_call:
    inputs:
      nodes_json: { required: true, type: string } # Terraform outputs (map hostname -> IP)
      ssh_user:   { required: true, type: string } # SSH username
      branch:     { required: true, type: string } # Git branch
    secrets:
      CLOUDLAB_PEM_B64: { required: true }         # Base64-encoded PEM

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repo
        uses: actions/checkout@v4
        with:
          ref: ${{ inputs.branch }}

      - name: Recreate SSH key from secret
        run: |
          echo "${CLOUDLAB_PEM_B64}" | base64 -d > cloudlab_decrypted.pem
          chmod 600 cloudlab_decrypted.pem
        env:
          CLOUDLAB_PEM_B64: ${{ secrets.CLOUDLAB_PEM_B64 }}

      - name: Install Ansible
        run: pip install "ansible>=9,<11"

      - name: Touch a file on remote nodes via playbook
        env:
          ANSIBLE_CONFIG: ${{ github.workspace }}/ansible/ansible.cfg
        run: |
          ansible-playbook ansible/playbooks/bootstrap.yml \
            -e nodes_json='${{ inputs.nodes_json }}' \
            -e ssh_user='${{ inputs.ssh_user }}' \
            --private-key "${GITHUB_WORKSPACE}/cloudlab_decrypted.pem"
```

---

## Using the Workflow

1. Push the files above.
2. Add the `CLOUDLAB_PEM_B64` secret.
3. Open **Actions → Terraform Apply Cloudlab → Run workflow**.
4. Pick:
   - **branch** (default: `main`)
   - **ssh_user** (your CloudLab username)
   - **num_workers** (e.g., 2)
   - **aggregate** (e.g., `emulab.net`)
   - **hardware_type** (must exist at chosen aggregate, e.g., `d430` at Emulab)
5. Watch Terraform create the experiment → Ansible touch a file on each node.

---

## Notes & Tips

- **Hardware & aggregate:** pick combinations that exist (e.g., `d430` @ Emulab, `r320` @ APT, `c6320` @ Clemson, `m510` @ Utah, `c220g5` @ Wisc).  
- **Idempotency:** This provider’s experiment topology is `ForceNew`; edits will recreate the experiment. Use `wait_for_status = "provisioned"` for faster iteration if ansible / ssh is not necessary.  
- **Security:** Keep PEMs in Secrets, not in the repo.  
- **Extending:** Replace the Ansible task with your real bootstrap (kubeadm, CNI, etc.) when ready.

That’s it — a simple, professional CI/CD path from **GitHub → Terraform → CloudLab → Ansible** that you can adapt to any project.

