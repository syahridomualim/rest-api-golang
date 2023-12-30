package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

var IsLoadConfigDone = false
var _ error

func InitConfig(force bool) {
	if IsLoadConfigDone && !force {
		return
	}
	envVariables()
	if getCloudConfigUrl() != "" {
		logrus.WithFields(logrus.Fields{
			"at": time.Now().Format(time.RFC3339),
		}).Infof("load cloud config URL: %s", getCloudConfigUrl())
		_ = springCloudConfig("", getCloudConfigUrl())
	} else {
		logrus.WithFields(logrus.Fields{
			"at": time.Now().Format(time.RFC3339),
		}).Warn("cloud config URL is not defined, please set cloud config url in env variable: CLOUD_CONFIG_URL")
	}
	IsLoadConfigDone = true
}

func getCloudConfigUrl() string {
	return viper.GetString("cloud.config.url")
}

func envVariables() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}

type cloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           interface{}      `json:"label"`
	Version         string           `json:"version"`
	State           interface{}      `json:"state"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func springCloudConfig(prefix, url string) error {
	body, err := callSpringCloudConfig(url)
	if err != nil {
		logrus.Errorf("SpringCloudConfig: %s\n", err)
		panic(err)
	}
	cloudConfig := new(cloudConfig)
	err = json.Unmarshal(body, cloudConfig)
	if err != nil {
		logrus.Errorf("SpringCloudConfig: %s\n", err)
		return err
	}
	for _, vps := range cloudConfig.PropertySources {
		for is, vs := range vps.Source {
			if prefix == "" {
				setIfNotExists(is, vs)
			} else {
				setIfNotExists(fmt.Sprintf("%s.%s", prefix, is), vs)
			}
		}
	}
	return nil
}

func setIfNotExists(k string, v interface{}) {
	if viper.Get(k) != nil {
		return
	}
	viper.Set(k, v)
	return
}

func callSpringCloudConfig(url string) ([]byte, error) {
	timeoutDur, err := time.ParseDuration(os.Getenv("CLOUD_CONFIG_TIMEOUT_DURATION"))
	if err != nil {
		timeoutDur = time.Minute
	}

	retryMax, err := strconv.Atoi(os.Getenv("CLOUD_CONFIG_RETRY_MAX"))
	if err != nil {
		retryMax = 3
	}

	rest := resty.New().
		SetTimeout(timeoutDur).
		SetRetryCount(retryMax).
		SetRetryWaitTime(time.Minute)
	trx, err := rest.R().SetHeader("Content-Type", "application/json").Get(url)
	if err != nil {
		if trx != nil {
			logrus.Debug(trx)
		}
		logrus.Error(err)
		return nil, fmt.Errorf("cloud config error services %s with error: %s", url, err)
	}
	return trx.Body(), nil
}
