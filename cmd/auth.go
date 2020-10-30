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
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veeva/vvfst/api"
	"github.com/veeva/vvfst/config"
	"github.com/veeva/vvfst/vlog"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

// loginCmd represents the Login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the vault",
	Long: `Login with username and password for the vault.  
For example:
  Login --domain_name myvalut.veevavault.com --Login myuser@mydomain.com --password mypassword
  Login -d myvault.veevavault.com -a v20.1 -u myuser@mydomain.com -p mypassword
`,
	RunE: loginCommand,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout form current cli session",
	Long:  `logout and delete all cached session data.`,
	Run:   logout,
}

var clearOpt bool

func init() {
	config.InitConfig()

	// Login
	rootCmd.AddCommand(loginCmd)
	buildCmdOption(loginCmd, "domain_name", "d", "Vault domain name", config.DomainName, config.ConfigKeyDomainName)
	buildCmdOption(loginCmd, "username", "u", "Vault username", config.Username, config.ConfigKeyUsername)
	buildCmdOption(loginCmd, "api_version", "a", "API Version", config.APIVersion, config.ConfigKeyAPIVersion)

	// logout
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().BoolVarP(&clearOpt, "clear", "c", false, "Clear all configuration data")
}

func loginCommand(_ *cobra.Command, _ []string) error {
	config.SetPassword(readPassword())
	return api.Login()
}

func logout(_ *cobra.Command, _ []string) {
	if clearOpt {
		fmt.Println("Clearing configuration..")
		config.ResetConfig()
		return
	}

	config.ResetAuthResult()
	config.UpdateConfig()

	vlog.Info("logout successful.")
}

func buildCmdOption(cmd *cobra.Command, flagName, flagNameShort, flagDescription string, configGetter func() string, configName string) {
	cmd.PersistentFlags().StringP(flagName, flagNameShort, configGetter(), flagDescription)
	if configGetter() == "" {
		_ = cmd.MarkPersistentFlagRequired(configName) // Required param if config is missing
	}
	_ = viper.BindPFlag(configName, cmd.PersistentFlags().Lookup(flagName)) // read from cli or config
}

func readPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic("Failed to read password")
	}
	password := string(bytePassword)

	return strings.TrimSpace(password)
}
