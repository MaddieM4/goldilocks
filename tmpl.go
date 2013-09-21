package main

import "fmt"
import "os"
import "html/template"

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

func TmplPrintUsage(tmpl *Tmpl) {
    fmt.Fprintf(
        os.Stderr,
        "usage: %s %s\n\n%s\n",
        "goldilocks",
        tmpl.GetUsage(),
        "See `goldilocks help tmpl` for usage.",
    )
}

func TmplPrint(config *GLConfig, names []string) {
    headers := (len(names) > 1)

    for _, name := range names {
        if headers {
            fmt.Printf("\n==== %s ====\n", name)
        }

        conf_tmpl := config.Templates[name]

        tmpl_obj, err := template.ParseFiles(conf_tmpl.Source)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "goldilocks: Failed to parse template '%s' (%s):\n%v\n",
                name,
                conf_tmpl.Source,
                err,
            )
            return
        }

        err = tmpl_obj.Execute(os.Stdout, conf_tmpl)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "goldilocks: Failed to exec template '%s' (%s):\n%v\n",
                name,
                conf_tmpl.Source,
                err,
            )
            return
        }
    }
}

func (tmpl Tmpl) Run(args []string) {
    if len(args) < 1 {
        TmplPrintUsage(&tmpl)
        return
    }

    directive := args[0]
    if directive != "print" && directive != "set" {
        TmplPrintUsage(&tmpl)
        return
    }

    config, err := GetConfig([]string{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "goldilocks: %v\n", err)
        return
    }

    templates := make([]string, 0)
    if len(args) > 1 {
        for _, name := range args[1:] {
            _, ok := config.Templates[name]
            if ! ok {
                fmt.Fprintf(
                    os.Stderr,
                    "goldilocks: No template '%s' in config\n",
                    name,
                )
                return
            }
            templates = append(templates, name)
        }
    } else {
        for name, _ := range config.Templates {
            templates = append(templates, name)
        }
    }

    if directive == "print" {
        TmplPrint(&config, templates)
    }
}
