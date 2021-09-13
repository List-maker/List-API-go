package database

import "fmt"

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         uint   `yaml:"port"`
	DatabaseName string `yaml:"databaseName"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

func (dbConfig DatabaseConfig) String() string {
	var final string
	final += fmt.Sprintf("\thost: %q\n", dbConfig.Host)
	final += fmt.Sprintf("\tport: %d\n", dbConfig.Port)
	final += fmt.Sprintf("\tdatabaseName: %q\n", dbConfig.DatabaseName)
	final += fmt.Sprintf("\tusername: %q\n", dbConfig.Username)
	final += "\tpassword: *******\n"
	return final
}
