package utils

import "flag"

type Arguments struct {
	ConfigPath string
}

func ParseArguments() *Arguments {
	args := &Arguments{}

	flag.StringVar(&args.ConfigPath, "config", ".env", "Path to config file")
	flag.Parse()

	return args
}
