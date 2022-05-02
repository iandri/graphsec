package main

import (
	"encoding/json"
	"net"

	"github.com/palantir/stacktrace"
)

type SG struct {
	Name         string   `json:"name"`
	GroupID      string   `json:"groupID"`
	VpcID        string   `json:"vpcID"`
	ExposedPorts []int    `json:"exposedPorts"`
	Direction    string   `json:"direction"`
	IPList       []string `json:"ipList"`
}

type AssetSG []*SG

// String convert struct to json string
func (sg SG) String() string {
	j, _ := json.Marshal(sg)
	return string(j)
}

// Get return object by id
func (asset AssetSG) Get(id string) *SG {
	for _, v := range asset {
		if v.GroupID == id {
			return v
		}
	}
	return nil
}

// LookupByName returns object by name
func (asset AssetSG) LookupByName(name string) func() interface{} {
	return func() interface{} {
		for _, v := range asset {
			if v.Name == name {
				return v
			}
		}
		return nil
	}
}

// IsExposed checks if this security has a public ip
func (sg SG) IsExposed() (bool, error) {
	return isPublicIPRange(sg.IPList...)
}

// IsExposed checks if any security group has a public ip
func (asset AssetSG) IsExposed() ([]string, bool, error) {
	var exposed []string
	for _, v := range asset {
		ok, err := v.IsExposed()
		if err != nil {
			return nil, false, stacktrace.Propagate(err, "")
		}
		if ok {
			exposed = append(exposed, v.GroupID)
		}
	}
	if len(exposed) > 0 {
		return exposed, true, nil
	}

	return nil, false, nil
}

// GetOpenPort returns a list of security groups ids with
// port open
func (asset AssetSG) GetOpenPort(port int) ([]string, bool) {
	var result []string
	for _, v := range asset {
		if v.isPortOpen(port) {
			result = append(result, v.GroupID)
		}
	}

	if len(result) > 0 {
		return result, true
	}
	return nil, false
}

// IsHTTPOpen returns all security groups ids with
// http port open
func (asset AssetSG) IsHTTPOpen() ([]string, bool) {
	return asset.GetOpenPort(80)
}

// isPortOpen check if security group has port open
func (sg SG) isPortOpen(port int) bool {
	for _, v := range sg.ExposedPorts {
		if sg.Direction != "inbound" {
			continue
		}
		if v == port {
			return true
		}
	}
	return false
}

// isPublicIP check if the ip is public
func isPublicIP(ipaddr string) (bool, error) {
	if ip, _, err := net.ParseCIDR(ipaddr); err == nil {
		return !ip.IsPrivate(), nil
	}

	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return false, stacktrace.NewError("IP Address: %s - Invalid\n", ip)
	}
	return !ip.IsPrivate(), nil
}

// isPublicIPRange check if any ip from the ips specified is public
func isPublicIPRange(ipaddr ...string) (bool, error) {
	for _, ip := range ipaddr {
		isPublic, err := isPublicIP(ip)
		if err != nil {
			return false, stacktrace.Propagate(err, "")
		}
		if isPublic {
			return true, nil
		}
	}
	return false, nil
}
