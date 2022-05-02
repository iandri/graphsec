package main

import "github.com/palantir/stacktrace"

type Store struct {
	AssetVM
	AssetNI
	AssetSG
	AssetVPC
}

// NewStore create a new place to store the objects files
func NewStore() (*Store, error) {
	assetvm, err := NewAsset[AssetVM](VMFILE)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	assetni, err := NewAsset[AssetNI](NIFILE)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	assetsg, err := NewAsset[AssetSG](SGFILE)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	assetvpc, err := NewAsset[AssetVPC](VPCFILE)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return &Store{
		*assetvm,
		*assetni,
		*assetsg,
		*assetvpc,
	}, nil
}
