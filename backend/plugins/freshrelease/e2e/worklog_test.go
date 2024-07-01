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

package e2e

import (
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/freshrelease/impl"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks"
	"testing"
)

func TestWorklogDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_worklogs.csv", "_raw_freshrelease_api_worklogs")

	// verify worklog extraction
	dataflowTester.FlushTabler(&models.FreshreleaseWorklog{})
	dataflowTester.Subtask(tasks.ExtractWorklogsMeta, taskData)
	dataflowTester.VerifyTable(
		models.FreshreleaseWorklog{},
		"./snapshot_tables/_tool_freshrelease_worklogs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"issue_id",
			"worklog_id",
			"author_id",
			"update_author_id",
			"time_spent",
			"time_spent_seconds",
			"updated",
			"started",
			"issue_updated",
		),
	)

	// verify worklog conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_worklogs_for_convertor.csv", &models.FreshreleaseWorklog{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_board_issues_for_worklog.csv", &models.FreshreleaseBoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueWorklog{})
	dataflowTester.Subtask(tasks.ConvertWorklogsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueWorklog{},
		"./snapshot_tables/worklogs.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"author_id",
			"comment",
			"time_spent_minutes",
			"logged_date",
			"started_date",
			"issue_id",
		),
	)
}
