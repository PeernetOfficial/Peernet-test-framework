package main

import (
	"encoding/hex"
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
	MainServerAddress string
}

// PeernetNode Node Stores all nodes generated
type PeernetNode struct {
	NodeConfig    *webapi.WebapiInstance
	FolderInfo    string
	RootNode      bool
	WebApiAddress string
}

// SetDefaults This function to be called only during a
// make install
func SetDefaults() error {
	//Setting current directory to default path
	defaultPath := "./"

	//Setting default paths for the config file
	defaults["NumberOfSlaveNode"] = 1
	defaults["NumberOfRootNode"] = 2
	defaults["MainServerAddress"] = "127.0.0.1:8000"
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
		// Sets default configuration to viper
		SetDefaults()
	}

	// Adds configuration to the struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// RunManager Starts Peernet test network for testing
func (c *Config) RunManager() (*[]PeernetNode, error) {
	// Create runs folder if it does not exist
	os.Mkdir("runs", os.ModePerm)

	now := time.Now()
	folder := "runs/test" + now.Format("20060102150405")

	// Create runs folder if it does not exist
	if err := os.Mkdir(folder, os.ModePerm); err != nil {
		return nil, err
	}

	var peernetNodes []PeernetNode
	var peernetRootNodes []PeernetNode
	var emptyRootNode []PeernetNode

	// Create Root Node
	for i := 0; i < c.NumberOfRootNode; i++ {
		// Start root node
		peernetNode, err := StartPeernet(&emptyRootNode, folder, i)
		if err != nil {
			return nil, err
		}
		// TODO remove duplicate
		peernetNodes = append(peernetNodes, *peernetNode)
		peernetRootNodes = append(peernetRootNodes, *peernetNode)
	}

	for i := 0; i < c.NumberOfSlaveNode; i++ {
		// Start regular node
		peernetNode, err := StartPeernet(&peernetRootNodes, folder, i)
		if err != nil {
			return nil, err
		}
		peernetNodes = append(peernetNodes, *peernetNode)
	}

	return &peernetNodes, nil

}

// StartPeernet Spawn Peernet instance
func StartPeernet(rootNodes *[]PeernetNode, RunPath string, number int) (*PeernetNode, error) {
	var peernetNode PeernetNode

	var folder string
	if len(*rootNodes) != 0 {
		folder = RunPath + "/run" + strconv.Itoa(number)
		peernetNode.RootNode = true
	} else {
		folder = RunPath + "/runroot" + strconv.Itoa(number)
		peernetNode.RootNode = false
	}

	// Set folder information
	peernetNode.FolderInfo = folder

	// Create runs folder if it does not exist
	if err := os.Mkdir(folder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	RunPath = folder

	var config core.Config

	_, err := core.LoadConfig("/Config.yaml", &config)
	if err != nil {
		return nil, err
	}

	// set paths for warehouse
	config.WarehouseMain = RunPath + "/data"

	UDPport, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	UDPportAddress := "127.0.0.1:" + strconv.Itoa(UDPport)
	config.Listen = []string{UDPportAddress}
	config.BlockchainGlobal = RunPath + "/data/blockchain global/"
	config.DataFolder = RunPath + "/data"
	config.BlockchainMain = RunPath + "/data/blockchain main/"
	config.LogFile = RunPath + "/data/log backend.txt"
	config.SearchIndex = RunPath + "/data/search index.txt"

	// if root node send back public key,udp address,webapi
	if len(*rootNodes) != 0 {
		var Seeds []core.PeerSeed
		rootNodeInterate := *rootNodes
		// Set root nodes to listen to
		for i, _ := range rootNodeInterate {
			var Seed core.PeerSeed
			Seed.PublicKey = hex.EncodeToString(rootNodeInterate[i].NodeConfig.Backend.PeerPublicKey.SerializeCompressed())
			Seed.Address = rootNodeInterate[i].NodeConfig.Backend.Config.Listen
			Seeds = append(Seeds, Seed)
		}
		// Iterate through all root nodes and generate
		config.SeedList = Seeds
	}

	// save configuration
	err = core.SaveConfig(RunPath+"/Config.yaml", &config)
	if err != nil {
		return nil, err
	}

	// Copy config to the appropriate folder
	backend, status, err := core.Init("Test framework/1.0", RunPath+"/Config.yaml", nil, nil)
	if status != core.ExitSuccess {
		fmt.Printf("Error %d initializing backend: %s\n", status, err.Error())
		return nil, err
	}

	backend.Connect()

	APIport, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	APIportAddress := "127.0.0.1:" + strconv.Itoa(APIport)

	api := webapi.Start(backend, []string{APIportAddress}, false, "", "", 0, 0, uuid.Nil)
	peernetNode.NodeConfig = api

	// testing reasons
	fmt.Println(APIportAddress)
	peernetNode.WebApiAddress = APIportAddress

	return &peernetNode, nil

}

// Spawn test Peernet instances for testing
func main() {

	r := mux.NewRouter()

	// Get Config information
	config, err := ConfigInit()
	if err != nil {
		fmt.Println(err)
	}

	// TODO: extend for future use-case for a embed dashboard tracker
	srv := &http.Server{
		Handler: r,
		Addr:    config.MainServerAddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	manager, err := config.RunManager()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(*manager))

	// Lister for the main server
	log.Fatal(srv.ListenAndServe())
}
