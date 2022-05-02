package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/palantir/stacktrace"
)

func openFile(path string) ([]byte, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}
	return f, nil
}

// NewAsset open the input file
func NewAsset[T AssetVM | AssetNI | AssetSG | AssetVPC](path string) (*T, error) {
	data, err := openFile(path)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	var asset T
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	return &asset, nil
}
