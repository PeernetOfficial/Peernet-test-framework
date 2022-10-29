package main

import (
	"fmt"
	"github.com/PeernetOfficial/core"
	"github.com/PeernetOfficial/core/webapi"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"time"
)

// Default variables
var (
	defaultPath    string
	defaults       = map[string]interface{}{}
	configName     = "config"
	configType     = "json"
	configFile     = "config.json"
	configPaths    []string
	defaultEnvName = "Peernet-Test-Framework"
)

// Config TestNetwork config file
type Config struct {
	NumberOfSlaveNode int
	NumberOfRootNode  int
}

// Spawn test Peernet instances for testing
func main() {
	//1. Create Config function to read configurations
	//2. Spawn root nodes with Custom configs with a folder structure
	// with UUID
	// Test Run 1 ( [/Run1/RootNode1/<config file>],[/Run1/SlaveNode1/<binary>|<config file>]
	// Start add nodes and manage go routine

	// Get Config information
	config, err := ConfigInit()
	if err != nil {
		fmt.Println(err)
	}
	// Create Peernet instances based on the config
	config.RunManager()
}

// ConfigInit Reads config file for settings
// for the test network
func ConfigInit() (*Config, error) {
	//Paths to search for config file
	configPaths = append(configPaths, ".")

	//Add all possible configurations paths
	for _, v := range configPaths {
		viper.AddConfigPath(v)
	}

	//Read config file
	if err := viper.ReadInConfig(); err != nil {
		// If the error thrown is config file not found
		//Sets default configuration to viper
		for k, v := range defaults {
			viper.SetDefault(k, v)
		}
		viper.SetConfigName(configName)
		viper.SetConfigFile(configFile)
		viper.SetConfigType(configType)

		if err = viper.WriteConfig(); err != nil {
			return nil, err
		}
	}

	// Adds configuration to the struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Stores all root nodes generated
type RootNodes struct {
	RootNode []RootNode
}

// RootNode Store information regarding root nodes
type RootNode struct {
	PublicKey     string
	UDPAddress    string
	WebAPIAddress string
}

// RunManager Starts Peernet test network for testing
func (c *Config) RunManager() {
	// Create runs folder if it does not exist
	os.Mkdir("runs", os.ModePerm)

	now := time.Now()
	folder := "runs/test" + now.Format("20060102150405")

	// Create runs folder if it does not exist
	if err := os.Mkdir(folder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Create Root Node
	for i := 0; i < c.NumberOfRootNode; i++ {
		// Start root node
		go StartPeernet(true, folder, i)
	}

}

// StartPeernet Spawn Peernet instance
func StartPeernet(rootNode bool, RunPath string, number int) {

	var config core.Config

	_, err := core.LoadConfig(RunPath+"/Config.yaml", &config)
	if err != nil {
		log.Fatal(err)
	}

	// set paths for warehouse
	config.WarehouseMain = RunPath + "/data"

	UDPport, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	UDPportAddress := "127.0.0.1:" + strconv.Itoa(UDPport)
	config.Listen = []string{UDPportAddress}

	// if root node send back public key,udp address,webapi
	if rootNode {

	} else {
		config.SeedList = []core.PeerSeed{{
			PublicKey: "",
			Address:   []string{""},
		}}
	}

	// Copy config to the appropriate folder
	backend, status, err := core.Init("Your application/1.0", RunPath+"/Config.yaml", nil, nil)
	if status != core.ExitSuccess {
		fmt.Printf("Error %d initializing backend: %s\n", status, err.Error())
		return
	}

	backend.Connect()

	APIport, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	APIportAddress := "127.0.0.1:" + strconv.Itoa(APIport)

	webapi.Start(backend, []string{APIportAddress}, false, "", "", 0, 0, nil)

	// Figure which folder to read config from
}
