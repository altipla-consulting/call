package main

import (
	"os"
	"os/exec"

	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

var cmdUpdate = &cobra.Command{
	Use: "update",
	RunE: func(cmd *cobra.Command, args []string) error {
		install := exec.Command("go", "install", "-v", "github.com/altipla-consulting/call@latest")
		install.Stdin = os.Stdin
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			return errors.Trace(err)
		}

		emoji.Sprint(":heavy_check_mark: CLI call updated successfully!")

		return nil
	},
}
