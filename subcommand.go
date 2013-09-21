package main

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
