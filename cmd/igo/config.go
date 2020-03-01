package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type configuration struct {
	Test string `yaml:"TEST,omitempty"`
}

var config = &configuration{
	Test: "hello world",
}

// PostProcess post-processes configuration after loading it.
func (c *configuration) PostProcess() error {
	return nil
}

// Validate validates the configuration.
func (c *configuration) Validate() error {
	return nil
}

func dumpDefaultConfig(fileName string) {
	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	defer fp.Close()
	if err != nil {
		return
	}
	var data []byte
	data, err = yaml.Marshal(config)
	if err != nil {
		return
	}
	fp.Write(data)
}

func configure() {
	home, err := homedir.Dir()
	if err != nil {
		er(err)
	}
	configPath := path.Join(home, ".igo")
	if ok, _ := pathExists(configPath); !ok {
		if err := os.Mkdir(configPath, os.ModePerm); err != nil {
			fmt.Println("Warning: can not create config dir:", configPath)
		}
	}
	configFullName := path.Join(configPath, "igo.conf")
	if ok, _ := pathExists(configFullName); !ok {
		dumpDefaultConfig(configFullName)
	}

	// Viper settings
	viper.AddConfigPath(fmt.Sprintf("$%s_CONFIG_DIR/", strings.ToUpper(envPrefix)))
	viper.AddConfigPath(configPath)
	viper.SetConfigName("igo.conf")
	viper.SetConfigType("yaml")

	// Environment variable settings
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	// Application constants
	viper.Set("appName", appName)

	err = viper.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if configFileNotFound {
		fmt.Println("configuration file not found")
	} else {
		emperror.Panic(errors.Wrap(err, "failed to process configuration"))
	}
	if err != nil && *verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(config)
	emperror.Panic(errors.Wrap(err, "failed to unmarshal configuration"))

	err = config.PostProcess()
	emperror.Panic(errors.WithMessage(err, "failed to post-process configuration"))
}
