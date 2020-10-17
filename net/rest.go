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
package net

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/veeva/vvfst/config"
	"github.com/veeva/vvfst/model"
	"github.com/veeva/vvfst/vlog"
	"strings"
	"sync"
	"time"
)

var (
	restClient *RestClient
	once       sync.Once
)

type RestClient struct {
	client *resty.Client
}

func InitRestClient(enableDebug bool) *RestClient {
	once.Do(func() {
		vlog.Debugf("Initialize the rest client with enableDebug: %t", enableDebug)

		remoteURL := fmt.Sprintf("https://%s/api/%s", config.DomainName(), config.APIVersion())
		restClient = NewRestClient(enableDebug, remoteURL)
	})

	return restClient
}

func NewRestClient(enableDebug bool, url string) *RestClient {
	client := resty.New()
	client.SetDebug(enableDebug)

	client.SetHostURL(url)
	client.SetHeader("User-Agent", "vvfst/20.2")

	return &RestClient{client: client}
}

func (rc *RestClient) BuildRestRequest(includeAuth bool) *resty.Request {
	if rc.client == nil {
		vlog.Fatal("Initialize the rest client")
	}

	req := rc.client.R()
	if includeAuth {
		req = req.SetAuthToken(config.AuthResult().SessionID)
	}

	return req
}

func LogTime(msg string, resp *resty.Response) {
	vlog.Infof("[Duration: %.3f seconds] %s ", float32(resp.Time())/float32(time.Second), msg)
}

func LogRestResultError(msg string, resp *model.RestResultError) {
	vlog.Errorf(FormatRestResultError(msg, resp))
}

func FormatRestResultError(msg string, resp *model.RestResultError) string {
	if msg == "" {
		return fmt.Sprintf("[%s]: %s ", resp.Type, resp.Message)
	}
	return fmt.Sprintf("%s - [%s]: %s ", msg, resp.Type, resp.Message)
}

func IsSessionExpired(err error) bool {
	return err != nil && strings.Contains(err.Error(), "INVALID_SESSION_ID")
}
