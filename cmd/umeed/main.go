package main

import (
	"os"
	"strings"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	nebulaapp "github.com/tessornetwork/nebula/v3/app"
	appparams "github.com/tessornetwork/nebula/v3/app/params"
	"github.com/tessornetwork/nebula/v3/cmd/nebud/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, strings.ToUpper(appparams.Name), nebulaapp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
