package main

import (
	"os"

	"cosmossdk.io/log"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/skip-mev/connect/v2/tests/simapp"
	cmd "github.com/skip-mev/connect/v2/tests/simapp/connectd/testappd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", simapp.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
