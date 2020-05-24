package gogitcmd

import (
	"os"

	"github.com/spf13/cobra"
)

// CMD offers a couple helpers to remove local replace directives during commit.
func CMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gogit",
		Short: "offers a couple helpers to remove local replace directives during commit.",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return cmd.Help()
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(GogitInstallHooksCMD(), GogitRemoveHooksCMD())
	cmd.AddCommand(GogitReplaceCMD())

	return cmd
}
