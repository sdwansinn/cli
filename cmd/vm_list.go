package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var vmlistCmd = &cobra.Command{
	Use:     "list",
	Short:   "List virtual machines instances",
	Aliases: gListAlias,
	Run: func(cmd *cobra.Command, args []string) {
		if err := listVMs(); err != nil {
			log.Fatal(err)
		}
	},
}

func listVMs() error {
	vms, err := cs.List(&egoscale.VirtualMachine{})
	if err != nil {
		return err
	}

	table := table.NewTable(os.Stdout)
	table.SetHeader([]string{"Name", "Security Group", "IP Address", "Status", "Zone", "ID"})

	for _, key := range vms {
		vm := key.(*egoscale.VirtualMachine)

		sgs := getSecurityGroup(vm)

		sgName := strings.Join(sgs, " - ")
		table.Append([]string{vm.Name, sgName, vm.IP().String(), vm.State, vm.ZoneName, vm.ID})
	}
	table.Render()

	return nil
}

func init() {
	vmCmd.AddCommand(vmlistCmd)
}