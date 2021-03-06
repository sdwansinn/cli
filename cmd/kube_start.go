package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// kubeStartCmd represents the start command
var kubeStartCmd = &cobra.Command{
	Use:   "start <cluster name>",
	Short: "Start a stopped standalone Kubernetes cluster instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		clusterName := args[0]

		vm, err := getKubeVM(clusterName)
		if err != nil {
			return err
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		if !force {
			if !askQuestion(fmt.Sprintf("sure you want to stop %q cluster instance", vm.Name)) {
				return nil
			}
		}

		if err := startVirtualMachine(vm.Name, ""); err != nil {
			fmt.Println("failed")
			return err
		}

		return nil
	},
}

func init() {
	kubeStartCmd.Flags().BoolP("force", "f", false, "Attempt to start cluster instance without prompting for confirmation")
	kubeCmd.AddCommand(kubeStartCmd)
}
