/*
 * Minio Cloud Storage, (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"fmt"

	"encoding/json"

	"github.com/Sirupsen/logrus"
)

type httpNotify struct {
	Enable   bool   `json:"enable"`
	Endpoint string `json:"endpoint"`
}

type httpConn struct {
	HTTPClient *http.Client
	Endpoint   string
}

func newHTTPNotify(accountID string) (*logrus.Logger, error) {
	rNotify := serverConfig.GetHTTPNotifyByID(accountID)
	fmt.Println(rNotify)

	connection := httpConn{
		Endpoint: rNotify.Endpoint,
	}

	// Configure aggressive timeouts for client posts.
	connection.HTTPClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	notifyLog := logrus.New()
	notifyLog.Out = ioutil.Discard

	// Set default JSON formatter.
	notifyLog.Formatter = new(logrus.JSONFormatter)

	notifyLog.Hooks.Add(connection)

	// Success
	return notifyLog, nil
}

// Fire is called when an event should be sent to the message broker.
func (n httpConn) Fire(entry *logrus.Entry) error {

	// Fetch event type upon reflecting on its original type.
	entryStr, ok := entry.Data["EventType"].(string)
	if !ok {
		return nil
	}
	fmt.Println("We got a HTTP event fired.")
	fmt.Println(entryStr)

	httpPostBody1 := httpPostBody{
		Records: entry.Data["Records"],
		Key:     entry.Data["Key"].(string),
	}

	itemsStr, err := json.Marshal(httpPostBody1)
	if err != nil {
		entry.Warn("Cannot convert request body to JSON.")
	}
	buf := bytes.NewBuffer(itemsStr)
	response, err := n.HTTPClient.Post(n.Endpoint, "application/json", buf)

	fmt.Println(response.Status, err)

	return nil
}

type httpPostBody struct {
	Records interface{}
	Key     string
}

// Levels are Required for logrus hook implementation
func (httpConn) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
	}
}
