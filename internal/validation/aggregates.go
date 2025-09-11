package validation

var Aggregates = map[string]string{
	"urn:publicid:IDN+emulab.net+authority+cm":          "emulab.net",
	"urn:publicid:IDN+utah.cloudlab.us+authority+cm":    "utah.cloudlab.us",
	"urn:publicid:IDN+clemson.cloudlab.us+authority+cm": "clemson.cloudlab.us",
	"urn:publicid:IDN+wisc.cloudlab.us+authority+cm":    "wisc.cloudlab.us",
	"urn:publicid:IDN+apt.emulab.net+authority+cm":      "apt.emulab.net",
}

func IsValidAggregate(urn string) bool {
	if urn == "" {
		return true
	}
	_, ok := Aggregates[urn]
	return ok
}
