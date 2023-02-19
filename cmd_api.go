package main

import (
	"io"
	"os"

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
	var flagBody string
	cmdAPI.Flags().StringVarP(&flagBody, "body", "d", "", "Body of the request.")

	cmdAPI.RunE = func(cmd *cobra.Command, args []string) error {
		if flagBody == "-" {
			body, err := io.ReadAll(os.Stdin)
			if err != nil {
				return errors.Trace(err)
			}
			flagBody = string(body)
		}
		return errors.Trace(sendRequest(cmd.Context(), "api.dev.remote", flagBody, args))
	}
}
