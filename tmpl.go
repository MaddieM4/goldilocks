package main

import "fmt"

type Tmpl struct {
}

const TMPL_DESC string = `
Goldilocks provides facilities for outputting stats
and information to static files, which you can then
serve directly from a webserver, or further process
however you want.

See README for more information on how to configure
template stuff. Template names, sources, and destinations
are all defined in the config file.

goldilocks tmpl print

    Print all templates in config, with headers.

goldilocks tmpl print tmplname

    Print specific template, without header.

goldilocks tmpl set

    Save all templates based on config.

goldilocks tmpl set tmplname

    Save specific template based on config.
`

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
