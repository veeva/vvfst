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
	"github.com/eiannone/keyboard"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/veeva/vvfst/api"
	"github.com/veeva/vvfst/config"
	"github.com/veeva/vvfst/model"
	"github.com/veeva/vvfst/net"
	"github.com/veeva/vvfst/util"
	"github.com/veeva/vvfst/vlog"
	os "os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	recursiveOpt bool
	overwriteOpt bool
	limitOpt     int64
	threadCnt    int
)

// lsCmd represents the listCommand command
var lsCmd = &cobra.Command{
	Use:   "ls <remote-file/folder>",
	Short: "List of files and folders",
	Long:  `List the content in the directory.  The listing is a flat list when including sub directories. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, listCommand)
	},
}

var mkdirCmd = &cobra.Command{
	Use:   "mkdir <remote-folder>",
	Short: "Create remote directory",
	Long:  `Create remote directory if not exists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, mkdirCommand)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload <local-file/folder> <remote-file/folder>",
	Short: "Copy a file or folder to remote directory",
	Long:  `Uploading a single file or all files from a folder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, uploadCommand)
	},
}
var downloadCmd = &cobra.Command{
	Use:   "download <remote-file/folder> <local-file/folder>",
	Short: "Download folder/files remote",
	Long:  `Download folder/files from remote staging folder to current folder `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, downloadCommand)
	},
}
var moveCmd = &cobra.Command{
	Use:   "mv <src-remote-file/folder> <dest-remote-file/folder>",
	Short: "Move files/folder in the remote directory",
	Long:  `Move files/folder from one location another location within remote location or rename file/folder in the remote directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, mvCommand)
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm <src-remote-file/folder> <dest-remote-file/folder>",
	Short: "Remove files/folder remote location",
	Long:  "Remove files/folder remote location",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, rmCommand)
	},
}

var mlsCmd = &cobra.Command{
	Use:   "mls <remote-file/folder>",
	Short: "List multipart upload sessions",
	Long:  "List of multipart upload sessions initiated by that user which is not expired",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, mlistCommand)
	},
}
var mRmCmd = &cobra.Command{
	Use:   "mrm <remote-file>",
	Short: "Delete upload sessions",
	Long:  "Delete upload session started for the given file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithAutoLogin(cmd, args, mrmCommand)
	},
}

func init() {
	config.InitConfig()

	// List
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&recursiveOpt, "recursive", "r", false, "Enable recursive mode to list all sub directories")
	lsCmd.Flags().Int64VarP(&limitOpt, "limit", "l", 100, "Limit number of items")

	// Mkdir
	rootCmd.AddCommand(mkdirCmd)
	mkdirCmd.Flags().BoolVarP(&overwriteOpt, "overwrite", "o", false, "Enable overwrite to overwrite existing directory")

	// Upload
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&overwriteOpt, "overwrite", "o", false, "Enable overwrite to overwrite if file/folder exists")
	uploadCmd.Flags().IntVarP(&threadCnt, "threadCount", "t", 1, "Number of concurrent thread to upload")

	// Download
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolVarP(&recursiveOpt, "recursive", "r", false, "Enable recursive mode to download all sub directories")
	downloadCmd.Flags().IntVarP(&threadCnt, "threadCount", "t", 1, "Number of concurrent thread to download")

	// Move
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().BoolVarP(&overwriteOpt, "overwrite", "o", false, "Enable overwrite to overwrite a file or merge existing folders")

	// Delete
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&recursiveOpt, "recursive", "r", false, "Enable recursive mode to delete all sub directory contents")

	// Multipart List
	rootCmd.AddCommand(mlsCmd)

	// Multipart Remove session
	rootCmd.AddCommand(mRmCmd)
}

func listCommand(cmd *cobra.Command, args []string) error {
	itemPath := "/"
	if len(args) == 1 {
		itemPath = strings.TrimSpace(args[0])
	}

	if limitOpt < 0 || limitOpt > 1000 {
		return fmt.Errorf("limit must be between 0 and 1000")
	}

	firstPageItemPath := itemPath
	nextPageURL := ""
	nextPage := false
	for ok := true; ok; ok = nextPage {
		itemsRestResult, err := api.ListPage(firstPageItemPath, nextPageURL, limitOpt, recursiveOpt, true)
		nextPage = false

		if err != nil {
			return err
		}

		if itemsRestResult == nil {
			return errors.Errorf("No items found")
		}

		fmt.Printf("listing: %s\n", itemPath)
		fmt.Printf("%-6.6s  %-50.50s  %s\n", "kind", "path", "size")
		fmt.Printf("====================================================================\n")
		for _, item := range itemsRestResult.Data {
			s := util.FixedWidth(item.Path, 50, true)
			fmt.Printf("%-6.6s  %-50.50s  %s\n", item.Kind, s, util.ByteCountSI(item.Size))
		}

		if itemsRestResult.ResponseDetails != nil {
			nextPageURL = itemsRestResult.ResponseDetails.NextPage
			nextPage = continueNextPage()
		}
	}

	return nil
}

func mkdirCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("must specify a <remote-folder>")
	}

	remoteItem := strings.TrimSpace(args[0])
	api.CreateFolder(remoteItem, overwriteOpt, true)

	return nil
}

func mvCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Must specify <src-remote-file/folder> and <dest-remote-file/folder>")
	}

	srcRemoteItem := strings.TrimSpace(args[0])
	destRemoteItem := strings.TrimSpace(args[1])
	_, srcName := util.SplitParentAndName(srcRemoteItem)
	destParent, destName := util.SplitParentAndName(destRemoteItem)

	if util.EndWithFileSeparator(destRemoteItem) {
		destParent = util.TrimLastChar(destRemoteItem)
		destName = srcName
	}

	params := map[string]string{
		"parent":    destParent,
		"name":      destName,
		"overwrite": strconv.FormatBool(overwriteOpt),
	}

	var jobRestResult *model.JobRestResult
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)
	resp, err := req.SetResult(&jobRestResult).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(params).
		Put(fmt.Sprintf("/services/file_staging/items%s", srcRemoteItem))

	if err != nil {
		return errors.Errorf("Failed to connect: %v", err)
	}

	if len(jobRestResult.Errors) != 0 {
		return errors.New(net.FormatRestResultError("", jobRestResult.Errors[0]))
	}

	net.LogTime("mv submitted successfully, waiting for job completion", resp)

	return api.WaitForJobCompletion(jobRestResult.Data.JobID,
		fmt.Sprintf("%s moved to %s successfully", srcRemoteItem, destRemoteItem), config.JobTimeoutSeconds)
}

func rmCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Must specify <remote-file/folder>")
	}

	remoteItem := strings.TrimSpace(args[0])

	var jobRestResult *model.JobRestResult
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)
	resp, err := req.SetResult(&jobRestResult).
		SetQueryParam("recursive", strconv.FormatBool(recursiveOpt)).
		Delete(fmt.Sprintf("/services/file_staging/items%s", remoteItem))

	if err != nil {
		return errors.Errorf("Failed to connect: %v", err)
	}

	if len(jobRestResult.Errors) != 0 {
		return errors.New(net.FormatRestResultError("", jobRestResult.Errors[0]))
	}

	net.LogTime("rm submitted successfully, waiting for job completion", resp)

	return api.WaitForJobCompletion(jobRestResult.Data.JobID, fmt.Sprintf("%s removed successfully", remoteItem), config.JobTimeoutSeconds)
}

func uploadCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Missing required args <local-folder/file> and/or <remote-folder/file>")
	}

	localItem := strings.TrimSpace(args[0])
	remoteItem := strings.TrimSpace(args[1])

	localItemStat, err := os.Stat(localItem)
	if err != nil {
		vlog.Errorf("%s not found", localItem)
		return err
	}

	if localItemStat.Mode().IsRegular() {
		if util.EndWithFileSeparator(remoteItem) {
			remoteItem = remoteItem + localItemStat.Name()
		}
		api.UploadSingleFile(&model.UploadItem{RemotePath: remoteItem, LocalPath: localItem}, overwriteOpt)
		return nil
	}

	if util.EndWithFileSeparator(remoteItem) {
		remoteItem = util.TrimLastChar(remoteItem)
	}

	var wg sync.WaitGroup
	ch := make(chan *model.UploadItem, threadCnt)

	// run worker pool
	for i := threadCnt; i > 0; i-- {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for item := range ch {
				api.UploadSingleFile(item, overwriteOpt)
			}
		}()
	}

	err = filepath.Walk(localItem, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir)
		remotePath1 := strings.Replace(path, localItem, "", 1)
		remotePath := remoteItem + remotePath1
		if info.Mode().IsRegular() {
			uploadItem := model.UploadItem{RemotePath: remotePath, LocalPath: path}
			ch <- &uploadItem
		}

		if info.Mode().IsDir() {
			api.CreateFolder(remotePath, true, false)
		}
		return nil
	})

	close(ch)
	wg.Wait()

	return err
}

func downloadCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Missing required args <remote-folder/file> and/or <local-folder/file>")
	}

	remoteItem := strings.TrimSpace(args[0])
	localItem := strings.TrimSpace(args[1])

	firstPageItemPath := remoteItem
	nextPageURL := ""
	nextPage := false
	for ok := true; ok; ok = nextPage {
		itemsRestResult, err := api.ListPage(firstPageItemPath, nextPageURL, limitOpt, recursiveOpt, false)
		nextPage = false

		if err != nil {
			return err
		}

		if itemsRestResult == nil {
			return errors.Errorf("No items found")
		}

		downloadItems := []*model.DownloadItem{}
		for _, item := range itemsRestResult.Data {
			if item.Kind == "folder" {
				continue
			}

			localPath1 := strings.Replace(item.Path, remoteItem, "", 1)
			localPath := filepath.Join(localItem, localPath1)

			downloadItem := &model.DownloadItem{RemotePath: item.Path, Size: item.Size, LocalPath: localPath}
			downloadItems = append(downloadItems, downloadItem)
			//DownloadSingleFile(downloadItem)
		}

		downloadInParallel(downloadItems)

		if itemsRestResult.ResponseDetails != nil {
			nextPageURL = itemsRestResult.ResponseDetails.NextPage
			nextPage = true
		}
	}

	return nil
}

func mlistCommand(cmd *cobra.Command, args []string) error {
	sessionsRestResult, err := api.MultipartList(true)
	if err != nil {
		return err
	}

	fmt.Printf("listing upload sessions: \n")
	fmt.Printf("%-30.30s  %-16.16s %-16.16s %-10.10s %-16.16s\n", "path", "size", "up size", "up parts", "expiration")
	fmt.Printf("==============================================================================================\n")
	for _, item := range sessionsRestResult.Data {
		path := util.FixedWidth(item.Path, 30, true)
		fmt.Printf("%-30.30s  %-16.16s %-16.16s %-10.1d %-16.16s\n", path, util.ByteCountSI(item.Size), util.ByteCountSI(item.UploadedSize), item.UploadedPartsCount, item.ExpirationDate)
	}

	return nil
}

func mrmCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Must specify <remote-file>")
	}

	remoteItem := strings.TrimSpace(args[0])
	sessionsRestResult, err := api.MultipartList(false)
	if err != nil {
		return err
	}

	var uploadSession *model.UploadSession
	for _, item := range sessionsRestResult.Data {
		if item.Path == remoteItem {
			uploadSession = item
			break
		}
	}

	if uploadSession == nil {
		return errors.Errorf("No upload session available for file %s", remoteItem)
	}

	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)
	var restResult model.RestResult
	resp, err := req.
		SetResult(&restResult).
		Delete(fmt.Sprintf("/services/file_staging/upload/%s", uploadSession.UploadSessionID))

	if err != nil {
		return errors.Errorf("Failed to connect: %v", err)
	}

	if len(restResult.Errors) != 0 {
		return errors.Errorf(net.FormatRestResultError(remoteItem, restResult.Errors[0]))
	}

	net.LogTime(fmt.Sprintf("Deleted upload session for %s", remoteItem), resp)

	return nil
}

func downloadInParallel(items []*model.DownloadItem) {
	var wg sync.WaitGroup
	ch := make(chan *model.DownloadItem, threadCnt)

	// run worker pool
	for i := threadCnt; i > 0; i-- {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for item := range ch {
				api.DownloadSingleFile(item)
			}
		}()
	}

	for _, i := range items {
		ch <- i
	}

	close(ch)
	wg.Wait()
}

func continueNextPage() bool {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	vlog.Infof("Press Ctl+C to stop or press space-bar for next page")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if event.Key == keyboard.KeyCtrlC {
			return false
		}
		if event.Key == keyboard.KeySpace {
			return true
		}
	}
}

func runWithAutoLogin(cmd *cobra.Command, args []string, cmdFunc func(cmd *cobra.Command, args []string) error) error {
	err := cmdFunc(cmd, args)
	if net.IsSessionExpired(err) {
		vlog.Infof("Session expired, auto Login")
		if err := api.Login(); err == nil {
			return cmdFunc(cmd, args)
		}
	}
	return err
}
