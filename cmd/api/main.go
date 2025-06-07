package main

import (
	"enuma-elish/api"
	"enuma-elish/config"
	"enuma-elish/infra"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	c, err := config.New("config.json")
	if err != nil {
		panic(err)
	}

	i, err := infra.New(c)
	if err != nil {
		log.Err(err).Msg("failed to initialize infra")
		panic(err)
	}
	cmd := &cobra.Command{
		Use: "enuma-elish",
	}

	a := api.New(c, i)
	cmd.AddCommand(a.Command())

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command execution failed")
		panic(err)
	}
}
