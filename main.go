package main

import "fmt"

func getSubcommands() []string {
    return []string{
        "cron ( set | print ) [ name ]",
        "tmpl ( set | print ) [ name ]",
        "help COMMAND",
    }
}

func usage() {
    fmt.Printf("usage: goldilocks COMMAND [ args ]\n\n")
    fmt.Printf("Known subcommands are:\n")
    for _, value := range getSubcommands() {
        fmt.Printf("    %s\n", value)
    }
}

func main() {
    usage()
}
