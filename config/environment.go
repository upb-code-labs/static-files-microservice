package config

import "github.com/kelseyhightower/envconfig"

type Environment struct {
	ArchivesVolumePath string `spilt_words:"true" default:"./archives"`
}

var env *Environment

func GetEnvironment() *Environment {
	if env == nil {
		env = &Environment{}

		err := envconfig.Process("", env)
		if err != nil {
			panic(err)
		}
	}

	return env
}
