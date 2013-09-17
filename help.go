package main

import "fmt"

type Help struct {
}

const HELP_DESC string = ""

func (help Help) GetName() string {
    return "help"
}

func (help Help) GetUsage() string {
    return help.GetName() + " COMMAND"
}

func (help Help) GetDescription() string {
    return HELP_DESC
}

func (help Help) Run(args []string) {
    fmt.Printf("Running help...\n")
}
