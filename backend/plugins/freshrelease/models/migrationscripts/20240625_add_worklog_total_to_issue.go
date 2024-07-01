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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type freshreleaseIssue20240625 struct {
	WorklogTotal int
}

func (freshreleaseIssue20240625) TableName() string {
	return "_tool_freshrelease_issues"
}

type addWorklogToIssue struct{}

func (script *addWorklogToIssue) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(&freshreleaseIssue20240625{})
	if err != nil {
		return err
	}
	// force full issue extraction so issue.worklog_total can be updated
	err = db.Exec("DELETE FROM _devlake_subtask_states WHERE plugin = ? AND subtask = ?", "freshrelease", "extractIssues")
	if err != nil {
		return err
	}
	// force full collection for all freshrelease worklogs
	return db.Exec("DELETE FROM _devlake_collector_latest_state WHERE raw_data_table = ?", "_raw_freshrelease_api_worklogs")
}

func (*addWorklogToIssue) Version() uint64 {
	return 20240625145131
}

func (*addWorklogToIssue) Name() string {
	return "add worklog_total to _tool_freshrelease_issues"
}
