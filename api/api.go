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
package api

import (
	"bufio"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/veeva/vvfst/config"
	"github.com/veeva/vvfst/model"
	"github.com/veeva/vvfst/net"
	"github.com/veeva/vvfst/util"
	"github.com/veeva/vvfst/vlog"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	remoteDirCache = map[string]bool{
		".": true,
		"":  true,
	}
)

// Login with username, password with configured in the config
func Login() error {
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(false)

	var authResult model.AuthResult
	resp, err := req.
		SetFormData(map[string]string{
			"username": config.Username(),
			"password": config.Password()}).
		SetResult(&authResult).
		Post("/auth")

	if err != nil {
		return errors.Errorf("Failed to connect: %v", err)
	}

	if len(authResult.Errors) != 0 {
		return errors.New(net.FormatRestResultError("", authResult.Errors[0]))
	}

	net.LogTime("Login successful.", resp)

	config.SetAuthResult(&authResult)
	config.UpdateConfig()
	return nil
}

// List items in the page, nextPageUrl is null then it will be the first page.
func ListPage(itemPath, nextPageURL string, limit int64, recursiveOpt, logStatus bool) (*model.ItemsRestResult, error) {
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var itemsRestResult *model.ItemsRestResult
	var resp *resty.Response
	var err error
	if nextPageURL == "" {
		resp, err = req.SetResult(&itemsRestResult).
			SetQueryParam("recursive", strconv.FormatBool(recursiveOpt)).
			SetQueryParam("limit", strconv.FormatInt(limit, 10)).
			Get(fmt.Sprintf("/services/file_staging/items%s", itemPath))
	} else {
		resp, err = req.SetResult(&itemsRestResult).Get(nextPageURL)
	}

	if err != nil {
		return nil, errors.Errorf("Failed to connect: %v", err)
	}

	if itemsRestResult == nil {
		return nil, errors.Errorf("Unknown error, response is empty")
	}

	if logStatus {
		net.LogTime("ls completed.", resp)
	}

	if len(itemsRestResult.Errors) != 0 {
		return nil, errors.New(net.FormatRestResultError("", itemsRestResult.Errors[0]))
	}

	return itemsRestResult, nil
}

// List items in the page, nextPageUrl is null then it will be the first page.
func ListExport(itemPath string, recursiveOpt bool) (*model.JobRestResult, error) {
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var jobRestResult *model.JobRestResult
	var err error
	_, err = req.SetResult(&jobRestResult).
		SetQueryParam("recursive", strconv.FormatBool(recursiveOpt)).
		SetQueryParam("format_result", "csv").
		Get(fmt.Sprintf("/services/file_staging/items%s", itemPath))

	if err != nil {
		return nil, errors.Errorf("Failed to connect: %v", err)
	}

	if jobRestResult == nil {
		return nil, errors.Errorf("Unknown error, response is empty")
	}

	if len(jobRestResult.Errors) != 0 {
		return nil, errors.New(net.FormatRestResultError("", jobRestResult.Errors[0]))
	}

	return jobRestResult, nil
}

// Make the directory and ignores, dot, empty space directory
func CreateFolder(remotePath string, overwrite, logStatus bool) {
	if _, ok := remoteDirCache[remotePath]; ok {
		return
	}

	formData := map[string]string{
		"path":      remotePath,
		"name":      filepath.Base(remotePath),
		"kind":      "folder",
		"overwrite": strconv.FormatBool(overwrite),
	}

	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var itemRestResult model.ItemRestResult
	resp, err := req.
		SetResult(&itemRestResult).
		SetMultipartFormData(formData).
		Post("/services/file_staging/items")

	if err != nil {
		vlog.Errorf("Failed to connect: %v", err)
		return
	}

	if len(itemRestResult.Errors) != 0 {
		if itemRestResult.Errors[0].Type == "ITEM_NAME_EXISTS" && !logStatus {
			return
		}
		net.LogRestResultError(remotePath, itemRestResult.Errors[0])
		return
	}

	if logStatus {
		net.LogTime(fmt.Sprintf("created folder: %s", remotePath), resp)
	}
}

