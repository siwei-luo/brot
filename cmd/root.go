/*
Copyright © 2021-2023 Siwei Luo <siwei@lu0.org>

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
	log "github.com/sirupsen/logrus"
	"github.com/siwei-luo/brot/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

var configurationFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "brot",
	Version: pkg.Version,
	Short:   "Make brot doing the repetitive tasks.",
	Long:    `Make brot doing the repetitive tasks.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("error")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&configurationFile, "config", "c", "", "configuration file")
	rootCmd.PersistentFlags().IntVarP(&pkg.Verbosity, "verbosity", "v", 0, "verbosity (1 ~ error, 2 ~ warn, 3 ~ info, 4 ~ debug)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Override cobra's version template
	rootCmd.SetVersionTemplate("{{ .Version }}\n")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("defaults.loglevel", log.ErrorLevel)
	viper.SetDefault("defaults.logformat", "text")

	if configurationFile != "" {
		viper.SetConfigFile(configurationFile) // use configuration from the flag
	} else {
		viper.SetConfigName("brot")          // name of config file (without extension)
		viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")             // look for a configuration file in the same directory
		viper.AddConfigPath("$HOME/.config") // else try in user's home
		viper.AddConfigPath("/etc/brot/")    // if still not found look in /etc/brot
	}

	viper.AutomaticEnv() // read in environment variables that match

	// if a configuration file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("error reading configuration")
	}

	// read configuration file into struct
	if err := viper.Unmarshal(&pkg.CurrentConfiguration); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("error parsing configuration")
	}

	// read version string and remove prefix
	confVersion := strings.TrimPrefix(pkg.CurrentConfiguration.ApiVersion, "v")
	// parse major version as int
	confVersionMajor, err := strconv.ParseInt(strings.Split(confVersion, ".")[0], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("error parsing configuration version")
	}

	// check for configuration file compatibility
	if pkg.VersionMajor > confVersionMajor {
		log.WithFields(log.Fields{
			"apiVersion": pkg.CurrentConfiguration.ApiVersion,
		}).Fatal("found outdated configuration")
	}

	// set log format
	if pkg.CurrentConfiguration.Defaults.Logformat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// set log level from configuration
	var level, _ = log.ParseLevel(pkg.CurrentConfiguration.Defaults.Loglevel)
	log.SetLevel(level)

	// overwrite log level from flag
	if pkg.Verbosity > 0 {
		switch pkg.Verbosity {
		case 1:
			log.SetLevel(log.ErrorLevel)
		case 2:
			log.SetLevel(log.WarnLevel)
		case 3:
			log.SetLevel(log.InfoLevel)
		case 4:
			log.SetLevel(log.DebugLevel)
		}
	}

	// do some debug outputs
	log.Debug("parsed config file: ", viper.ConfigFileUsed())
	log.Debug("set log level: ", level)
	log.Debug("set log format: ", pkg.CurrentConfiguration.Defaults.Logformat)
}
