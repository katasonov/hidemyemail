package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Host    string
	Port	string
	DbConnectionString string
}

var g_config Configuration

func LoadConfig() error {
	file, err := os.Open("hidemyemail.cfg")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	g_config = Configuration{}
	err = decoder.Decode(&g_config)
	if err != nil {
		return err
	}
	return err
}
