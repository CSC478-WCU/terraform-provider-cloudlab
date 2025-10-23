# Getting Started

If you’re new to Terraform, first read [What Is Terraform](what-is-terraform.md) to understand how configuration files, providers, and resources work together.

This guide walks you through installing Terraform, connecting your CloudLab account, and launching your first experiment.

!!! tip
    For full examples and resource schemas, see the [Usage Guide](usage.md).

---

## Step 1 — Install Terraform

You’ll need **Terraform ≥ 1.5.0**.

=== "Linux"

    ```bash
    sudo apt-get update && sudo apt-get install -y gnupg software-properties-common curl
    curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
      https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
    sudo apt update && sudo apt install terraform
    terraform -version
    ```

=== "macOS"

    ```bash
    brew tap hashicorp/tap
    brew install hashicorp/tap/terraform
    terraform -version
    ```

=== "Windows (PowerShell)"

    ```powershell
    choco install terraform
    terraform -version
    ```

---

## Step 2 — Download and Decrypt Your PEM

Follow the [Authentication Guide](auth.md) to obtain and decrypt your **CloudLab PEM certificate**.

You’ll need the path to this decrypted PEM when configuring the provider.

---

## Step 3 — Create a Terraform Workspace

```bash
mkdir my-experiment && cd my-experiment
touch main.tf provider.tf versions.tf
```

---

## Step 4 — Define Terraform and Provider Versions

`versions.tf`:

```hcl
terraform {
  required_version = ">= 1.5.0"
  required_providers {
    cloudlab = {
      source  = "CSC478-WCU/cloudlab"
      version = ">= 1.0.4"
    }
  }
}
```

---

## Step 5 — Configure the Provider

`provider.tf`:

```hcl
provider "cloudlab" {
  project  = "your-project"
  pem_path = "~/cloudlab_decrypted.pem"
  server   = "boss.emulab.net"
  port     = 3069
  path     = "/usr/testbed"
  timeout  = "15m"
}
```

---

## Step 6 — Define Your First Experiment

`main.tf`:

```hcl
resource "cloudlab_portal_experiment" "demo" {
  name            = "tf-demo"
  project         = "your-project"
  wait_for_status = "provisioned"

  rawpc {
    name          = "node0"
    hardware_type = "r320"
    exclusive     = true
    aggregate     = "urn:publicid:IDN+apt.emulab.net+authority+cm"
    routable_ip   = true
  }
}
```

---

## Step 7 — Initialize Terraform

```bash
terraform init
```

---

## Step 8 — Apply

```bash
terraform apply
```

Terraform will create the experiment and output the CloudLab URL and node IPs.

```
Apply complete! Resources: 1 added.

Outputs:

experiment_url = "https://www.cloudlab.us/showexp.php?pid=your-project&eid=tf-demo"
nodes = {
  "node0" = "155.98.123.45"
}
```

---

## Step 9 — Teardown and Destroy

When done:

```bash
terraform destroy
```

---

**Next Steps:**  
See the [Usage Guide](usage.md) for full examples, resource schemas, and advanced experiment definitions.
