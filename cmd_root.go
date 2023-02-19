package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{
	Use:          "call",
	SilenceUsage: true,
}

func init() {
	var flagDebug bool
	CmdRoot.PersistentFlags().BoolVar(&flagDebug, "debug", false, "Enable debug logging for this tool")

	CmdRoot.AddCommand(cmdAPI)
	CmdRoot.AddCommand(cmdInstall)
	CmdRoot.AddCommand(cmdUpdate)
	CmdRoot.AddCommand(cmdBackend)

	CmdRoot.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if flagDebug {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}
}
