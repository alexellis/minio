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
	"io/ioutil"

	"fmt"

	"github.com/Sirupsen/logrus"
)

type httpNotify struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"address"`
}

func newHTTPNotify(accountID string) (*logrus.Logger, error) {
	rNotify := serverConfig.GetHTTPNotifyByID(accountID)
	fmt.Println(rNotify)

	notifyLog := logrus.New()

	notifyLog.Out = ioutil.Discard

	// Set default JSON formatter.
	notifyLog.Formatter = new(logrus.JSONFormatter)

	// Success
	return notifyLog, nil
}

// Fire is called when an event should be sent to the message broker.
func Fire(entry *logrus.Entry) error {

	// Fetch event type upon reflecting on its original type.
	entryStr, ok := entry.Data["EventType"].(string)
	if !ok {
		return nil
	}
	fmt.Println("We got a HTTP event fired.")
	fmt.Println(entryStr)

	// buf := bytes.NewBufferString(entryStr)
	// response, err := http.Post("http://requestb.in/1i9al7m1", "application/json", buf)

	// fmt.Println(response.Status, err)

	return nil
}

// Levels are Required for logrus hook implementation
func Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
	}
}
