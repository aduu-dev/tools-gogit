package gogitcmd

import (
	"os"

	"github.com/spf13/cobra"

	"aduu.dev/tools/gogit"
)

// GogitReplaceCMD replaces the local go.mod with one containing no go.mod files.
func GogitReplaceCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace <path>",
		Short: "replaces the local go.mod with one containing no go.mod files",
		Args:  cobra.ExactArgs(1),
	}

	undo := cmd.Flags().Bool("undo", false, "undoes a prior replace on the path")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		if *undo {
			return gogit.UndoRemovingLocalReplacesFromGomod(args[0])
		}

		return gogit.RemoveLocalReplacesFromGomod(args[0])
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand()

	return cmd
}
