package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func AskSelect(message string, options []string) (string, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	var selected string
	err := survey.AskOne(prompt, &selected, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("failed to select option: %v", err)
	}
	return selected, nil
}

func AskInput(message string) (string, error) {
	answer := ""
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return "", fmt.Errorf("failed to ask input: %v", err)
	}
	return answer, nil
}
