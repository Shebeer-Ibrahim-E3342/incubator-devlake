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
	"github.com/apache/incubator-devlake/core/plugin"
)

// All return all the migration scripts
func All() []plugin.MigrationScript {
	return []plugin.MigrationScript{
		new(addSourceTable20240602),
		new(renameSourceTable20240603),
		new(addInitTables20240604),
		new(addTransformationRule20240605),
		new(addProjectName20240606),
		new(addFreshreleaseMultiAuth20240607),
		new(removeIssueStdStoryPoint),
		new(addCommitRepoPattern),
		new(expandRemotelinkUrl),
		new(addConnectionIdToTransformationRule),
		new(addChangeTotal20240612),
		new(expandRemotelinkSelfUrl),
		new(addDescAndComments),
		new(renameTr2ScopeConfig),
		new(addRepoUrl),
		new(addApplicationType),
		new(clearRepoPattern),
		new(addRawParamTableForScope),
		new(addIssueRelationship),
		new(dropIssueAllFields),
		new(modifyIssueRelationship),
		new(addComponents20240613),
		new(addFilterJQL),
		new(addWorklogToIssue),
		new(addSubtaskToIssue),
	}
}
