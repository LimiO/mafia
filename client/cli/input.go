package cli

import (
	"fmt"
	"github.com/dixonwille/wlog/v3"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dixonwille/wmenu/v5"
)

func AskSelect(message string, options []string) string {
	menu := wmenu.NewMenu(message)
	var result string
	for _, option := range options {
		menu.Option(option, nil, false, nil)
	}
	menu.Action(func(opts []wmenu.Opt) error {
		fmt.Printf("You selected: %q\n", opts[0].Text)
		result = opts[0].Text
		return nil
	})
	menu.AddColor(wlog.Yellow, wlog.Green, wlog.None, wlog.Red)
	menu.ClearOnMenuRun()
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
