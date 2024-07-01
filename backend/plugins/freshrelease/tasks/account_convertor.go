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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"reflect"
)

var ConvertAccountsMeta = plugin.SubTaskMeta{
	Name:             "convertAccounts",
	EntryPoint:       ConvertAccounts,
	EnabledByDefault: true,
	Description:      "convert Freshrelease accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertAccounts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FreshreleaseTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert account")
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.FreshreleaseAccount{}),
		dal.Where("account_id != ? AND connection_id = ?", "", connectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseAccount{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FreshreleaseApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_USERS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.FreshreleaseAccount{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			freshreleaseAccount := inputRow.(*models.FreshreleaseAccount)
			u := &crossdomain.Account{
				DomainEntity: domainlayer.DomainEntity{
					Id: accountIdGen.Generate(connectionId, freshreleaseAccount.AccountId),
				},
				FullName:  freshreleaseAccount.Name,
				UserName:  freshreleaseAccount.Name,
				Email:     freshreleaseAccount.Email,
				AvatarUrl: freshreleaseAccount.AvatarUrl,
			}
			return []interface{}{u}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
