![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dt/846963?style=plastic&logo=terraform&logoColor=%237B42BC&label=Terraform%20Downloads&color=%237B42BC&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2FCSC478-WCU%2Fcloudlab%2Flatest)

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
