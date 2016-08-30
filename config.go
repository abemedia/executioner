package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "flag"
    "log"
    "strings"
)

type config struct {
    Host string
    Port int
    Secret string
    LogPath string `yaml:"log_path"`
    Endpoints map[string]string
    CMD map[string][]string
}

var Config config

func init() {
    // get config file path from command line flag
    path := flag.String("c", "./config.yml", "path to config file")
    flag.Parse()

    // open config file
    file, err := ioutil.ReadFile(*path)
    if err != nil {
        log.Fatal(err.Error())
    }

    // parse yaml into struct
    err = yaml.Unmarshal(file, &Config)
    if err != nil {
        log.Fatal(err.Error())
    }

    Config.CMD = make(map[string][]string)
    for key, val := range Config.Endpoints {
        Config.CMD[key] = strings.Split(val, " ")
    }

}
