package main

import (
	"io"
	"os"

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
	var flagBody string
	cmdBackend.Flags().StringVarP(&flagBody, "body", "d", "", "Body of the request.")

	cmdBackend.RunE = func(cmd *cobra.Command, args []string) error {
		if flagBody == "-" {
			body, err := io.ReadAll(os.Stdin)
			if err != nil {
				return errors.Trace(err)
			}
			flagBody = string(body)
		}
		return errors.Trace(sendRequest(cmd.Context(), "backend.dev.remote", flagBody, args))
	}
}
