package main

import "encoding/json"

type VPC struct {
	Name  string `json:"name"`
	VpcID string `json:"vpcID"`
}

// AssetVPC a slice of vpc's
type AssetVPC []*VPC

// convert structure to json string
func (v VPC) String() string {
	j, _ := json.Marshal(v)
	return string(j)
}

// Get check struct by id  and returns it
// returns nil if it doesn't exist
func (a AssetVPC) Get(id string) *VPC {
	for _, v := range a {
		if v.VpcID == id {
			return v
		}
	}
	return nil
}

// LookupByName returns struct by name if exists,otherwise nil
func (a AssetVPC) LookupByName(name string) func() interface{} {
	return func() interface{} {
		for _, v := range a {
			if v.Name == name {
				return v
			}
		}
		return nil
	}
}
