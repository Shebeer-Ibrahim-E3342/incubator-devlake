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

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/freshrelease/impl"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
	"github.com/apache/incubator-devlake/plugins/freshrelease/tasks"
)

func TestConvertIssueRepoCommitsDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
			ScopeConfig: &models.FreshreleaseScopeConfig{
				RemotelinkCommitShaPattern: `.*/commit/(.*)`,
				RemotelinkRepoPattern: []models.CommitUrlPattern{
					{"", `https://bitbucket.org/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commits/(?P<commit_sha>\w{40})`},
					{"", `https://gitlab.com/(?P<namespace>\S+)/(?P<repo_name>\S+)/-/commit/(?P<commit_sha>\w{40})`},
					{"", `https://github.com/(?P<namespace>[^/]+)/(?P<repo_name>[^/]+)/commit/(?P<commit_sha>\w{40})`},
				},
			},
		},
	}
	dataflowTester.FlushTabler(&crossdomain.IssueRepoCommit{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_issue_commits_for_ConvertIssueRepoCommits.csv", &models.FreshreleaseIssueCommit{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_board_issues_for_ConvertIssueRepoCommits.csv", &models.FreshreleaseBoardIssue{})
	dataflowTester.Subtask(tasks.ConvertIssueRepoCommitsMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.IssueRepoCommit{},
		"./snapshot_tables/issue_repo_commits.csv",
		e2ehelper.ColumnWithRawData(
			"issue_id",
			"repo_url",
			"commit_sha",
			"host",
			"namespace",
			"repo_name",
		),
	)
}
