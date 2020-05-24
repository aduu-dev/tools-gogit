// Package main.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

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
	klog.InitFlags(nil)
	flag.Parse()
	// Make cobra aware of select glog flags
	// Enabling all flags causes unwanted deprecation warnings from glog to always print in plugin mode
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("logtostderr"))
	pflag.CommandLine.Set("logtostderr", "true")
	pflag.CommandLine.Set("v", "5")

	errorExitCode := 1

	if err := RootCMD().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(errorExitCode)
	}
}
