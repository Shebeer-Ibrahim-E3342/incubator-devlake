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

func TestIssueRelationshipDataFlow(t *testing.T) {
	var plugin impl.Freshrelease
	dataflowTester := e2ehelper.NewDataFlowTester(t, "freshrelease", plugin)

	taskData := &tasks.FreshreleaseTaskData{
		Options: &tasks.FreshreleaseOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_freshrelease_api_issue_relationships.csv", "_raw_freshrelease_api_issues")
	// verify issue extraction
	dataflowTester.FlushTabler(&models.FreshreleaseIssueRelationship{})
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
	dataflowTester.Subtask(tasks.ExtractIssuesMeta, taskData)

	dataflowTester.VerifyTableWithOptions(&models.FreshreleaseIssueRelationship{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_freshrelease_issue_relationships.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	// verify issue conversion
	dataflowTester.FlushTabler(&models.FreshreleaseBoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueRelationship{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_freshrelease_board_issues_relations.csv", &models.FreshreleaseBoardIssue{})

	dataflowTester.Subtask(tasks.ConvertIssueRelationshipsMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&ticket.IssueRelationship{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_relationships.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
