package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdAPI = &cobra.Command{
	Use:               "api",
	Short:             "Call a method of https://api.dev.remote",
	SilenceUsage:      true,
	Example:           "call api foo.bar.FooService/BarMethod param=foo other=bar My-Header:value",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: sendValidArgs,
}

func init() {
	cmdAPI.RunE = func(cmd *cobra.Command, args []string) error {
		return errors.Trace(sendRequest(cmd.Context(), "api.dev.remote", args))
	}
}
