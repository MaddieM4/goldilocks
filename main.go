package main

import "fmt"
import "flag"

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
			cmnd.Run(args[1:])
		} else {
			fmt.Printf("No command '%s' known\n\n", subcommand)
			usage()
		}
	}
}
