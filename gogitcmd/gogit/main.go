// Package main.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"aduu.dev/tools/gogit/gogitcmd"
)

// RootCMD creates a root command of a program.
func RootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gogit",
		Short: "offers a couple helpers to remove local replace directives during commit",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		return cmd.Help()
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand(gogitcmd.GogitInstallHooksCMD())
	cmd.AddCommand(gogitcmd.GogitReplaceCMD())
	return cmd
}

func main() {
	kmain

	errorExitCode := 1

	if err := RootCMD().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(errorExitCode)
	}
}
