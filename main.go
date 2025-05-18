package main

import (
	"embed"
	"enuma-elish/api"
	"enuma-elish/config"
	"enuma-elish/infra"
	"enuma-elish/pkg/migration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:embed migration/*
var MigrationFS embed.FS

//go:embed seeder/*
var SeederFS embed.FS

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

	m := migration.New(i.Postgres, MigrationFS, SeederFS)
	cmd.AddCommand(m.Command())

	a := api.New(c, i)
	cmd.AddCommand(a.Command())

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command execution failed")
		panic(err)
	}
}
