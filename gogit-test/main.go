// Package contains exploration code to find out how go-git assigns StatusCodes to file.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// RootCMD creates a root command of a program.
func RootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gogit-test",
		Short: "lists for each local file the StatusCode",
		Args:  cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		r, err := git.PlainOpen(args[0])
		if err != nil {
			return
		}

		w, err := r.Worktree()
		if err != nil {
			return
		}

		status, err := w.Status()
		if err != nil {
			return
		}

		files, err := ioutil.ReadDir(args[0])
		if err != nil {
			return
		}

		for _, file := range files {
			fileStatus := status.File(file.Name())

			fmt.Println("Name:", file.Name(), "Staus Worktree:", string(rune(fileStatus.Worktree)), "Status Staging:", string(rune(fileStatus.Staging)),
				"Untracked:", status.IsUntracked(file.Name()))
		}

		return nil
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.AddCommand()
	return cmd
}

func main() {

	errorExitCode := 1

	if err := RootCMD().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(errorExitCode)
	}
}
