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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type freshreleaseIssue20240606 struct {
	ProjectName string `gorm:"type:varchar(255)"`
}

func (freshreleaseIssue20240606) TableName() string {
	return "_tool_freshrelease_issues"
}

type addProjectName20240606 struct{}

func (script *addProjectName20240606) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &freshreleaseIssue20240606{})
}

func (*addProjectName20240606) Version() uint64 {
	return 20240606142316
}

func (*addProjectName20240606) Name() string {
	return "add project_name to _tool_freshrelease_issues"
}
