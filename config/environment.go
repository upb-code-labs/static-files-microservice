package config

import "github.com/kelseyhightower/envconfig"

type Environment struct {
	ArchivesVolumePath  string `split_words:"true" default:"./files/archives"`
	TemplatesVolumePath string `split_words:"true" default:"./files/templates"`
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
