package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdBackend = &cobra.Command{
	Use:               "backend",
	Short:             "Call a method of https://backend.dev.remote",
	SilenceUsage:      true,
	Example:           "call backend foo.bar.FooService/BarMethod param=foo other=bar My-Header:value",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: sendValidArgs,
}

func init() {
	cmdBackend.RunE = func(cmd *cobra.Command, args []string) error {
		return errors.Trace(sendRequest(cmd.Context(), "backend.dev.remote", args))
	}
}
