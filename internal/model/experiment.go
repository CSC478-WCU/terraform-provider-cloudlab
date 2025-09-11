package model

type ExperimentSpec struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links,omitempty"`
}

type Node struct {
	Kind          string       `json:"kind"` // "rawpc" | "xenvm"
	Name          string       `json:"name"`
	HardwareType  string       `json:"hardware_type,omitempty"` // rawpc
	Exclusive     *bool        `json:"exclusive,omitempty"`     // rawpc
	Cores         *int         `json:"cores,omitempty"`         // xenvm
	RamMB         *int         `json:"ram,omitempty"`           // MB
	DiskGB        *int         `json:"disk,omitempty"`          // GB
	InstantiateOn string       `json:"instantiate_on,omitempty"`
	DiskImage     string       `json:"disk_image,omitempty"`
	Aggregate     string       `json:"aggregate,omitempty"`   // optional
	RoutableIP    *bool        `json:"routable_ip,omitempty"` // optional
	Blockstores   []Blockstore `json:"blockstores,omitempty"`
}

type Blockstore struct {
	Name  string `json:"name"`
	Mount string `json:"mount,omitempty"`
	Size  int    `json:"size"` // GB
}

type Link struct {
	Kind       string   `json:"kind"` // "link" | "lan" | "bridged_link"
	Name       string   `json:"name"`
	Interfaces []Iface  `json:"interfaces"`
	Bandwidth  *int     `json:"bandwidth,omitempty"` // Mbps
	Latency    *int     `json:"latency,omitempty"`   // ms
	Plr        *float64 `json:"plr,omitempty"`       // 0..1
}

type Iface struct {
	Node   string `json:"node"`
	IfName string `json:"ifname,omitempty"`
}
