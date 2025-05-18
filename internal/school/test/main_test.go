package test

import (
	"enuma-elish/api"
	"enuma-elish/config"
	"enuma-elish/infra"
	"log"
	"os"
	"testing"
)

var (
	testInfra  *infra.Infra
	testConfig *config.Config
	testApi    *api.API
)

func TestMain(m *testing.M) {
	c, err := config.New("../../../config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	testConfig = c

	i, err := infra.New(testConfig)
	if err != nil {
		log.Fatalf("failed to init infra: %v", err)
	}
	testInfra = i

	testApi = api.New(testConfig, testInfra)

	code := m.Run()
	os.Exit(code)
}
