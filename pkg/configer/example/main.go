package main

import (
	"encoding/json"
	"fmt"
	"github.com/gitsang/flight-metrics/pkg/configer"

	"github.com/spf13/cobra"
)

type Config struct {
	Timezone string `json:"timezone" yaml:"timezone" default:"UTC"`

	Server struct {
		Http struct {
			Host string `json:"host" yaml:"host"`
			Port int    `json:"port" yaml:"port"`
		}
		Grpc struct {
			Host string `json:"host" yaml:"host"`
			Port int    `json:"port" yaml:"port"`
		}
	}

	Data struct {
		Database struct {
			Driver  string `json:"driver" yaml:"driver"`
			Mongodb struct {
				Url      string `json:"url" yaml:"url" env:"MONGODB_URI" flag:"mongodb.uri"`
				Database string `json:"database" yaml:"database" env:"MONGODB_DATABASE" flag:"mongodb.database"`
			} `json:"mongodb" yaml:"mongodb"`
		} `json:"database" yaml:"database"`
	} `json:"data" yaml:"data"`
}

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "This is an example",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var rootFlags = struct {
	Config string
}{}

var cfger *configer.Configer

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootFlags.Config, "config", "c", "config.yml", "config file")

	cfger = configer.New(
		configer.WithTemplate((*Config)(nil)),
		configer.WithEnvBind(
			configer.WithEnvPrefix("CONFIGER"),
			configer.WithEnvDelim("_"),
		),
		configer.WithFlagBind(
			configer.WithCommand(rootCmd),
			configer.WithFlagPrefix("configer"),
			configer.WithFlagDelim("."),
		),
	)
}

func run() {
	var c Config
	err := cfger.Load(&c, rootFlags.Config)
	if err != nil {
		panic(err)
	}
	configJsonBytes, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(configJsonBytes))
}

func main() {
	rootCmd.Execute()
}
