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
	"fmt"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks/apiv2models"
)

type FreshreleaseRemotePagination struct {
	MaxResults int `json:"maxResults"`
	StartAt    int `json:"startAt"`
}

func queryFreshreleaseAgileBoards(
	apiClient plugin.ApiClient,
	keyword string,
	page FreshreleaseRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.FreshreleaseBoard],
	nextPage *FreshreleaseRemotePagination,
	err errors.Error,
) {
	if page.MaxResults == 0 {
		page.MaxResults = 50
	}
	res, err := apiClient.Get("agile/1.0/board", url.Values{
		"maxResults": {fmt.Sprintf("%v", page.MaxResults)},
		"startAt":    {fmt.Sprintf("%v", page.StartAt)},
		"name":       {keyword},
	}, nil)
	if err != nil {
		return
	}

	resBody := struct {
		FreshreleaseRemotePagination
		IsLast bool                `json:"isLast"`
		Values []apiv2models.Board `json:"values"`
	}{}

	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return
	}

	for _, board := range resBody.Values {
		children = append(children, dsmodels.DsRemoteApiScopeListEntry[models.FreshreleaseBoard]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			Id:       fmt.Sprintf("%v", board.ID),
			ParentId: nil,
			Name:     board.Name,
			FullName: board.Name,
			Data:     board.ToToolLayer(0),
		})
	}

	if !resBody.IsLast {
		nextPage = &FreshreleaseRemotePagination{
			MaxResults: page.MaxResults,
			StartAt:    page.StartAt + page.MaxResults,
		}
	}

	return
}

func listFreshreleaseRemoteScopes(
	_ *models.FreshreleaseConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page FreshreleaseRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.FreshreleaseBoard],
	nextPage *FreshreleaseRemotePagination,
	err errors.Error,
) {
	return queryFreshreleaseAgileBoards(apiClient, "", page)
}

// RemoteScopes list all available scopes on the remote server
// @Summary list all available scopes on the remote server
// @Description list all available scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.FreshreleaseBoard]
// @Tags plugins/freshrelease
// @Router /plugins/freshrelease/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

func searchFreshreleaseRemoteBoards(
	apiClient plugin.ApiClient,
	params *dsmodels.DsRemoteApiScopeSearchParams,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.FreshreleaseBoard],
	err errors.Error,
) {
	if params.Page == 0 {
		params.Page = 1
	}
	page := FreshreleaseRemotePagination{
		MaxResults: params.PageSize,
		StartAt:    (params.Page - 1) * params.PageSize,
	}
	children, _, err = queryFreshreleaseAgileBoards(apiClient, params.Search, page)
	return
}

// SearchRemoteScopes searches scopes on the remote server
// @Summary searches scopes on the remote server
// @Description searches scopes on the remote server
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.FreshreleaseBoard] "the parentIds are always null"
// @Tags plugins/freshrelease
// @Router /plugins/freshrelease/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeSearch.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Router /plugins/freshrelease/connections/{connectionId}/proxy/{path} [GET]
// @Tags plugins/freshrelease
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
