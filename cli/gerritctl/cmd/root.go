/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shijl0925/go-gerrit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gerritctl",
	Short:   "A client for gerrit",
	Version: "v0.0.1",
	Long:    `Client for gerrit, manage resources by the gerrit`,
}

// GerritMod connection object
type GerritMod struct {
	Instance *gerrit.Gerrit
	Url      string
	Username string
	Password string
	Context  context.Context
}

// Init will initilialize connection with gerrit server
//
// Args:
//
// Returns
func (g *GerritMod) Init(config Config) error {
	g.Username = config.Username
	g.Url = config.Url
	g.Password = config.Password
	g.Context = context.Background()

	client, err := gerrit.NewClient(g.Url)
	if len(g.Username) != 0 && len(g.Password) != 0 {
		client.SetBasicAuth(g.Username, g.Password)
	}

	g.Instance = client

	return err
}

// Config is focused in the configuration json file
type Config struct {
	Url            string `mapstructure: Url`
	Username       string `mapstructure: Username`
	Password       string `mapstructure: Password`
	ConfigPath     string
	ConfigFileName string
}

// SetConfigPath set the default config path
//
// Args:
//
// Returns
//
//	string or error
func (c *Config) SetConfigPath(path string) {
	dir, file := filepath.Split(path)
	c.ConfigPath = dir
	c.ConfigFileName = file
}

// CheckIfExists check if file exists
//
// Args:
//
//	path - string
//
// Returns
//
//	error
func (c *Config) CheckIfExists() error {
	var err error
	if _, err = os.Stat(c.ConfigPath + c.ConfigFileName); err == nil {
		return nil

	}
	return err
}

// LoadConfig read the JSON configuration from specified file
//
// Example file:
//
// $HOME/.config/gerritctl/config.json
//
// Args:
//
// Returns
//
//	nil or error
func (c *Config) LoadConfig() (config Config, err error) {
	viper.AddConfigPath(c.ConfigPath)
	viper.SetConfigName(c.ConfigFileName)
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

var gerritConfig Config
var gerritMod GerritMod
var configFile string

var Verbose bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "", "", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if configFile != "" {
		gerritConfig.SetConfigPath(configFile)
	} else {
		gerritConfig.SetConfigPath(dirname + "/.config/gerritctl/config.json")
	}

	config, err := gerritConfig.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gerritMod = GerritMod{}
	err = gerritMod.Init(config)
	if err != nil {
		fmt.Println("❌ gerrit server unreachable: " + gerritMod.Url)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
