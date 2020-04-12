/*
 * Copyright 2018 The Service Manager Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package postgres

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type BaseEntity struct {
	ID             string    `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	PagingSequence int64     `db:"paging_sequence,auto_increment"`
	Ready          bool      `db:"ready"`
}

func (e *BaseEntity) GetID() string {
	return e.ID
}

type BaseLabelEntity struct {
	ID        sql.NullString `db:"id"`
	Key       sql.NullString `db:"key"`
	Val       sql.NullString `db:"val"`
	CreatedAt pq.NullTime    `db:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at"`
}

func (el BaseLabelEntity) GetKey() string {
	return el.Key.String
}

func (el BaseLabelEntity) GetValue() string {
	return el.Val.String
}

func (el BaseLabelEntity) LabelsPrimaryColumn() string {
	return "id"
}
