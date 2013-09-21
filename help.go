package main

import "fmt"

type Help struct {
}

const HELP_DESC string = `
List more detail about a subcommand.

goldilocks help

    Print usage information

goldilocks help cron

    Print specific information about the cron subcommand.

goldilocks help help

    Breaks the internet. Don't do it, Jen!
`

func (help Help) GetName() string {
	return "help"
}

func (help Help) GetUsage() string {
	return help.GetName() + " COMMAND"
}

func (help Help) GetDescription() string {
	return HELP_DESC
}

func doHelp(cmnd GLSubcommand) {
	fmt.Printf(":: %s\n%s\n", cmnd.GetUsage(), cmnd.GetDescription())
}

func (help Help) Run(args []string) {
	if len(args) == 0 {
		usage()
	} else {
		cmnd_map := getSubcommandMap()
		for _, value := range args {
			cmnd, ok := cmnd_map[value]
			if ok {
				doHelp(cmnd)
			} else {
				fmt.Printf("No subcommand '%s'\n", value)
			}
		}
	}
}
