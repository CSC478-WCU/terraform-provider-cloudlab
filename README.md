<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

[![GitHub Tag](https://img.shields.io/github/v/tag/CSC478-WCU/terraform-provider-cloudlab?style=plastic&logo=terraform&logoColor=%23844FBA&label=latest&color=%23844FBA&link=https%3A%2F%2Fgithub.com%2FCSC478-WCU%2Fterraform-provider-cloudlab%2Freleases)](https://github.com/CSC478-WCU/terraform-provider-fabric/releases) [![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dt/846963?style=plastic&logo=terraform&logoColor=%237B42BC&label=Terraform%20Downloads&color=%237B42BC&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2FCSC478-WCU%2Fcloudlab%2Flatest)](https://registry.terraform.io/providers/CSC478-WCU/cloudlab)

# CloudLab Terraform Provider

A Terraform provider that automates **CloudLab** experiments using standard IaC workflows.

**Full documentation:** [csc478-wcu.github.io/terraform-provider-cloudlab](https://csc478-wcu.github.io/terraform-provider-cloudlab)

---

## Quick Start

Add the provider and configure credentials:

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

provider "cloudlab" {
  project  = "your-project"
  pem_path = "~/cloudlab_decrypted.pem"
  server   = "boss.emulab.net"
  port     = 3069
  path     = "/usr/testbed"
  timeout  = "15m"
}

resource "cloudlab_portal_experiment" "demo" {
  name            = "tf-demo"
  project         = "your-project"
  wait_for_status = "provisioned"

  rawpc {
    name          = "node0"
    hardware_type = "r320"  # e.g., APT
    exclusive     = true
    aggregate     = "urn:publicid:IDN+apt.emulab.net+authority+cm"
    routable_ip   = true
  }
}
````

```bash
terraform init
terraform apply
```

---

## Links

* Docs: [https://csc478-wcu.github.io/terraform-provider-cloudlab/#overview](https://csc478-wcu.github.io/terraform-provider-cloudlab/#overview)
* Registry: [https://registry.terraform.io/providers/CSC478-WCU/cloudlab](https://registry.terraform.io/providers/CSC478-WCU/cloudlab)
* Releases: [https://github.com/CSC478-WCU/terraform-provider-cloudlab/releases](https://github.com/CSC478-WCU/terraform-provider-cloudlab/releases)

---

## License

MIT — see [`LICENSE`](./LICENSE).
