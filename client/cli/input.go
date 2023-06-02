package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dixonwille/wmenu/v5"
)

func AskSelect(message string, options []string) string {
	menu := wmenu.NewMenu(message)

	var result string
	menu.Action(func(opts []wmenu.Opt) error {
		fmt.Printf("You selected: %q\n", opts[0].Text)
		result = opts[0].Text
		return nil
	})
	for _, option := range options {
		menu.Option(option, nil, false, nil)
	}
	_ = menu.Run()
	return result
}

func AskInput(message string) string {
	answer := ""
	prompt := &survey.Input{
		Message: message,
	}
	_ = survey.AskOne(prompt, &answer)
	return answer
}