// Download single from the file staging area
func DownloadSingleFile(downloadItem *model.DownloadItem) {
	vlog.Debugf("Download file: %s, size: %d ", downloadItem.RemotePath, downloadItem.Size)
	req := net.InitRestClient(config.EnableDebug).
		BuildRestRequest(true).
		SetDoNotParseResponse(true)

	var err error
	var resp *resty.Response
	if downloadItem.RemoteHref != "" {
		resp, err = req.Get(fmt.Sprintf("https://%s%s", config.DomainName(), downloadItem.RemoteHref))
	} else {
		resp, err = req.Get(fmt.Sprintf("/services/file_staging/items/content%s", downloadItem.RemotePath))
	}

	if err != nil {
		vlog.Errorf("Failed to download file: %s, error: %v", downloadItem.RemotePath, err)
		return
	}
	defer func() {
		err = resp.RawBody().Close()
		if err != nil {
			vlog.Errorf("Error closing http response")
		}
	}()

	localParentDir := filepath.Dir(downloadItem.LocalPath)
	localParentStat, err := os.Stat(localParentDir)
	if localParentStat == nil {
		err := os.MkdirAll(localParentDir, 0755)
		if err != nil {
			vlog.Errorf("Failed to create directory: %s, err: %v", localParentDir, err)
			return
		}
	}
	if localParentStat != nil && !localParentStat.IsDir() {
		vlog.Errorf("Cannot create directory, a same filename exists: %s", localParentDir)
		return
	}

	f, _ := os.OpenFile(downloadItem.LocalPath, os.O_CREATE|os.O_WRONLY, 0644)
	bar := buildProgressbar(filepath.Base(downloadItem.LocalPath), downloadItem.Size)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.RawBody())

	if err != nil {
		vlog.Errorf("Failed to download file: %s, err: %v", downloadItem.RemotePath, err)
		return
	}
}

//UploadSingleFile - uploads single file using if size is less than 50MB
func UploadSingleFile(uploadItem *model.UploadItem, overwriteOpt bool) {
	fi, err := os.Stat(uploadItem.LocalPath)
	if err != nil {
		vlog.Errorf("%s file not found, err: %v", uploadItem.LocalPath, err)
		return
	}

	if fi.Size() > config.Size50MB {
		err := MultipartUploadSingleFile(uploadItem.LocalPath, uploadItem.RemotePath, overwriteOpt)
		if err != nil {
			vlog.Errorf("%v", err)
		}
		return
	}

	formData := map[string]string{
		"path":      uploadItem.RemotePath,
		"name":      util.GetFilename(uploadItem.RemotePath),
		"size":      strconv.FormatInt(fi.Size(), 10),
		"kind":      "file",
		"overwrite": strconv.FormatBool(overwriteOpt),
	}

	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var itemRestResult model.ItemRestResult
	resp, err := req.
		SetResult(&itemRestResult).
		SetMultipartFormData(formData).
		SetFile("file", uploadItem.LocalPath).
		Post("/services/file_staging/items")

	if err != nil {
		vlog.Errorf("Failed to connect: %v", err)
		return
	}

	if len(itemRestResult.Errors) != 0 {
		net.LogRestResultError(uploadItem.RemotePath, itemRestResult.Errors[0])
		return
	}

	net.LogTime(fmt.Sprintf("uploaded file: %s", uploadItem.RemotePath), resp)
}

//MultipartList - list all active multipart session
func MultipartList(logStatus bool) (*model.UploadSessionsRestResult, error) {
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var sessionsRestResult *model.UploadSessionsRestResult
	var resp *resty.Response
	var err error
	resp, err = req.SetResult(&sessionsRestResult).
		Get("/services/file_staging/upload")

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	if logStatus {
		net.LogTime("mListCmd completed.", resp)
	}

	if len(sessionsRestResult.Errors) != 0 {
		return nil, errors.New(net.FormatRestResultError("", sessionsRestResult.Errors[0]))
	}
	return sessionsRestResult, nil
}

