# Terraform Basics

Terraform is an **Infrastructure-as-Code (IaC)** tool developed by **HashiCorp** that allows you to define, provision, and manage infrastructure using simple, **declarative configuration files**.

Instead of manually creating servers or networks through a web portal, you write configuration files that describe **what** you want — and Terraform determines **how** to create or update it.

!!! info
    The CloudLab Terraform Provider communicates with the CloudLab Portal using **XML-RPC** over **mutual TLS (mTLS)**, making Terraform a powerful automation layer on top of the academic testbed.

---

## Why Terraform?

Terraform enables **reproducible**, **version-controlled**, and **automated** infrastructure deployment.

It’s especially valuable for environments like CloudLab, where resources such as nodes, links, and profiles can be recreated automatically from configuration files — improving collaboration, teaching reproducibility, and experiment repeatability.

!!! example "Key Advantages"
    - **Reproducibility** — Every experiment can be recreated exactly as defined in code.  
    - **Automation** — Integrate with CI/CD pipelines for consistent deployments.  
    - **Version Control** — Use Git to track infrastructure changes.  
    - **Portability** — Works across providers (CloudLab, AWS, Azure, GCP).  
    - **Safety** — Terraform plans and previews every change before applying it.

---

## Core Terraform Concepts

| Concept | Description |
|---|---|
| **Provider** | A plugin that exposes Terraform resources for a specific platform (e.g., CloudLab, AWS, Azure). |
| **Resource** | A single managed object — such as a CloudLab experiment, node, or link. |
| **Data Source** | A read-only object that queries existing information (e.g., aggregate lists or existing profiles). |
| **Module** | A reusable collection of Terraform resources (similar to a function). |
| **State** | A local or remote file that tracks the current deployed infrastructure and maps real-world objects to your Terraform configuration. |
| **Plan** | The preview stage where Terraform shows what will change. |
| **Apply** | Executes the actions needed to reach the desired state. |
| **Destroy** | Removes managed resources and cleans up the environment. |

---

## Example Workflow

```bash
# 1. Initialize provider plugins
terraform init

# 2. Preview changes (dry run)
terraform plan

# 3. Apply configuration
terraform apply

# 4. Review outputs and experiment URLs
terraform output

# 5. Clean up resources
terraform destroy
```

For most providers terraform executes these steps idempotently — meaning the same configuration can be applied multiple times, and only drifted resources will be updated. Cloudlab experiment modification is hopefully coming soon. 

---

## Declarative Configuration Example

Here’s a minimal example for deploying one node on CloudLab:

```hcl
terraform {
  required_providers {
    cloudlab = {
      source  = "CSC478-WCU/cloudlab"
      version = ">= 1.0.4"
    }
  }
}

provider "cloudlab" {
  project  = "your-project"
  pem_path = "~/cloudlab_decrypted.pem"
}

resource "cloudlab_portal_experiment" "example" {
  name            = "demo"
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

## How Terraform Works (Under the Hood)

1. **Initialization (`terraform init`)** — Downloads the required provider plugins (e.g., the CloudLab provider).  
2. **Planning (`terraform plan`)** — Compares your configuration with the current state file.  
3. **Execution (`terraform apply`)** — Invokes the provider’s APIs to provision or modify resources.  
4. **State Management** — Saves real-world resource IDs (UUIDs, URLs, IPs) in the `.tfstate` file.  
5. **Destruction (`terraform destroy`)** — Calls provider APIs to remove resources gracefully.

In the CloudLab provider, the “apply” step sends a **JSON spec** to the CloudLab Portal via **XML-RPC**, waits for provisioning, and returns experiment metadata back into Terraform outputs.

---

## Learn More

- [Terraform Official Introduction — developer.hashicorp.com](https://developer.hashicorp.com/terraform/intro)  
- [Terraform Overview — GeeksforGeeks](https://www.geeksforgeeks.org/devops/what-is-terraform/)  

Both are excellent resources for learning Terraform fundamentals and best practices.

---


