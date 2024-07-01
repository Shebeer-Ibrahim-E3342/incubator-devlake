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
	"testing"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/freshrelease/impl"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks"
)

func TestIssueDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
			ScopeConfig: &models.FreshreleaseScopeConfig{
				StoryPointField: "customfield_10024",
				TypeMappings: map[string]models.TypeMapping{
					"子任务": {
						StandardType: "Sub-task",
						StatusMappings: map[string]models.StatusMapping{
							"done": {StandardStatus: "你好世界"},
							"new":  {StandardStatus: "\u6069\u5E95\u6EF4\u68AF\u6B38\u592B\u5178\u4EA2\u59C6"},
						},
					},
					"任务": {
						StandardType: "Task",
						StatusMappings: map[string]models.StatusMapping{
							"done": {StandardStatus: "hello world"},
							"new":  {StandardStatus: "110 100 100 116 102 46 99 111 109"},
						},
					},
					// issueType "Test Execution" in raw_data and not fill here to test issueType not be defined
				},
			},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_issues.csv", "_raw_freshrelease_api_issues")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_issue_types.csv", "_raw_freshrelease_api_issue_types")

	// verify issue extraction
	dataflowTester.FlushTabler(&models.FreshreleaseIssue{})
	dataflowTester.FlushTabler(&models.FreshreleaseBoardIssue{})
	dataflowTester.FlushTabler(&models.FreshreleaseSprintIssue{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueComment{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueChangelogs{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueChangelogItems{})
	dataflowTester.FlushTabler(&models.FreshreleaseWorklog{})
	dataflowTester.FlushTabler(&models.FreshreleaseAccount{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueType{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueLabel{})
	dataflowTester.Subtask(tasks.ExtractIssueTypesMeta, taskData)
	dataflowTester.Subtask(tasks.ExtractIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueType{},
		"./snapshot_tables/_tool_freshrelease_issue_types.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"self",
			"id",
			"description",
			"icon_url",
			"name",
			"untranslated_name",
			"subtask",
			"avatar_id",
			"hierarchy_level",
		),
	)

	dataflowTester.VerifyTable(
		models.FreshreleaseIssue{},
		"./snapshot_tables/_tool_freshrelease_issues.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"issue_id",
			"project_id",
			"project_name",
			"self",
			"issue_key",
			"summary",
			"description",
			"type",
			"epic_key",
			"status_name",
			"status_key",
			"story_point",
			"original_estimate_minutes",
			"aggregate_estimate_minutes",
			"remaining_estimate_minutes",
			"creator_account_id",
			"creator_account_type",
			"creator_display_name",
			"assignee_account_id",
			"assignee_account_type",
			"assignee_display_name",
			"priority_id",
			"priority_name",
			"parent_id",
			"parent_key",
			"sprint_id",
			"sprint_name",
			"resolution_date",
			"created",
			"updated",
			"spent_minutes",
			"lead_time_minutes",
			"std_type",
			"std_status",
			"icon_url",
			"changelog_total",
			"comment_total",
		),
	)

	dataflowTester.VerifyTable(
		models.FreshreleaseBoardIssue{},
		"./snapshot_tables/_tool_freshrelease_board_issues.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"board_id",
			"issue_id",
		),
	)
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueChangelogs{},
		"./snapshot_tables/_tool_freshrelease_issue_changelogs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"changelog_id",
			"issue_id",
			"author_account_id",
			"author_display_name",
			"author_active",
			"created",
			"issue_updated",
		),
	)
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueChangelogItems{},
		"./snapshot_tables/_tool_freshrelease_issue_changelog_items.csv",
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
		models.FreshreleaseSprintIssue{},
		"./snapshot_tables/_tool_freshrelease_sprint_issues.csv",
		[]string{
			"connection_id",
			"sprint_id",
			"issue_id",
			"resolution_date",
			"issue_created_date",
		},
	)

	dataflowTester.VerifyTableWithRawData(
		models.FreshreleaseIssueLabel{},
		"./snapshot_tables/_tool_freshrelease_issue_labels.csv",
		[]string{
			"connection_id",
			"issue_id",
			"label_name",
		})

	// verify issue conversion
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueAssignee{})
	dataflowTester.Subtask(tasks.ConvertIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		"./snapshot_tables/issues.csv",
		[]string{
			"id",
			"url",
			"icon_url",
			"issue_key",
			"title",
			"description",
			"epic_key",
			"type",
			"original_type",
			"status",
			"original_status",
			"story_point",
			"resolution_date",
			"created_date",
			"updated_date",
			"lead_time_minutes",
			"parent_issue_id",
			"priority",
			"original_estimate_minutes",
			"time_spent_minutes",
			"time_remaining_minutes",
			"creator_id",
			"creator_name",
			"assignee_id",
			"assignee_name",
			"severity",
			"component",
			"original_project",
		},
	)
	dataflowTester.VerifyTable(
		ticket.BoardIssue{},
		"./snapshot_tables/board_issues.csv",
		[]string{"board_id", "issue_id"},
	)
	dataflowTester.VerifyTableWithOptions(ticket.IssueAssignee{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_assignees.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
