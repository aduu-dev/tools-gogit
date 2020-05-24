package gogitcmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"aduu.dev/tools/gogit/replace"
)

func gomodFilepath(base string) string {
	return filepath.Join(base, "go.mod")
}

// GogitReplaceCMD replaces the local go.mod with one containing no go.mod files.
func GogitReplaceCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace <path>",
		Short: "replaces the local go.mod with one containing no go.mod files",
		Long: `The command only works on the status of the staged file and not
on the file's status in the working directory itself to avoid doing work on a non-staged go.mod'`,
		Args: cobra.ExactArgs(1),
	}

	undo := cmd.Flags().Bool("undo", false, "undoes a prior replace on the path")
	workOnStaged := cmd.Flags().Bool("replace-only-if-staged", false, "modifies only the staged go.mod if this is set to true")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		if *undo {
			return replace.UndoRemovingLocalReplacesFromGomod(args[0], *workOnStaged)
		}

		return replace.RemoveLocalReplacesFromGomod(args[0], *workOnStaged)
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand()

	return cmd
}
