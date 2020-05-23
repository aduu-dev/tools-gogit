package gogitcmd

import (
	"os"

	"github.com/spf13/cobra"

	"aduu.dev/tools/gogit/gogitinstall"
)

// gogitInstallHooksCMD installs pre-commit and post-commit hooks to temporarily chang gogit.
func GogitInstallHooksCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install-hooks <repo>",
		Short: "installs pre-commit and post-commit hooks which remove local go.mod directives",
		Args:  cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return gogitinstall.InstallHooks(args[0])
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand()

	return cmd
}
