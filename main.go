package main

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"github.com/spf13/cobra"
	"io"
	"log"
)

const (
	VMFILE  = "./testdata/VM.json"
	SGFILE  = "./testdata/SecurityGroup.json"
	NIFILE  = "./testdata/NetworkInterface.json"
	VPCFILE = "./testdata/VPC.json"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "graphsec",
		Long: "Root command",
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "check options",
	}

	vmCmd := &cobra.Command{
		Use:   "vm",
		Short: "parse vm config",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			name, _ := cmd.Flags().GetString("name")
			full, _ := cmd.Flags().GetBool("full")
			debug, _ := cmd.Flags().GetBool("debug")
			if err := parseVM(name, full, out); err != nil {
				if debug {
					log.Fatal(err)
				} else {
					log.Fatal(stacktrace.RootCause(err))
				}
			}
		},
	}

	exposeCmd := &cobra.Command{
		Use:   "expose",
		Short: "find vms that are exposed to the internet",
		Run: func(cmd *cobra.Command, args []string) {
			output := cmd.OutOrStdout()
			debug, _ := cmd.Flags().GetBool("debug")
			if err := checkExposed(output); err != nil {
				if debug {
					log.Fatal(err)
				} else {
					log.Fatal(stacktrace.RootCause(err))
				}
			}
		},
	}

	portCmd := &cobra.Command{
		Use:   "port",
		Short: "find vms with http port open",
		Run: func(cmd *cobra.Command, args []string) {
			output := cmd.OutOrStdout()
			port := 80
			debug, _ := cmd.Flags().GetBool("debug")
			if err := checkPort(port, output); err != nil {
				if debug {
					log.Fatal(err)
				} else {
					log.Fatal(stacktrace.RootCause(err))
				}
			}
		},
	}

	rootCmd.AddCommand(checkCmd)
	checkCmd.AddCommand(exposeCmd)
	checkCmd.AddCommand(portCmd)
	checkCmd.AddCommand(vmCmd)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug")

	vmCmd.PersistentFlags().StringP("name", "n", "", "vm name")
	vmCmd.PersistentFlags().BoolP("full", "f", false, "full resources output")
	_ = vmCmd.MarkPersistentFlagRequired("name")

	return rootCmd
}

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func parseVM(name string, fullOutput bool, output io.Writer) error {
	store, err := NewStore()
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	vm := store.AssetVM.Get(name)
	if vm == nil {
		return stacktrace.NewError("%s not found", name)
	}

	rel := NewRelationship()
	rel.Add(Relationship{Name: vm.Name, get: store.AssetVM.LookupByName(vm.Name)})

	if nets, ok := vm.HasNetwork(store.AssetNI); ok {
		for _, v := range nets {
			rel.Add(Relationship{Name: v, get: store.AssetNI.LookupByName(v)})
		}
	}

	if sgs, ok := vm.HasSecurityGroup(store.AssetSG); ok {
		for _, v := range sgs {
			rel.Add(Relationship{Name: v, get: store.AssetSG.LookupByName(v)})
		}
	}

	if vpc, ok := vm.HasVPC(store.AssetVPC); ok {
		rel.Add(Relationship{Name: vpc, get: store.AssetVPC.LookupByName(vpc)})
	}

	rel.Traverse(fullOutput, output)
	return nil
}

func checkExposed(output io.Writer) error {
	store, err := NewStore()
	if err != nil {
		return stacktrace.Propagate(err, "")
	}
	sgs, found, err := store.AssetSG.IsExposed()
	if err != nil {
		return stacktrace.Propagate(err, "")
	}
	if found {
		if vms, ok := store.AssetVM.LookupBySG(sgs...); ok {
			for _, v := range vms {
				fprintf(output, "\t- %s\n", v)
			}
		}
	}

	return nil
}

func checkPort(port int, output io.Writer) error {
	store, err := NewStore()
	if err != nil {
		return stacktrace.Propagate(err, "")
	}
	if sgs, found := store.AssetSG.GetOpenPort(port); found {
		if vms, ok := store.AssetVM.LookupBySG(sgs...); ok {
			for _, v := range vms {
				fprintf(output, "\t- %s\n", v)
			}
		} else {
			fprintf(output, "no vms found with port %d open", port)
		}
	} else {
		fprintf(output, "port %d not found", port)
	}
	return nil
}

func fprintf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}
