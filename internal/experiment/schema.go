package experiment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		DeleteContext: resourceDelete,
		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true, ForceNew: true},
			"project":  {Type: schema.TypeString, Optional: true, ForceNew: true},
			"pem_path": {Type: schema.TypeString, Optional: true, ForceNew: true},

			"wait_for_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "provisioned",
				ValidateFunc: validation.StringInSlice([]string{"provisioned", "ready"}, false),
				ForceNew:     true,
			},

			// outputs
			"uuid":    {Type: schema.TypeString, Computed: true},
			"url":     {Type: schema.TypeString, Computed: true},
			"status":  {Type: schema.TypeString, Computed: true},
			"expires": {Type: schema.TypeString, Computed: true},
			"nodes":   {Type: schema.TypeMap, Elem: &schema.Schema{Type: schema.TypeString}, Computed: true},

			"rawpc":        {Type: schema.TypeList, Optional: true, Elem: rawpcBlock(), ForceNew: true},
			"xenvm":        {Type: schema.TypeList, Optional: true, Elem: xenvmBlock(), ForceNew: true},
			"link":         {Type: schema.TypeList, Optional: true, Elem: linkBlock(), ForceNew: true},
			"lan":          {Type: schema.TypeList, Optional: true, Elem: lanBlock(), ForceNew: true},
			"bridged_link": {Type: schema.TypeList, Optional: true, Elem: bridgedLinkBlock(), ForceNew: true},
		},
	}
}

func rawpcBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":          {Type: schema.TypeString, Required: true},
		"hardware_type": {Type: schema.TypeString, Optional: true},
		"exclusive":     {Type: schema.TypeBool, Optional: true},
		"disk_image":    {Type: schema.TypeString, Optional: true},
		"aggregate":     {Type: schema.TypeString, Optional: true}, // optional
		"routable_ip":   {Type: schema.TypeBool, Optional: true},
		"blockstore":    {Type: schema.TypeList, Optional: true, Elem: blockstoreBlock()},
	}}
}

func xenvmBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":           {Type: schema.TypeString, Required: true},
		"cores":          {Type: schema.TypeInt, Optional: true},
		"ram_mb":         {Type: schema.TypeInt, Optional: true},
		"disk_gb":        {Type: schema.TypeInt, Optional: true},
		"instantiate_on": {Type: schema.TypeString, Optional: true},
		"disk_image":     {Type: schema.TypeString, Optional: true},
		"aggregate":      {Type: schema.TypeString, Optional: true}, // optional
		"routable_ip":    {Type: schema.TypeBool, Optional: true},
		"blockstore":     {Type: schema.TypeList, Optional: true, Elem: blockstoreBlock()},
	}}
}

func blockstoreBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":    {Type: schema.TypeString, Required: true},
		"mount":   {Type: schema.TypeString, Optional: true},
		"size_gb": {Type: schema.TypeInt, Required: true},
	}}
}

func ifaceBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"node":   {Type: schema.TypeString, Required: true},
		"ifname": {Type: schema.TypeString, Optional: true},
	}}
}

func linkBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":      {Type: schema.TypeString, Required: true},
		"interface": {Type: schema.TypeList, Required: true, MinItems: 2, MaxItems: 2, Elem: ifaceBlock()},
	}}
}

func lanBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":      {Type: schema.TypeString, Required: true},
		"interface": {Type: schema.TypeList, Required: true, MinItems: 2, Elem: ifaceBlock()},
	}}
}

func bridgedLinkBlock() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"name":           {Type: schema.TypeString, Required: true},
		"bandwidth_mbps": {Type: schema.TypeInt, Optional: true},
		"latency_ms":     {Type: schema.TypeInt, Optional: true},
		"plr":            {Type: schema.TypeFloat, Optional: true},
		"interface":      {Type: schema.TypeList, Required: true, MinItems: 2, Elem: ifaceBlock()},
	}}
}
