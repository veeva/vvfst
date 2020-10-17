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
	"github.com/spf13/cobra"
	"github.com/veeva/vvfst/config"
	"github.com/veeva/vvfst/vlog"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vvfst",
	Short: "A cli tool to manage files in the File Staging Area using File Staging REST API",
	Long: `This cli tool connects the File Staging Area using the newly introduce File Staging REST API. 
The cli authenticates using REST API then session will be cached locally for subsequent REST API calls. 
Each command has unique functionality, and help doc is obtained by -h argument
Example:
  vvfst Login -h

LICENSE:
=============================================================================================
This code serves as an example and is not meant for production use.

Copyright 2020 Veeva Systems Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions
and limitations under the License.
=============================================================================================
`,
}

func init() {
	config.InitConfig()

	rootCmd.PersistentFlags().BoolVarP(&config.EnableDebug, "debug", "x", false, "Enable debug")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		vlog.Errorf("%v", err)
		os.Exit(1)
	}
}
