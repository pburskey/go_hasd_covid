package main

import (
	"github.com/cucumber/godog"
)

func iConsumeTheUrl() error {
	return godog.ErrPending
}

func theHASDHasACovidSite() error {
	return godog.ErrPending
}

func theHttpStatusShouldBe(arg1 string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {

}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(
		func() {
			// clean the state before every suite
		})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		// clean the state before every scenario
	})

	ctx.Step(`^I consume the url$`, iConsumeTheUrl)
	ctx.Step(`^the HASD has a covid site$`, theHASDHasACovidSite)
	ctx.Step(`^the http status should be "([^"]*)"$`, theHttpStatusShouldBe)
}
