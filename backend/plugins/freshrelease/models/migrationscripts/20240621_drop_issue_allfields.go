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

type dropIssueAllFields struct{}

func (script *dropIssueAllFields) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().DropColumns("_tool_freshrelease_issues", "all_fields")
}

func (*dropIssueAllFields) Version() uint64 {
	return 20240621162534
}

func (*dropIssueAllFields) Name() string {
	return "add table _tool_freshrelease_issues.all_fields column"
}