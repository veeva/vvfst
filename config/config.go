/*
This code serves as an example and is not meant for production use.

Copyright 2020 Veeva Systems Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions
and limitations under the License.
*/
package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/veeva/vvfst/model"
	"github.com/veeva/vvfst/vlog"
	"os"
	"path/filepath"
	"strconv"
)

const (
	Size5MB           = 5 * 1024 * 1024
	Size50MB          = 50 * 1024 * 1024
	JobTimeoutSeconds = 60
)

var EnableDebug bool

var cfgFile string
var initialized bool

const (
	ConfigKeyDomainName   = "domain_name"
	ConfigKeyAPIVersion   = "api_version"
	ConfigKeyUsername     = "username"
	ConfigKeyPassword     = "password"
	ConfigAuthResult      = "auth_result"
	ConfigUploadSessionID = "upload_session_id"
)

func DomainName() string {
	return viper.GetString(ConfigKeyDomainName)
}

func APIVersion() string {
	return viper.GetString(ConfigKeyAPIVersion)
}

func Username() string {
	return viper.GetString(ConfigKeyUsername)
}

func PasswordMasked() string {
	if Password() == "" {
		return ""
	}
	return "****"
}
func Password() string {
	return viper.GetString(ConfigKeyPassword)
}

func UploadSessionID() string {
	return viper.GetString(ConfigUploadSessionID)
}

func SaveUploadSessionID(sessionID string) {
	viper.Set(ConfigUploadSessionID, sessionID)
}

func SaveAuthResult(result *model.AuthResult) {
	authResult := map[string]interface{}{
		"session_id": result.SessionID,
		"vault_id":   result.VaultID,
		"user_id":    result.UserID,
	}
	viper.Set(ConfigAuthResult, authResult)
}

func AuthResult() *model.AuthResult {
	if !viper.IsSet(ConfigAuthResult) {
		return nil
	}
	authMap := viper.GetStringMapString(ConfigAuthResult)
	authResult := &model.AuthResult{
		SessionID: authMap["session_id"],
	}

	if userID, ok := authMap["vault_id"]; ok {
		authResult.VaultID, _ = strconv.Atoi(userID)
	}

	if userID, ok := authMap["user_id"]; ok {
		authResult.UserID, _ = strconv.Atoi(userID)
	}

	return authResult
}

func UpdateConfig() {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if _, err := os.Create(cfgFile); err != nil {
			vlog.Info("Config not exists, creating it")
		}
	}

	err := viper.WriteConfig()
	if err != nil {
		vlog.Errorf("Error updating config file: %v", err)
	}
}

func ResetAuthResult() {
	viper.Set(ConfigAuthResult, "")
}

func ResetConfig() {
	viper.Reset()
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		vlog.Info("Config file not exists.")
		return
	}

	err := os.Remove(cfgFile)
	if err != nil {
		vlog.Errorf("Fail to remove config file: %v", err)
	}
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if initialized {
		return
	}

	vlog.InitLog(false)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vvfst" (without extension).
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.SetConfigName(".vvfst")

		cfgFile = filepath.Join(home, ".vvfst.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		vlog.Errorf("Configuration not initialized")
	}

	initialized = true
}