//MultipartUploadSingleFile - Upload single file using multipart
func MultipartUploadSingleFile(localPath, remotePath string, overwriteOpt bool) error {
	fi, err := os.Stat(localPath)
	if err != nil {
		return errors.Errorf("%s file not found", localPath)
	}

	if fi.Size() < config.Size5MB {
		return errors.Errorf("%s file is less than %d", localPath, config.Size5MB)
	}

	sessionsRestResult, err := MultipartList(false)
	if err != nil {
		return err
	}

	var uploadSession *model.UploadSession
	for _, item := range sessionsRestResult.Data {
		if item.Path == remotePath {
			uploadSession = item
			break
		}
	}

	if uploadSession == nil {
		uploadSession, err = MultipartUploadBegin(localPath, remotePath, overwriteOpt)
		if err != nil {
			return err
		}
	}

	config.SetUploadSessionID(uploadSession.UploadSessionID)
	config.UpdateConfig()

	err = MultipartUploadFilePart(localPath, uploadSession)
	if err != nil {
		return err
	}

	return MultipartUploadCommit(uploadSession)
}

//MultipartUploadBegin - Begin multipart upload session
func MultipartUploadBegin(localPath, remotePath string, overwriteOpt bool) (*model.UploadSession, error) {
	fi, err := os.Stat(localPath)
	if err != nil {
		return nil, errors.Errorf("%s file not found", localPath)
	}

	formData := map[string]string{
		"path":      remotePath,
		"name":      filepath.Base(remotePath),
		"size":      strconv.FormatInt(fi.Size(), 10),
		"overwrite": strconv.FormatBool(overwriteOpt),
	}

	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var sessionRestResult model.UploadSessionRestResult
	resp, err := req.
		SetResult(&sessionRestResult).
		SetMultipartFormData(formData).
		Post("/services/file_staging/upload")

	if err != nil {
		return nil, errors.Errorf("Failed to connect: %v", err)
	}

	if len(sessionRestResult.Errors) != 0 {
		return nil, errors.Errorf(net.FormatRestResultError(remotePath, sessionRestResult.Errors[0]))
	}

	net.LogTime(fmt.Sprintf("upload session created for file: %s", remotePath), resp)

	uploadSession := sessionRestResult.Data
	return uploadSession, nil
}

//MultipartUploadFilePart - Upload multipart file
func MultipartUploadFilePart(localPath string, session *model.UploadSession) error {
	uploadSessionID := config.UploadSessionID()

	if uploadSessionID == "" {
		return errors.Errorf("Upload session not found for filepath: %s", localPath)
	}

	size, chunkSize, err := findChunkSize(localPath)
	if err != nil {
		return err
	}

	file, err := os.Open(localPath)
	if err != nil {
		return errors.Errorf("Failed to open file: %s", localPath)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, chunkSize)
	remainingSize := size - session.UploadedSize
	partNumber := 1
	totalParts := int(math.Ceil(float64(size) / float64(chunkSize)))

	if session.UploadedPartsCount > 0 {
		partNumber = session.UploadedPartsCount + 1
	}

	for ; ; partNumber++ {
		if remainingSize == 0 {
			return nil
		}

		if remainingSize < chunkSize {
			buffer = make([]byte, remainingSize)
		}

		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Errorf("cannot read chunk to buffer, err: %v", err)
		}

		var partRestResult model.UploadPartRestResult
		req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)
		_, err = req.
			SetResult(&partRestResult).
			SetHeader("X-VaultAPI-FilePartNumber", strconv.FormatInt(int64(partNumber), 10)).
			SetHeader("Content-Length", strconv.FormatInt(int64(n), 10)).
			SetHeader("Content-Type", "application/octet-stream").
			SetBody(buffer).
			Put(fmt.Sprintf("/services/file_staging/upload/%s", uploadSessionID))

		if err != nil {
			return errors.Errorf("Failed to upload file part: %s", localPath)
		}

		remainingSize = remainingSize - int64(n)

		vlog.Infof("[%s] Uploaded part: %d of %d, size: %s, partContentMD5: %s", session.Path,
			partRestResult.Data.PartNumber, totalParts, util.ByteCountSI(partRestResult.Data.PartSize), partRestResult.Data.PartContentMD5)
	}

	return nil
}

