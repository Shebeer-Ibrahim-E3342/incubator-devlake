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

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/freshrelease/impl"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks"
)

func TestRemotelinkDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
			ScopeConfig: &models.FreshreleaseScopeConfig{
				RemotelinkCommitShaPattern: ".*/commit/(.*)",
				RemotelinkRepoPattern: []models.CommitUrlPattern{
					{
						Pattern: "",
						Regex:   `https://example.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commits/(?P<commit_sha>\w{40})`,
					},
				},
			},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_remotelinks.csv", "_raw_freshrelease_api_remotelinks")

	// verify remotelink extraction
	dataflowTester.FlushTabler(&models.FreshreleaseRemotelink{})
	dataflowTester.FlushTabler(&models.FreshreleaseIssueCommit{})
	dataflowTester.Subtask(tasks.ExtractRemotelinksMeta, taskData)
	dataflowTester.VerifyTable(
		models.FreshreleaseRemotelink{},
		"./snapshot_tables/_tool_freshrelease_remotelinks.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"remotelink_id",
			"issue_id",
			"self",
			"title",
			"url",
			"issue_updated",
		),
	)
	dataflowTester.VerifyTable(
		models.FreshreleaseIssueCommit{},
		"./snapshot_tables/_tool_freshrelease_issue_commits.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"issue_id",
			"commit_sha",
			"commit_url",
		),
	)
}
