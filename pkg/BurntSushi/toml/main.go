package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	Age        int
	Cats       []string
	Pi         float64
	Perfection []int
	DOB        time.Time // requires `import time`
	Song       []song
	tomlConfig
}

type song struct {
	Name     string
	Duration duration
}

type duration struct {
	time.Duration
}

type tomlConfig struct {
	Title string
	Owner ownerInfo
	DB database `toml:"database"`
	Servers map[string]server
	Clients clients
}

type ownerInfo struct {
	Name string
	Org string `toml:"organization"`
	Bio string
	DOB time.Time
}

type database struct {
	Server string
	Ports []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data [][]interface{}
	Hosts []string
}

func main() {
	var Conf Config

	configPath := flag.String("f", "./BurntSushi/toml/config/config.toml", "config file")
	flag.Parse()

	tomlDataByte, err := ioutil.ReadFile(*configPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	if _, err := toml.Decode(string(tomlDataByte), &Conf); err != nil {
		fmt.Println(err.Error())
	}
	log.Println(Conf)

	for _, s := range Conf.Song {
		fmt.Printf("%s (%s)\n", s.Name, s.Duration)
	}

}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
