package registries

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type ConfigurationRegistry struct {
	c *Configuration
	e error
}

type Configuration struct {
	ItemExpireInterval uint64 `yaml:"itemExpireInterval"`
	ItemExpireCheck    uint64 `yaml:"itemExpireCheck"`
}

var configurationRegistryOnce sync.Once
var configurationRegistry *ConfigurationRegistry

func GetConfiguration() (*Configuration, error) {
	configurationRegistryOnce.Do(func() {
		configurationRegistry = &ConfigurationRegistry{}
		err := configurationRegistry.loadConfiguration()
		configurationRegistry.e = err
	})
	return configurationRegistry.c, configurationRegistry.e
}

func (c *ConfigurationRegistry) loadConfiguration() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	con := &Configuration{}
	err = yaml.Unmarshal(yamlFile, con)
	if err != nil {
		return err
	}
	c.c = con
	return nil
}
