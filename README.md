<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

[![GitHub Tag](https://img.shields.io/github/v/tag/CSC478-WCU/terraform-provider-cloudlab?style=plastic&logo=terraform&logoColor=%23844FBA&label=latest&color=%23844FBA&link=https%3A%2F%2Fgithub.com%2FCSC478-WCU%2Fterraform-provider-cloudlab%2Freleases)](https://github.com/CSC478-WCU/terraform-provider-fabric/releases) [![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dt/846963?style=plastic&logo=terraform&logoColor=%237B42BC&label=Terraform%20Downloads&color=%237B42BC&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2FCSC478-WCU%2Fcloudlab%2Flatest)](https://registry.terraform.io/providers/CSC478-WCU/cloudlab)

# Terraform Provider for Cloudlab 

```
provider "cloudlab" {
  project  = "YOUR_PROJECT"         # required (here or on the resource)
  pem_path = "~/.ssl/emulab.pem"
  server   = "boss.emulab.net"
  port     = 3069
  path     = "/usr/testbed"
  timeout  = "30m"
}

resource "cloudlab_portal_experiment" "demo" {
  name            = "tf-demo-001"
  wait_for_status = "provisioned"   # or "ready"

  rawpc {
    name           = "host1"
    hardware_type  = "d430"
    exclusive      = true
    routable_ip    = true

    blockstore {
      name    = "bs1"
      mount   = "/data1"
      size_gb = 20
    }
  }

  xenvm {
    name           = "vm1"
    cores          = 2
    ram_mb         = 4096
    disk_gb        = 20
    instantiate_on = "host1"
    routable_ip    = true
  }

  lan {
    name = "lan0"
    interface { node = "host1" }
    interface { node = "vm1" }
  }
}

output "url"   { value = cloudlab_portal_experiment.demo.url }
output "nodes" { value = cloudlab_portal_experiment.demo.nodes }
```
