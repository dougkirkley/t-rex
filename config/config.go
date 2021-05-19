package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	//"github.com/spf13/viper"
)

var config TrexConfig

// TrexConfig config struct for trex
type TrexConfig struct {
	Provider  Provider    `yaml:"provider" json:"module_provider"`
	APIName   string      `yaml:"api_name" json:"api_name"`
	Functions []Functions `yaml:"functions" json:"functions"`
}

// Provider struct for aws provider
type Provider struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Region  string `yaml:"region" json:"region"`
	Profile string `yaml:"profile" json:"profile"`
}

// Functions struct for api gateway endpoints
type Functions struct {
	LambdaName    string     `yaml:"lambda_name" json:"lambda_name"`
	LambdaRuntime string     `yaml:"lambda_runtime" json:"lambda_runtime"`
	Method        string     `yaml:"method" json:"method"`
	RootPath      string     `yaml:"root_path" json:"root_path"`
	Path          string     `yaml:"path" json:"path"`
	Authorization string     `yaml:"authorization" json:"authorization"`
	IAMPolicy     string     `yaml:"iam_policy" json:"iam_policy"`
}

var defaultConfig = `
provider:
  name: aws
  version: ~> 3.2
  region: us-east-1 
  profile: default

api_name: apiname
   
functions:
- lambda_name: getOrder
  lambda_runtime: go1.x
  method: GET
  path: orders
- lambda_name: addOrder
  lambda_runtime: go1.x
  method: POST
  path: orders
- lambda_name: getOrders
  lambda_runtime: go1.x
  method: GET
  path: "{id}"
`

// exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// New returns Trexconfig
func New(configFile string) *TrexConfig {
	cfg := &TrexConfig{}
	trexCfg, _ := cfg.LoadConfig()
	return &trexCfg
}

// WorkingDirectory set the working directory
func (cfg *TrexConfig) WorkingDirectory() string {
    wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
    return wd
}

//GenerateDefaultConfig generates trex.yaml file
func (cfg *TrexConfig) GenerateDefaultConfig() error {
	filePath := filepath.Join(cfg.WorkingDirectory(), "trex.yaml")
	err := ioutil.WriteFile(filePath, []byte(defaultConfig), 0755)
	if err != nil {
		return err
	}
	return nil
}

// // LoadConfig loads yaml config for trex
// func (cfg *TrexConfig) LoadConfig() (TrexConfig, error) {
// 	viper.AddConfigPath(".")
// 	viper.SetConfigName("trex")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		return config, err
// 	}

// 	configFile := &TrexConfig{}
// 	err = viper.Unmarshal(configFile)
// 	if err != nil {
// 		return *configFile, err
// 	}
// 	return *configFile, nil
// }

// LoadConfig loads yaml config for trex
func (cfg *TrexConfig) LoadConfig() (TrexConfig, error) {
	filePath := filepath.Join(cfg.WorkingDirectory(), "/trex.yaml")
    yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil

}

// MakeVarJSON converts yaml file to terraform.tfvars.json file
func (cfg *TrexConfig) MakeVarJSON() error {
	yamlPath := filepath.Join(cfg.WorkingDirectory(), "trex.yaml")
	filePath := filepath.Join(cfg.WorkingDirectory(), "terraform.tfvars.json")
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Printf(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}
	prettyJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	err = ioutil.WriteFile(filePath, prettyJSON, 0755)
	if err != nil {
		return err
	}
	return nil
}
