package main

import (
	"testing"

	"github.com/openshift-pipelines/tektoncd-catalog/acceptance-tests/controller"

	"github.com/cucumber/godog"
)

func TestFeature(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: controller.InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
