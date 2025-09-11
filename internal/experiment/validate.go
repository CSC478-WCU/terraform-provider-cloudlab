package experiment

import (
	"fmt"

	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/model"
	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/validation"
)

func validateSpec(s model.ExperimentSpec) error {
	names := map[string]string{}
	for _, n := range s.Nodes {
		if n.Name == "" {
			return fmt.Errorf("node name is required")
		}
		if _, dup := names[n.Name]; dup {
			return fmt.Errorf("duplicate node name: %s", n.Name)
		}
		names[n.Name] = n.Kind

		// numeric
		if n.Kind == "xenvm" {
			if n.Cores != nil && *n.Cores < 1 {
				return fmt.Errorf("xenvm %q cores must be >= 1", n.Name)
			}
			if n.RamMB != nil && *n.RamMB < 1 {
				return fmt.Errorf("xenvm %q ram_mb must be >= 1", n.Name)
			}
			if n.DiskGB != nil && *n.DiskGB < 1 {
				return fmt.Errorf("xenvm %q disk_gb must be >= 1", n.Name)
			}
		}
		// blockstores
		for _, b := range n.Blockstores {
			if b.Name == "" {
				return fmt.Errorf("node %q blockstore missing name", n.Name)
			}
			if b.Size < 1 {
				return fmt.Errorf("node %q blockstore %q size_gb must be >= 1", n.Name, b.Name)
			}
		}
		// aggregate (optional) must be in list when set
		if n.Aggregate != "" && !validation.IsValidAggregate(n.Aggregate) {
			return fmt.Errorf("node %q aggregate %q is not a recognized URN", n.Name, n.Aggregate)
		}
	}

	for _, n := range s.Nodes {
		if n.Kind == "xenvm" && n.InstantiateOn != "" {
			if k, ok := names[n.InstantiateOn]; !ok || k != "rawpc" {
				return fmt.Errorf("xenvm %q instantiate_on must reference an existing rawpc (got %q)", n.Name, n.InstantiateOn)
			}
		}
	}

	for _, l := range s.Links {
		if l.Name == "" {
			return fmt.Errorf("link name is required")
		}
		if len(l.Interfaces) < 2 {
			return fmt.Errorf("%s %q must have at least 2 interfaces", l.Kind, l.Name)
		}
		if l.Kind == "link" && len(l.Interfaces) != 2 {
			return fmt.Errorf("link %q must have exactly 2 interfaces", l.Name)
		}
		for _, ifc := range l.Interfaces {
			if _, ok := names[ifc.Node]; !ok {
				return fmt.Errorf("link %q references unknown node %q", l.Name, ifc.Node)
			}
		}
		if l.Plr != nil && (*l.Plr < 0.0 || *l.Plr > 1.0) {
			return fmt.Errorf("bridged_link %q plr must be between 0.0 and 1.0", l.Name)
		}
	}
	return nil
}
