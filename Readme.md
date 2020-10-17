# Overview

A Veeva Vault File Staging Tool (vvfst) is a cli client using the REST API to access the staging area using newly introduced File Staging API, REST API.  Like FTP, this tool has all this CLI commands to upload, download, list, etc., files in the staging area.  The CLI wraps around the REST API for login/logout and executing REST API.

## Features
The CLI wraps around all File Staging REST API and makes easier for end user to consume the API.  High level features are:

* Upload a directory with a single command, file of any size is handled automatically.
* Download an entire directory from staging area.
* Listing of all files and folders in the staging area.
* Move/Delete any file/directory.
* Upload/Download with concurrent processes.
* Auto login if the session expired for uninterrupted usage.

## Config
The configuration and status are cache in the `$HOME/.vvfst.yaml`

## Usage
The command has self documentation
```
This cli tool connects the File Staging Area using the newly introduce File Staging REST API.
The cli authenticates using REST API then session will be cached locally for subsequent REST API calls.
Each command has unique functionality, and help doc is obtained by -h argument
Example:
  vvfst Login -h

Usage:
  vvfst [command]

Available Commands:
  download    Download folder/files remote
  help        Help about any command
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