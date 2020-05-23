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

	baseCommand := cmd.Flags().String("base-command", "", "sets the base command to use for fixing go.mod: default=gogit. Can also be set via $GOGIT_REPLACE_CMD")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		baseCMD := *baseCommand
		if len(baseCMD) == 0 {
			baseCMD = os.Getenv("GOGIT_REPLACE_CMD")
		}
		if len(baseCMD) == 0 {
			baseCMD = "gogit"
		}

		return gogitinstall.InstallHooks(args[0], baseCMD)
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand()

	return cmd
}
