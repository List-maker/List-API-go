package main

import (
	"fmt"
	"listes_back/src/database"
	"strconv"
)

type ServerConfig struct {
	Port             uint   `yaml:"port"`
	AvatarsDirectory string `yaml:"avatarsDirectory"`
}

func (srvConfig ServerConfig) String() string {
	var final string
	final += fmt.Sprintf("\tport: %d\n", srvConfig.Port)
	final += fmt.Sprintf("\tavatarsDirectory: %s\n", srvConfig.AvatarsDirectory)
	return final
}

func (srvConfig ServerConfig) GetStringAddress() string {
	return ":" + strconv.FormatUint(uint64(srvConfig.Port), 10)
}

type Config struct {
	Database database.DatabaseConfig `yaml:"database"`
	Server   ServerConfig            `yaml:"server"`
}

func (config Config) String() string {
	var final string
	final += "database:\n"
	final += config.Database.String()
	final += "server:\n"
	final += config.Server.String()
	return final
}
