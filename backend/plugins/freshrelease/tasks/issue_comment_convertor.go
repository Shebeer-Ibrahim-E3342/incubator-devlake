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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
)

var ConvertIssueCommentsMeta = plugin.SubTaskMeta{
	Name:             "ConvertIssueComments",
	EntryPoint:       ConvertIssueComments,
	EnabledByDefault: false,
	Description:      "convert Freshrelease issue comments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConvertIssueComments(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FreshreleaseTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert issue comments")

	clauses := []dal.Clause{
		dal.Select("jic.*"),
		dal.From("_tool_freshrelease_issue_comments jic"),
		dal.Join(`left join _tool_freshrelease_board_issues jbi on (
			jbi.connection_id = jic.connection_id
			AND jbi.issue_id = jic.issue_id
		)`),
		dal.Where("jbi.connection_id = ? AND jbi.board_id = ?", connectionId, boardId),
		dal.Orderby("jbi.connection_id, jbi.issue_id"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.FreshreleaseAccount{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FreshreleaseApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_ISSUE_COMMENT_TABLE,
		},
		InputRowType: reflect.TypeOf(models.FreshreleaseIssueComment{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			var result []interface{}
			issueComment := inputRow.(*models.FreshreleaseIssueComment)
			domainIssueComment := &ticket.IssueComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(data.Options.ConnectionId, issueComment.IssueId),
				},
				IssueId:     issueIdGen.Generate(data.Options.ConnectionId, issueComment.IssueId),
				Body:        issueComment.Body,
				AccountId:   accountIdGen.Generate(data.Options.ConnectionId, issueComment.CreatorAccountId),
				CreatedDate: issueComment.Created,
			}
			if !issueComment.Updated.IsZero() {
				domainIssueComment.UpdatedDate = &issueComment.Updated
			}
			result = append(result, domainIssueComment)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
