package main

import "fmt"
import "flag"

type GLSubcommand interface {
    GetName() string
    GetUsage() string
    GetDescription() string

    Run(args []string)
}

func getSubcommands() []GLSubcommand {
    return []GLSubcommand{
        Cron{},
        Tmpl{},
        Help{},
    }
}

func getSubcommandMap() map[string]GLSubcommand {
    cmnd_list := getSubcommands()
    cmnd_map := make(map[string]GLSubcommand, len(cmnd_list))

    for _, value := range cmnd_list {
        cmnd_map[value.GetName()] = value
    }
    return cmnd_map
}

func usage() {
    fmt.Printf("usage: goldilocks COMMAND [ args ]\n\n")
    fmt.Printf("Known subcommands are:\n")
    for _, value := range getSubcommands() {
        fmt.Printf("    %s\n", value.GetUsage())
    }
}

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) == 0 {
        usage()
    } else {
        subcommand := args[0]
        cmnd_map := getSubcommandMap()
        cmnd, ok := cmnd_map[subcommand] 
        if ok {
            cmnd.Run(args[0:])
        } else {
            fmt.Printf("No command '%s' known\n\n", subcommand)
            usage()
        }
    }
}
