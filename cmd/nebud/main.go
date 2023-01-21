package main

import (
	"os"
	"strings"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	nebulaapp "github.com/tessornetwork/nebula/app"
	appparams "github.com/tessornetwork/nebula/app/params"
	"github.com/tessornetwork/nebula/cmd/nebud/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, strings.ToUpper(appparams.Name), nebulaapp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
