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

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks/apiv2models"
)

type PutScopesReqBody api.PutScopesReqBody[models.FreshreleaseBoard]
type ScopeDetail api.ScopeDetail[models.FreshreleaseBoard, models.FreshreleaseScopeConfig]

// PutScope create or update freshrelease board
// @Summary create or update freshrelease board
// @Description Create or update Freshrelease board
// @Tags plugins/freshrelease
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param searchTerm query string false "search term for scope name"
// @Param scope body PutScopesReqBody true "json"
// @Success 200  {object} []models.FreshreleaseBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/freshrelease/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.PutMultiple(input)
}

// UpdateScope patch to freshrelease board
// @Summary patch to freshrelease board
// @Description patch to freshrelease board
// @Tags plugins/freshrelease
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Param scope body models.FreshreleaseBoard true "json"
// @Success 200  {object} models.FreshreleaseBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/freshrelease/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Patch(input)
}

// GetScopeList get Freshrelease boards
// @Summary get Freshrelease boards
// @Description get Freshrelease boards
// @Tags plugins/freshrelease
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/freshrelease/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetPage(input)
}

// GetScope get one Freshrelease board
// @Summary get one Freshrelease board
// @Description get one Freshrelease board
// @Tags plugins/freshrelease
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Success 200  {object} ScopeDetail
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/freshrelease/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.GetScopeDetail(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/freshrelease
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/freshrelease/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ScopeApi.Delete(input)
}

func GetApiFreshrelease(op *tasks.FreshreleaseOptions, apiClient plugin.ApiClient) (*apiv2models.Board, errors.Error) {
	boardRes := &apiv2models.Board{}
	res, err := apiClient.Get(fmt.Sprintf("agile/1.0/board/%d", op.BoardId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting repo detail from %s", res.Request.URL.String()))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}
	err = errors.Convert(json.Unmarshal(body, boardRes))
	if err != nil {
		return nil, err
	}
	return boardRes, nil
}
