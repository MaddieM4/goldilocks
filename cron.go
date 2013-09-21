package main

import "fmt"

type Cron struct {
}

const CRON_DESC string = `
Goldilocks does not run consistently as a daemon.
It can be driven manually or by cron, and includes
the ability to generate and store cron configs based
on its own config, which will wake it for specific
tasks when needed.

goldilocks cron print

    Print the full cron configuration.

goldilocks cron print somename

    Print cron line for just one scheduled action.

goldilocks cron set

    Set full cron config in OS.

goldilocks cron set somename

    Set a limited cron config with just the scheduled action.
`

func (cron Cron) GetName() string {
	return "cron"
}

func (cron Cron) GetUsage() string {
	return cron.GetName() + " ( set | print ) [ name ]"
}

func (cron Cron) GetDescription() string {
	return CRON_DESC
}

func (cron Cron) Run(args []string) {
	fmt.Printf("Running cron...\n")
}
