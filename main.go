package main

import (
	"fmt"
	"github.com/PeernetOfficial/core"
	"github.com/PeernetOfficial/core/webapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
	"log"
	"net/http"
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

// SetDefaults This function to be called only during a
// make install
func SetDefaults() error {
	//Setting current directory to default path
	defaultPath := "./"

	//Setting default paths for the config file
	defaults["NumberOfSlaveNode"] = 1
	defaults["NumberOfRootNode"] = 2
	//defaults["NetworkInterface"] = "wlp0s20f3"
	//defaults["NetworkInterfaceIPV6Index"] = "2"

	//Paths to search for config file
	configPaths = append(configPaths, defaultPath)

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetConfigName(configName)
	viper.SetConfigFile(configFile)
	viper.SetConfigType(configType)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
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
		SetDefaults()
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
		StartPeernet(true, folder, i)
		// Create delay 1 second
		// Calling Sleep method
		//time.Sleep(1 * time.Second)
	}

	for i := 0; i < c.NumberOfSlaveNode; i++ {
		// Start root node
		StartPeernet(false, folder, i)
		// Create delay 1 second
		// Calling Sleep method
		//time.Sleep(1 * time.Second)
	}

}

// StartPeernet Spawn Peernet instance
func StartPeernet(rootNode bool, RunPath string, number int) {

	var folder string
	if rootNode {
		folder = RunPath + "/runroot" + strconv.Itoa(number)
	} else {
		folder = RunPath + "/run" + strconv.Itoa(number)
	}

	// Create runs folder if it does not exist
	if err := os.Mkdir(folder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	RunPath = folder

	var config core.Config

	_, err := core.LoadConfig("/Config.yaml", &config)
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
	config.BlockchainGlobal = RunPath + "/data/blockchain global/"
	config.DataFolder = RunPath + "/data"
	config.BlockchainMain = RunPath + "/data/blockchain main/"
	config.LogFile = RunPath + "/data/log backend.txt"
	config.SearchIndex = RunPath + "/data/search index.txt"

	// if root node send back public key,udp address,webapi
	if rootNode {

	} else {
		config.SeedList = []core.PeerSeed{{
			PublicKey: "",
			Address:   []string{""},
		}}
	}

	// save configuration
	err = core.SaveConfig(RunPath+"/Config.yaml", &config)
	if err != nil {
		log.Fatal(err)
	}

	// Copy config to the appropriate folder
	backend, status, err := core.Init("Test framework/1.0", RunPath+"/Config.yaml", nil, nil)
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

	webapi.Start(backend, []string{APIportAddress}, false, "", "", 0, 0, uuid.Nil)

	fmt.Println(APIportAddress)

}

// Spawn test Peernet instances for testing
func main() {
	//1. Create Config function to read configurations
	//2. Spawn root nodes with Custom configs with a folder structure
	// with UUID
	// Test Run 1 ( [/Run1/RootNode1/<config file>],[/Run1/SlaveNode1/<binary>|<config file>]
	// Start add nodes and manage go routine

	//finish := make(chan bool)

	r := mux.NewRouter()

	//// Get Config information
	config, err := ConfigInit()
	if err != nil {
		fmt.Println(err)
	}

	// Create Peernet instances based on the config
	//go func() {
	//}()

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	config.RunManager()

	// go func() {
	log.Fatal(srv.ListenAndServe())
	//}()

	//<-finish
}
