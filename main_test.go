package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBehavior(t *testing.T) {
	t.Run("check expose", func(t *testing.T) {
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"check", "expose"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
		gotOutput := output.String()
		wantOutput := "\t- VM_1\n\t- VM_2\n"
		assert.Equal(t, wantOutput, gotOutput, "graphsec check expose")
	})

	t.Run("check path", func(t *testing.T) {
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"check", "vm", "-n", "VM_1"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		gotOutput := output.String()
		wantOutput := "VM_1 -> NetworkInterface_1 -> NetworkInterface_2 -> SecurityGroup_1 -> VPC_1\n"
		assert.Equal(t, wantOutput, gotOutput, "")
	})

	t.Run("check http port", func(t *testing.T) {
		cmd := NewRootCommand()
		output := &bytes.Buffer{}
		cmd.SetOut(output)
		cmd.SetArgs([]string{"check", "port"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}

		gotOutput := output.String()
		wantOutput := "\t- VM_1\n"
		assert.Equal(t, wantOutput, gotOutput, "")
	})
}
