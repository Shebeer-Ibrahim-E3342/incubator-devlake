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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type freshreleaseMultiAuth20240607 struct {
	AuthMethod string `gorm:"type:varchar(20)"`
	Token      string `gorm:"type:varchar(255)"`
}

func (freshreleaseMultiAuth20240607) TableName() string {
	return "_tool_freshrelease_connections"
}

type addFreshreleaseMultiAuth20240607 struct{}

func (script *addFreshreleaseMultiAuth20240607) Up(basicRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(basicRes, &freshreleaseMultiAuth20240607{})
	if err != nil {
		return err
	}
	return basicRes.GetDal().UpdateColumn(
		&freshreleaseMultiAuth20240607{},
		"auth_method", plugin.AUTH_METHOD_BASIC,
		dal.Where("auth_method IS NULL"),
	)
}

func (*addFreshreleaseMultiAuth20240607) Version() uint64 {
	return 20240607115901
}

func (*addFreshreleaseMultiAuth20240607) Name() string {
	return "add multiauth to _tool_freshrelease_connections"
}