// Commit the Multipart session
func MultipartUploadCommit(uploadSession *model.UploadSession) error {
	req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)

	var jobRestResult model.JobRestResult
	resp, err := req.
		SetResult(&jobRestResult).
		Post(fmt.Sprintf("/services/file_staging/upload/%s", uploadSession.UploadSessionID))

	if err != nil {
		return errors.Errorf("Failed to connect: %v", err)
	}

	if len(jobRestResult.Errors) != 0 {
		return errors.Errorf(net.FormatRestResultError(uploadSession.Path, jobRestResult.Errors[0]))
	}

	net.LogTime(fmt.Sprintf("upload session completed for file: %s, waiting for job completion", uploadSession.Path), resp)
	msg := fmt.Sprintf("%s file upload sucessfully", uploadSession.Path)

	_, err = WaitForJobCompletion(jobRestResult.Data.JobID, msg, config.JobTimeoutSeconds)
	return err
}

// Check for job status every 10 seconds
func WaitForJobCompletion(jobID int64, message string, timeoutSec int) (*model.Link, error) {
	completionTime := time.Now().Add(time.Second * time.Duration(timeoutSec))
	jobIDStr := strconv.FormatInt(jobID, 10)
	config.UpdateActiveJob(jobIDStr, message)

	time.Sleep(time.Second) // first sleep for a second

	for time.Now().Before(completionTime) {
		req := net.InitRestClient(config.EnableDebug).BuildRestRequest(true)
		var jobStatusRestResult model.JobStatusRestResult
		_, err := req.
			SetResult(&jobStatusRestResult).
			Get(fmt.Sprintf("/services/jobs/%d", jobID))

		if err != nil {
			return nil, errors.Errorf("Failed to connect: %v", err)
		}

		if len(jobStatusRestResult.Errors) != 0 {
			return nil, errors.Errorf(net.FormatRestResultError(jobIDStr, jobStatusRestResult.Errors[0]))
		}

		if jobStatusRestResult.Data.Status == "SUCCESS" {
			if message != "" {
				vlog.Infof(message)
			}

			config.RemoveActiveJob(jobIDStr)

			var resultLink *model.Link
			for _, link := range jobStatusRestResult.Data.Links {
				if link.Rel == "results" {
					resultLink = link
					break
				}
			}

			return resultLink, nil
		}

		vlog.Infof("Current job status: %s", jobStatusRestResult.Data.Status)

		time.Sleep(10 * time.Second)
	}

	return nil, errors.Errorf("Job not completed within %d seconds", timeoutSec)
}

func buildProgressbar(name string, size int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(fmt.Sprintf("downloading %s ", name)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerPadding: " ", BarStart: ">", BarEnd: "|"}),
		progressbar.OptionFullWidth(),
	)
	_ = bar.RenderBlank()
	return bar
}

func findChunkSize(localPath string) (int64, int64, error) {
	fi, err := os.Stat(localPath)
	if err != nil {
		return 0, 0, errors.Errorf("%s file not found", localPath)
	}

	var MB int64 = 1024 * 1024
	var GB = MB * 1024

	if fi.Size() < 5*GB {
		return fi.Size(), 5 * MB, nil
	}

	if fi.Size() < 100*GB {
		return fi.Size(), 25 * MB, nil
	}

	return fi.Size(), 50 * MB, nil
}
