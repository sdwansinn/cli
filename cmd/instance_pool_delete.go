package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

var instancePoolDeleteCmd = &cobra.Command{
	Use:     "delete <name | id>+",
	Short:   "Delete an instance pool",
	Aliases: gDeleteAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		zoneflag, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}

		zone, err := getZoneByName(zoneflag)
		if err != nil {
			return err
		}

		tasks := make([]task, 0, len(args))
		for _, arg := range args {
			if !force {
				if !askQuestion(fmt.Sprintf("sure you want to delete %q", arg)) {
					continue
				}
			}

			i, err := getInstancePoolByName(arg, zone.ID)
			if err != nil {
				return err
			}

			tasks = append(tasks, task{
				egoscale.DestroyInstancePool{
					ID:     i.ID,
					ZoneID: zone.ID,
				},
				fmt.Sprintf("Deleting instance pool %q", args[0]),
			})
		}

		r := asyncTasks(tasks)
		errs := filterErrors(r)
		if len(errs) > 0 {
			return errs[0]
		}

		return nil
	},
}

func init() {
	// Required Flags
	instancePoolDeleteCmd.Flags().StringP("zone", "z", "", "Instance pool zone")
	if err := instancePoolDeleteCmd.MarkFlagRequired("zone"); err != nil {
		log.Fatal(err)
	}

	instancePoolDeleteCmd.Flags().BoolP("force", "f", false, "Attempt to remove instance pool without prompting for confirmation")
	instancePoolCmd.AddCommand(instancePoolDeleteCmd)
}
