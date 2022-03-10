package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	cfgReader *configReader
)

// config.yaml ile burdakileri map etmek için viper kullanılıyor.
type (
	Configuration struct {
		MongoSettings MongoSettings
	}
	MongoSettings struct {
		DatabaseName string
		Uri          string
		Timeout      int
	}

	configReader struct {
		configFile string
		v          *viper.Viper
	}
)

func GetAllValues(configPath, configFile string) (configuration *Configuration, err error) {

	newConfigReader(configPath, configFile)
	if err = cfgReader.v.ReadInConfig(); err != nil {
		fmt.Println("Failed to read config file,Error:", err)
		return nil, errors.Wrap(err, "Config: Failed to read config file.")
	}

	err = cfgReader.v.Unmarshal(&configuration)
	if err != nil {
		fmt.Println("Failed to parse config file.", err)
		return nil, errors.Wrap(err, "Config: Failed to unmarshal yaml file to configuration object.")
	}
	return
}

func newConfigReader(folderPath, configFile string) {

	vip := viper.GetViper()
	vip.SetConfigType("yaml")
	vip.SetConfigName(configFile)
	vip.AddConfigPath(folderPath)
	cfgReader = &configReader{
		configFile: configFile,
		v:          vip,
	}
}