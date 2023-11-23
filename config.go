package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
}

func parseProfile(filepath string) (*Profile, error) {
	prf := &Profile{}
	f, err := os.Open(filepath)

	if err != nil {
		return prf, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return prf, err
	}

	if err := yaml.Unmarshal(data, prf); err != nil {
		return prf, err
	}

	return prf, nil
}

func getProfile(profName string) (*Profile, error) {
	usr, err := user.Current()
	if err != nil {
		log.Printf("error: user.Current(). %s", err)
		os.Exit(1)
	}

	var profile = &Profile{}
	profPath := path.Join(usr.HomeDir, ".config", "discord-quickpost", profName+".yaml")
	profile, err = parseProfile(profPath)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
