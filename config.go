package main

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"log"
)

type Configuration struct {
	Host    string
	Port	string
	DbConnectionString string
	ResourcePath string
}

var g_config Configuration

func LoadConfig() error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(path.Join(dir, "hidemyemail.cfg"))
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
