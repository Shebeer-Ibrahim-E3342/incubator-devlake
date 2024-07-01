/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"github.com/apache/incubator-devlake/core/utils"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type EpicResponse struct {
	Id    int
	Title string
	Value string
}

type BoardResponse struct {
	Id    int
	Title string
	Value string
}

// FreshreleaseConn holds the essential information to connect to the Freshrelease API
type FreshreleaseConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.MultiAuth      `mapstructure:",squash"`
	helper.BasicAuth      `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

func (jc *FreshreleaseConn) Sanitize() FreshreleaseConn {
	jc.Password = ""
	jc.AccessToken.Token = utils.SanitizeString(jc.AccessToken.Token)
	return *jc
}

// SetupAuthentication implements the `IAuthentication` interface by delegating
// the actual logic to the `MultiAuth` struct to help us write less code
func (jc *FreshreleaseConn) SetupAuthentication(req *http.Request) errors.Error {
	return jc.MultiAuth.SetupAuthenticationForConnection(jc, req)
}

// FreshreleaseConnection holds FreshreleaseConn plus ID/Name for database storage
type FreshreleaseConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	FreshreleaseConn              `mapstructure:",squash"`
}

func (FreshreleaseConnection) TableName() string {
	return "_tool_freshrelease_connections"
}

func (connection *FreshreleaseConnection) MergeFromRequest(target *FreshreleaseConnection, body map[string]interface{}) error {
	token := target.Token
	password := target.Password
	authMethod := target.AuthMethod

	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}

	modifiedToken := target.Token
	modifiedPassword := target.Password
	modifiedAuthMethod := target.AuthMethod

	// maybe auth method has changed
	if authMethod == modifiedAuthMethod {
		if modifiedToken == "" || modifiedToken == utils.SanitizeString(token) {
			target.Token = token
		}
		if modifiedPassword == "" || modifiedPassword == utils.SanitizeString(password) {
			target.Password = password
		}
	}

	return nil
}

func (connection FreshreleaseConnection) Sanitize() FreshreleaseConnection {
	connection.FreshreleaseConn = connection.FreshreleaseConn.Sanitize()
	return connection
}
