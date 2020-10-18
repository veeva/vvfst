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
package model

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

// Response Status Type to represent FAILURE or SUCCESS
type ResponseStatusType string

const (
	FAILURE ResponseStatusType = "FAILURE"
	SUCCESS ResponseStatusType = "SUCCESS"
)

func (rst *ResponseStatusType) UnmarshalJSON(b []byte) error {
	var s string
	_ = json.Unmarshal(b, &s)
	rsType := ResponseStatusType(s)
	switch rsType {
	case FAILURE, SUCCESS:
		return nil
	}
	return errors.New("Invalid ResponseStatus type")
}

func (rst *ResponseStatusType) IsSuccess() bool {
	return *rst == SUCCESS
}

type RestResult struct {
	ResponseStatus ResponseStatusType `json:"responseStatus"`
	Errors         []*RestResultError `json:"errors"`
}

type RestResultError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ResponseDetails struct {
	NextPage string `json:"next_page"`
}

type Item struct {
	Path         string     `json:"path"`
	Name         string     `json:"name"`
	Kind         string     `json:"kind"`
	Size         int64      `json:"size"`
	ModifiedDate *time.Time `json:"modified_date"`
	MD5          string     `json:"file_content_md5"`
}

type ItemsRestResult struct {
	RestResult
	ResponseDetails *ResponseDetails `json:"responseDetails"`
	Data            []*Item          `json:"data"`
}

type ItemRestResult struct {
	RestResult
	Data *Item `json:"data"`
}

type UploadSession struct {
	UploadSessionID    string     `json:"id"`
	Path               string     `json:"path"`
	Name               string     `json:"name"`
	Size               int64      `json:"size"`
	UploadedSize       int64      `json:"uploaded"`
	UploadedPartsCount int        `json:"uploaded_parts"`
	CreatedDate        *time.Time `json:"created_date"`
	ExpirationDate     *time.Time `json:"expiration_date"`
	LastUploadedDate   *time.Time `json:"last_uploaded_date"`
}

type UploadSessionsRestResult struct {
	RestResult
	Data []*UploadSession `json:"data"`
}

type UploadSessionRestResult struct {
	RestResult
	Data *UploadSession `json:"data"`
}

type UploadPart struct {
	PartNumber     int    `json:"part_number"`
	PartSize       int64  `json:"size"`
	PartContentMD5 string `json:"part_content_md5"`
}

type UploadPartRestResult struct {
	RestResult
	Data *UploadPart `json:"data"`
}

type UploadPartsRestResult struct {
	RestResult
	Data []*UploadPart `json:"data"`
}

type Job struct {
	JobID int64  `json:"job_id"`
	URL   string `json:"url"`
}

type JobRestResult struct {
	RestResult
	Data *Job `json:"data"`
}

// DownloadItem - Download from RemotePath or Remote Href, either one should be available.
type DownloadItem struct {
	RemoteHref string
	RemotePath string
	Size       int64
	LocalPath  string
}

type UploadItem struct {
	RemotePath string
	LocalPath  string
}

type JobStatusData struct {
	Status string  `json:"status"`
	Links  []*Link `json:"links"`
}

type Link struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method"`
	Accept string `json:"accept"`
}

type JobStatusRestResult struct {
	RestResult
	Data *JobStatusData
}
