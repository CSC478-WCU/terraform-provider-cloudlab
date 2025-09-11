package experiment

import (
	"github.com/csc478-wcu/terraform-provider-cloudlab/internal/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildSpec(d *schema.ResourceData) model.ExperimentSpec {
	var spec model.ExperimentSpec
	// rawpc
	for _, v := range toList(d.Get("rawpc")) {
		m := v.(map[string]interface{})
		spec.Nodes = append(spec.Nodes, model.Node{
			Kind:         "rawpc",
			Name:         s(m["name"]),
			HardwareType: s(m["hardware_type"]),
			Exclusive:    pBool(m, "exclusive"),
			DiskImage:    s(m["disk_image"]),
			Aggregate:    s(m["aggregate"]),
			RoutableIP:   pBool(m, "routable_ip"),
			Blockstores:  expandBlockstores(m["blockstore"]),
		})
	}
	// xenvm
	for _, v := range toList(d.Get("xenvm")) {
		m := v.(map[string]interface{})
		spec.Nodes = append(spec.Nodes, model.Node{
			Kind:          "xenvm",
			Name:          s(m["name"]),
			Cores:         pInt(m, "cores"),
			RamMB:         pInt(m, "ram_mb"),
			DiskGB:        pInt(m, "disk_gb"),
			InstantiateOn: s(m["instantiate_on"]),
			DiskImage:     s(m["disk_image"]),
			Aggregate:     s(m["aggregate"]),
			RoutableIP:    pBool(m, "routable_ip"),
			Blockstores:   expandBlockstores(m["blockstore"]),
		})
	}
	// links
	spec.Links = append(spec.Links, expandLinks("link", d.Get("link"))...)
	spec.Links = append(spec.Links, expandLinks("lan", d.Get("lan"))...)
	spec.Links = append(spec.Links, expandBridged(d.Get("bridged_link"))...)
	return spec
}

func expandBlockstores(v interface{}) []model.Blockstore {
	var out []model.Blockstore
	for _, it := range toList(v) {
		m := it.(map[string]interface{})
		out = append(out, model.Blockstore{
			Name:  s(m["name"]),
			Mount: s(m["mount"]),
			Size:  m["size_gb"].(int),
		})
	}
	return out
}

func expandLinks(kind string, v interface{}) []model.Link {
	var out []model.Link
	for _, it := range toList(v) {
		m := it.(map[string]interface{})
		out = append(out, model.Link{
			Kind:       kind,
			Name:       s(m["name"]),
			Interfaces: expandIfaces(m["interface"]),
		})
	}
	return out
}

func expandBridged(v interface{}) []model.Link {
	var out []model.Link
	for _, it := range toList(v) {
		m := it.(map[string]interface{})
		l := model.Link{
			Kind:       "bridged_link",
			Name:       s(m["name"]),
			Interfaces: expandIfaces(m["interface"]),
			Bandwidth:  pInt(m, "bandwidth_mbps"),
			Latency:    pInt(m, "latency_ms"),
			Plr:        pFloat(m, "plr"),
		}
		out = append(out, l)
	}
	return out
}

func expandIfaces(v interface{}) []model.Iface {
	var out []model.Iface
	for _, it := range toList(v) {
		m := it.(map[string]interface{})
		out = append(out, model.Iface{
			Node:   s(m["node"]),
			IfName: s(m["ifname"]),
		})
	}
	return out
}

// small helpers
func toList(v interface{}) []interface{} { if v == nil { return nil }; return v.([]interface{}) }
func s(v interface{}) string             { if v == nil { return "" }; return v.(string) }
func pInt(m map[string]interface{}, k string) *int {
	if v, ok := m[k]; ok && v != nil {
		iv := v.(int)
		return &iv
	}
	return nil
}
func pBool(m map[string]interface{}, k string) *bool {
	if v, ok := m[k]; ok && v != nil {
		bv := v.(bool)
		return &bv
	}
	return nil
}
func pFloat(m map[string]interface{}, k string) *float64 {
	if v, ok := m[k]; ok && v != nil {
		fv := v.(float64)
		return &fv
	}
	return nil
}
