package main

import (
	"bufio"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config | config.yaml  file struct
type ConfigYaml struct {
	Main struct {
		Amplify        int    `yaml:"amplify"`
		Token          string `yaml:"token"`
		Webhook        string `yaml:"webhook"`
		GuildID        int    `yaml:"guildid"`
		SocketUsage    bool   `yaml:"usesockets"`
		SocketChannels int    `yaml:"socketchannels"`
	} `yaml:"main"`
}

// Load Yaml config
func LoadYaml() {
	cf, err := os.Open("config.yaml")
	if err != nil {
		LogFatal(fmt.Sprintf("error while reading config file ensure that the file exists and is properly formatted, %s", err.Error()))
	}
	defer cf.Close()

	decoder := yaml.NewDecoder(cf)
	err = decoder.Decode(&Config)
	if err != nil {
		LogFatal(fmt.Sprintf("error while reading config file ensure that the file exists and is properly formatted, %s", err.Error()))
	}
}

// Setup the sniper
func SetupSniper() {
	// Loads the proxies into the ProxyChannel
	PrintLogo()

	// Read proxies from file
	proxyFile, err := os.Open("proxies.txt")
	if err != nil {
		LogFatal("no proxies.txt file found! Ensure that the proxies.txt file exists in the current directory")
	}

	scanner := bufio.NewScanner(proxyFile)

	for scanner.Scan() {
		ProxyList = append(ProxyList, scanner.Text())

	}

	// Read vanities from file
	vanityFile, err := os.Open("vanities.txt")
	if err != nil {
		LogFatal("no vanities.txt file found! Ensure that the vanities.txt file exists in the current directory")
	}

	scanner = bufio.NewScanner(vanityFile)

	for scanner.Scan() {
		VanityList = append(VanityList, scanner.Text())

	}
	LoadYaml()

	LogInfo(fmt.Sprintf("initiated sniper with %d proxies and %d vanities", len(ProxyList), len(VanityList)))
	LogInfo("ensure you have enough proxies or this sniper will not work")
	UserInput("enter any key to start the sniper")

}
