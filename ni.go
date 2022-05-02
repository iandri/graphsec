package main

import "encoding/json"

type NI struct {
	Name               string   `json:"name"`
	NetworkInterfaceID string   `json:"networkInterfaceID"`
	SecurityGroupIDs   []string `json:"securityGroupIDs"`
	VpcID              string   `json:"vpcID"`
}

type AssetNI []*NI

func (n NI) String() string {
	j, _ := json.Marshal(n)
	return string(j)
}

func (a AssetNI) Get(id string) *NI {
	for _, v := range a {
		if v.NetworkInterfaceID == id {
			return v
		}
	}
	return nil
}

func (a AssetNI) LookupByName(name string) func() interface{} {
	return func() interface{} {
		for _, v := range a {
			if v.Name == name {
				return v
			}
		}
		return nil
	}
}
