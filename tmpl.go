package main

import "fmt"
import "io"
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

type TmplProcessingError struct {
    Name string
    Action string
    Source string
    Original error
}

func (e *TmplProcessingError) Error() string {
    return "Failed to " + e.Action +
        " template '" + e.Name +"' (" + e.Source + "):" +
        "\n" + e.Original.Error()
}

func TmplOutput(w io.Writer, config *GLConfig, name string) (error) {
    conf_tmpl := config.Templates[name]
    source := conf_tmpl.Source

    if source == "core.json" {
        // Output with special JSON dumper
        err := ConfigDump(config, w)
        if err != nil { 
            return &TmplProcessingError{name,"execute",source,err}
        }
    } else {
        // Output with normal template
        
        tmpl_obj, err := template.ParseFiles(source)
        if err != nil { 
            return &TmplProcessingError{name,"parse",source,err}
        }

        err = tmpl_obj.Execute(w, config)
        if err != nil { 
            return &TmplProcessingError{name,"execute",source,err}
        }
    }

    return nil
}

func TmplPrint(config *GLConfig, names []string) (error) {
    headers := (len(names) > 1)

    for _, name := range names {
        if headers {
            fmt.Printf("\n==== %s ====\n", name)
        }

        err := TmplOutput(os.Stdout, config, name)
        if err != nil { return err }
    }
    return nil
}

func TmplSet(config *GLConfig, names []string) (error) {
    for _, name := range names {
        conf_tmpl   := config.Templates[name]
        source      := conf_tmpl.Source
        output_path := conf_tmpl.Output

        output_file, err := os.Create(output_path)
        if err != nil {
            return &TmplProcessingError{name,"open output of",source,err}
        }
        defer func() {
            err := output_file.Close(); if err != nil { panic(err) }
        }

        err = TmplOutput(output_file, config, name)
        if err != nil { return err }
    }
    return nil
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

    ConfigSanitize(&config)
    if directive == "print" {
        err = TmplPrint(&config, templates)
    } else {
        err = TmplSet(&config, templates)
    }
    if err != nil {
        fmt.Fprintf(os.Stderr,"goldilocks: %v\n",err)
    }
}
