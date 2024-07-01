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

package tasks

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
)

type FreshreleaseOptions struct {
	ConnectionId  uint64                  `json:"connectionId" mapstructure:"connectionId"`
	BoardId       uint64                  `json:"boardId" mapstructure:"boardId"`
	ScopeConfig   *models.FreshreleaseScopeConfig `json:"scopeConfig" mapstructure:"scopeConfig"`
	ScopeConfigId uint64                  `json:"scopeConfigId" mapstructure:"scopeConfigId"`
	PageSize      int                     `json:"pageSize" mapstructure:"pageSize"`
}

type FreshreleaseTaskData struct {
	Options        *FreshreleaseOptions
	ApiClient      *api.ApiAsyncClient
	FreshreleaseServerInfo models.FreshreleaseServerInfo
}

type FreshreleaseApiParams models.FreshreleaseApiParams

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*FreshreleaseOptions, errors.Error) {
	var op FreshreleaseOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid connectionId:%d", op.ConnectionId))
	}
	if op.BoardId == 0 {
		return nil, errors.BadInput.New(fmt.Sprintf("invalid boardId:%d", op.BoardId))
	}
	return &op, nil
}
