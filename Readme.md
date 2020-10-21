
![Go](https://github.com/veeva/vvfst/workflows/Go/badge.svg?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/veeva/vvfst)](https://goreportcard.com/report/github.com/veeva/vvfst) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/veeva/vvfst/blob/main/LICENSE)  
  
  
# Overview  
  
A Veeva Vault File Staging Tool (vvfst) is a cli client using the REST API to access the staging area using newly introduced File Staging API, REST API.  Like FTP, this tool has all CLI commands to upload, download, list, etc., files in the staging area.  The CLI wraps around the REST API for login/logout and executing REST API.  
  
## Features  
The CLI wraps around all File Staging REST API and makes easier for end user to consume the API.  High level features are:  
  
* Upload a directory with a single command, file of any size is handled automatically.  
* Download an entire directory from staging area.  
* Listing of all files and folders in the staging area.  
* Move/Delete any file/directory.  
* Upload/Download with concurrent processes.  
* Auto login if the session expired for uninterrupted usage.  

# Demo
[![asciicast](https://asciinema.org/a/iWzJve3MUH69EpFZZZqmlHas5.svg)](https://asciinema.org/a/iWzJve3MUH69EpFZZZqmlHas5)
  
# Installation  
The CLI is built into your native platform and available in the release section [Release Section](https://github.com/veeva/vvfst/releases).  

* Download the distribution as per your OS 
* Extract and copy the vvfst into your accessible path
	* For linux/osx it would be, `/usr/local/bin`
	* For windows it would be, `C:\Windows\System32`
	
**Note mac user**  
  This tool is not distributed through app store hence the mac will complain about security.  If you see a security dialog,
  ![Security Warning](https://github.com/veeva/vvfst/blob/main/security-warning.png)
  
  go to system preference and click on the security and allow `vvfst` to run
  ![Unblock Security warning](https://github.com/veeva/vvfst/blob/main/security-allow.png)
  	
  
## Usage  
The command has self documentation  
```  
This cli tool connects the File Staging Area using the newly introduce File Staging REST API.
The cli authenticates using REST API then session will be cached locally for subsequent REST API calls.
Each command has unique functionality, and help doc is obtained by -h argument
Example:
  vvfst login -h

Usage:
  vvfst [command]

Available Commands:
  download    Download folder/files remote
  help        Help about any command
  jobs        Display list of active jobs and check status
  login       Login to the vault
  logout      logout form current cli session
  ls          List of files and folders
  mkdir       Create remote directory
  mls         List multipart upload sessions
  mrm         Delete upload sessions
  mv          Move files/folder in the remote directory
  rm          Remove files/folder remote location
  upload      Copy a file or folder to remote directory

Flags:
  -x, --debug   Enable debug
  -h, --help    help for vvfst

Use "vvfst [command] --help" for more information about a command.
```  
  
Note: 
The configuration, such as login credential, domain and status are cache in the `$HOME/.vvfst.yaml`

# Commands
Usage of each commands with example found here [Commands](https://github.com/veeva/vvfst/blob/main/commands.md)

  
# Development  
If you wanted to add more functionality and develop from scratch then, these steps are for you.

* Install the golang 
* Checkout the source `git checkout https://github.com/veeva/vvfst`
* Run the make command `make build` or run other make commands as per your OS as well.
  

# TODO 
There are multiple nice to have open items

* Support multiple login like a profile
* Add progress for upload progress
* Upload/download resume from a directory
* Add config command and move flags from login to config command


  
# License  
This code serves as an example and is not meant for production use.  
  
Copyright 2020 Veeva Systems Inc.  
  
Licensed under the Apache License, Version 2.0 (the "License"); you may not use  
this file except in compliance with the License. You may obtain a copy of the License at  
  
```  
http://www.apache.org/licenses/LICENSE-2.0  
```  
  
Unless required by applicable law or agreed to in writing, software distributed under  
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,  
either express or implied. See the License for the specific language governing permissions  
and limitations under the License.
