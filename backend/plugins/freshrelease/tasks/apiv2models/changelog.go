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

package apiv2models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/plugins/freshrelease/models"
)

type Changelog struct {
	ID      uint64             `json:"id,string"`
	Author  Account            `json:"author"`
	Created common.Iso8601Time `json:"created"`
	Items   []ChangelogItem    `json:"items"`
}

func (c Changelog) ToToolLayer(connectionId, issueId uint64, issueUpdated *time.Time) (*models.FreshreleaseIssueChangelogs, *models.FreshreleaseAccount) {
	changelog := &models.FreshreleaseIssueChangelogs{
		ConnectionId:      connectionId,
		ChangelogId:       c.ID,
		IssueId:           issueId,
		AuthorAccountId:   c.Author.getAccountId(),
		AuthorDisplayName: c.Author.DisplayName,
		AuthorActive:      c.Author.Active,
		Created:           c.Created.ToTime(),
		IssueUpdated:      issueUpdated,
	}
	return changelog, c.Author.ToToolLayer(connectionId)
}

type ChangelogItem struct {
	Field      string `json:"field"`
	Fieldtype  string `json:"fieldtype"`
	FromValue  string `json:"from"`
	FromString string `json:"fromString"`
	ToValue    string `json:"to"`
	ToString   string `json:"toString"`
}

func (c ChangelogItem) ToToolLayer(connectionId, changelogId uint64) *models.FreshreleaseIssueChangelogItems {
	item := &models.FreshreleaseIssueChangelogItems{
		ConnectionId: connectionId,
		ChangelogId:  changelogId,
		Field:        c.Field,
		FieldType:    c.Fieldtype,
		FromValue:    c.FromValue,
		FromString:   c.FromString,
		ToValue:      c.ToValue,
		ToString:     c.ToString,
	}
	return item
}

func (c ChangelogItem) ExtractUser(connectionId uint64) []*models.FreshreleaseAccount {
	if c.Field != "assignee" {
		return nil
	}
	var result []*models.FreshreleaseAccount
	if c.FromValue != "" {
		result = append(result, &models.FreshreleaseAccount{ConnectionId: connectionId, AccountId: c.FromValue})
	}
	if c.ToValue != "" {
		result = append(result, &models.FreshreleaseAccount{ConnectionId: connectionId, AccountId: c.ToValue})
	}
	return result
}
