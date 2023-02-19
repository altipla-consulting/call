package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

const installLine = `. <(call completion bash)`

var cmdInstall = &cobra.Command{
	Use: "install",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.Trace(err)
		}
		content, err := ioutil.ReadFile(filepath.Join(home, ".bashrc"))
		if err != nil {
			return errors.Trace(err)
		}
		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			if line == installLine {
				return nil
			}
		}

		f, err := os.OpenFile(filepath.Join(home, ".bashrc"), os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return errors.Trace(err)
		}
		defer f.Close()

		fmt.Fprintln(f)
		fmt.Fprintln(f, installLine)

		log.Info("CLI autocomplete is now installed in ~/.bashrc")
		log.Info("Restart your terminal to finish your setup")

		return nil
	},
}
