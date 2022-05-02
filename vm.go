package main

import "encoding/json"

type VM struct {
	Name                string   `json:"name"`
	NetworkInterfaceIDs []string `json:"networkInterfaceIDs"`
	SecurityGroupIDs    []string `json:"securityGroupIDs"`
	VpcID               string   `json:"vpcID"`
}

type AssetVM []*VM

func (vm VM) String() string {
	j, _ := json.Marshal(vm)
	return string(j)
}

func (a AssetVM) Get(name string) *VM {
	for _, v := range a {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// HasNetwork check if the vm has a corresponding network id
// in the network asset
func (vm *VM) HasNetwork(asset AssetNI) ([]string, bool) {
	var result []string
	nets := vm.NetworkInterfaceIDs
	for _, v := range nets {
		n := asset.Get(v)
		if n != nil {
			result = append(result, n.Name)
		}
	}
	if len(result) > 0 {
		return result, true
	}
	return nil, false
}

// HasSecurityGroup check if vm has a corresponding security group id
// in the security group asset
func (vm *VM) HasSecurityGroup(asset AssetSG) ([]string, bool) {
	var result []string
	sgs := vm.SecurityGroupIDs
	for _, v := range sgs {
		n := asset.Get(v)
		if n != nil {
			result = append(result, n.Name)
		}
	}
	if len(result) > 0 {
		return result, true
	}
	return nil, false
}

// HasVPC check if vm has a corresponding vpc id
// in the vpc asset
func (vm *VM) HasVPC(asset AssetVPC) (string, bool) {
	vpc := vm.VpcID
	n := asset.Get(vpc)
	if n != nil {
		return n.Name, true
	}

	return "", false
}

// Get returns vm name
func (vm *VM) Get() string {
	return vm.Name
}

// LookupByName return vm object by name
func (a AssetVM) LookupByName(name string) func() interface{} {
	return func() interface{} {
		for _, v := range a {
			if v.Name == name {
				return v
			}
		}
		return nil
	}
}

// LookupBySG returns the vm names that have the specified
//security group ids
func (a AssetVM) LookupBySG(sgID ...string) ([]string, bool) {
	var result []string
	for _, v := range a {
		for _, s := range sgID {
			if found := checkSlice(v.SecurityGroupIDs, s); found {
				result = append(result, v.Name)
			}
		}
	}

	if len(result) > 0 {
		return result, true
	}

	return nil, false
}

// check if string is found in slice
func checkSlice(in []string, x string) bool {
	for _, v := range in {
		if v == x {
			return true
		}
	}
	return false
}
