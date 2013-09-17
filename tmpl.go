package main

import "fmt"

type Tmpl struct {
}

const TMPL_DESC string = ""

func (cron Tmpl) GetName() string {
    return "tmpl"
}

func (tmpl Tmpl) GetUsage() string {
    return tmpl.GetName() + " ( set | print ) [ name ]"
}

func (tmpl Tmpl) GetDescription() string {
    return TMPL_DESC
}

func (tmpl Tmpl) Run(args []string) {
    fmt.Printf("Running tmpl...\n")
}
