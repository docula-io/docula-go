package initialize

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
)

type defaultSurvey struct{}

var questions = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "What should we name this dir?"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name: "index",
		Prompt: &survey.Select{
			Message: "Choose an index type",
			Options: []string{"timestamp", "sequential"},
			Default: "timestamp",
		},
	},
}

func (s *defaultSurvey) Ask(opts ...survey.AskOpt) (Configuration, error) {
	var answers Configuration

	if err := survey.Ask(questions, &answers, opts...); err != nil {
		return answers, fmt.Errorf("asking survey: %w", err)
	}

	return answers, nil
}
