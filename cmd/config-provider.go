package cmd

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "gobyoall",
		Short: "Gobyoall is a flexible stress-tester for servers",
		/// TODO: provide more info, documentations
		Long: `See https://github.com/runar-rkmedia/gabyoall`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.SetConfigName("gobyoall-conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path.Join(home, "gobyall"))
		viper.AddConfigPath(path.Join(home, ".config", "gobyall"))
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix("gobyoall")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
}

func ReadConfig() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/gobyoall-conf.yaml)")
	var cfg Config
	t := reflect.TypeOf(cfg)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		cfgName := field.Tag.Get("cfg")
		if cfgName == "-" {

			continue
		}
		if cfgName == "" {
			panic(fmt.Sprintf("field %d has no name %#v", i+1, field))
		}
		short := field.Tag.Get("short")
		mapstructure := field.Tag.Get("mapstructure")
		if mapstructure == "" {
			mapstructure = strings.ToLower(field.Name[0:1]) + field.Name[1:]
		}
		defaultStr := field.Tag.Get("default")
		desc := field.Tag.Get("description")
		kind := field.Type.Name()
		if kind == "" {
			kind = field.Type.String()
		}
		switch kind {
		case "bool":
			defaultValue := defaultStr == "true"
			rootCmd.PersistentFlags().BoolP(cfgName, short, defaultValue, desc)
		case "string":
			rootCmd.PersistentFlags().StringP(cfgName, short, defaultStr, desc)

		case "int":
			defaultInt := 0
			if defaultStr != "" {
				n, err := strconv.ParseInt(defaultStr, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("failed to convert default-tag (%s) on config-field %s", defaultStr, field.Name))
				}
				defaultInt = int(n)
			}
			rootCmd.PersistentFlags().IntP(cfgName, short, defaultInt, desc)
		case "[]int":
			var defaultInts []int
			if defaultStr != "" {
				split := strings.Split(defaultStr, ",")
				for i := 0; i < len(split); i++ {
					n, err := strconv.ParseInt(split[i], 10, 64)
					if err != nil {
						panic(fmt.Sprintf("failed to convert default-tag (%s) on config-field %s", defaultStr, field.Name))
					}
					defaultInts = append(defaultInts, int(n))
				}
			}
			rootCmd.PersistentFlags().IntSliceP(cfgName, short, defaultInts, desc)
		case "interface {}":
			rootCmd.PersistentFlags().StringP(cfgName, short, defaultStr, desc)
		case "map[string]string":
			var defaultStrings map[string]string
			if defaultStr != "" {
				split := strings.Split(defaultStr, ";")
				for i := 0; i < len(split); i++ {
					kv := strings.Split(defaultStr, "=")
					defaultStrings[kv[0]] = kv[1]
				}
			}
			rootCmd.PersistentFlags().StringToStringP(cfgName, short, defaultStrings, desc)
		case "map[string]interface {}":
			continue
		default:
			panic(fmt.Sprintf("no handler for %s, %s", field.Name, kind))
		}
		viper.BindPFlag(mapstructure, rootCmd.PersistentFlags().Lookup(cfgName))
	}

}

func Execute() error {
	ReadConfig()
	return rootCmd.Execute()
}
