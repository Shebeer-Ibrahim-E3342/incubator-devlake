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

func TestIssueChangelogDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_issue_changelogs.csv", "_raw_freshrelease_api_issue_changelogs")
	dataflowTester.FlushTabler(&models.FreshreleaseIssueChangelogs{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueChangelogItems{})
	dataflowTester.FlushTabler(&models.FreshreleaseAccount{})
	dataflowTester.Subtask(tasks.ExtractIssueChangelogsMeta, taskData)
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueChangelogs{},
		"./snapshot_tables/_tool_freshrelease_issue_changelogs_extractor.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"changelog_id",
			"issue_id",
			"author_account_id",
			"author_display_name",
			"author_active",
			"created",
			"issue_updated",
		))
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueChangelogItems{},
		"./snapshot_tables/_tool_freshrelease_issue_changelog_items_extractor.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"changelog_id",
			"field",
			"field_type",
			"field_id",
			"from_value",
			"from_string",
			"to_value",
			"to_string",
		),
	)
	dataflowTester.VerifyTable(
		models.FreshreleaseAccount{},
		"./snapshot_tables/_tool_freshrelease_accounts_changelog.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"account_id",
			"account_type",
			"name",
			"email",
			"avatar_url",
			"timezone",
		),
	)
	// verify changelog conversion
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_issue_changelogs.csv", &models.FreshreleaseIssueChangelogs{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_issue_changelog_items.csv", &models.FreshreleaseIssueChangelogItems{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_statuses_for_changelog.csv", &models.FreshreleaseStatus{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_board_issues_for_changelog.csv", &models.FreshreleaseBoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertIssueChangelogsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueChangelogs{},
		"./snapshot_tables/issue_changelogs.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"issue_id",
			"author_id",
			"author_name",
			"field_id",
			"field_name",
			"original_from_value",
			"original_to_value",
			"from_value",
			"to_value",
			"created_date",
		),
	)
}
