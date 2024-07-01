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

type freshreleaseIssue20240626 struct {
	Subtask bool
}

func (freshreleaseIssue20240626) TableName() string {
	return "_tool_freshrelease_issues"
}

type addSubtaskToIssue struct{}

func (script *addSubtaskToIssue) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(&freshreleaseIssue20240626{})
	if err != nil {
		return err
	}
	// force full issue extraction so issue.worklog_total can be updated
	return db.Exec("DELETE FROM _devlake_subtask_states WHERE plugin = ? AND subtask = ?", "freshrelease", "extractIssues")
}

func (*addSubtaskToIssue) Version() uint64 {
	return 20240625145131
}

func (*addSubtaskToIssue) Name() string {
	return "add subtask to _tool_freshrelease_issues"
}
