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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"reflect"
)

var ConvertWorklogsMeta = plugin.SubTaskMeta{
	Name:             "convertWorklogs",
	EntryPoint:       ConvertWorklogs,
	EnabledByDefault: true,
	Description:      "convert Freshrelease work logs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FreshreleaseTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert worklog")
	// select all worklogs belongs to the board
	clauses := []dal.Clause{
		dal.From(&models.FreshreleaseWorklog{}),
		dal.Select("_tool_freshrelease_worklogs.*"),
		dal.Join(`LEFT JOIN _tool_freshrelease_board_issues
              ON _tool_freshrelease_board_issues.connection_id = _tool_freshrelease_worklogs.connection_id
                   AND _tool_freshrelease_board_issues.issue_id = _tool_freshrelease_worklogs.issue_id`),
		dal.Where("_tool_freshrelease_board_issues.connection_id = ? AND _tool_freshrelease_board_issues.board_id = ?", connectionId, boardId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "convert worklog error")
		return err
	}
	defer cursor.Close()

	worklogIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseWorklog{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseAccount{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FreshreleaseApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.FreshreleaseWorklog{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			freshreleaseWorklog := inputRow.(*models.FreshreleaseWorklog)
			worklog := &ticket.IssueWorklog{
				DomainEntity:     domainlayer.DomainEntity{Id: worklogIdGen.Generate(freshreleaseWorklog.ConnectionId, freshreleaseWorklog.IssueId, freshreleaseWorklog.WorklogId)},
				IssueId:          issueIdGen.Generate(freshreleaseWorklog.ConnectionId, freshreleaseWorklog.IssueId),
				TimeSpentMinutes: freshreleaseWorklog.TimeSpentSeconds / 60,
				StartedDate:      &freshreleaseWorklog.Started,
				LoggedDate:       &freshreleaseWorklog.Updated,
			}
			if freshreleaseWorklog.AuthorId != "" {
				worklog.AuthorId = accountIdGen.Generate(connectionId, freshreleaseWorklog.AuthorId)
			}
			return []interface{}{worklog}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
