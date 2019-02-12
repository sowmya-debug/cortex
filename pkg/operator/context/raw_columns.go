/*
Copyright 2019 Cortex Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package context

import (
	"bytes"

	"github.com/cortexlabs/cortex/pkg/api/context"
	"github.com/cortexlabs/cortex/pkg/api/resource"
	s "github.com/cortexlabs/cortex/pkg/api/strings"
	"github.com/cortexlabs/cortex/pkg/api/userconfig"
	"github.com/cortexlabs/cortex/pkg/utils/errors"
	"github.com/cortexlabs/cortex/pkg/utils/util"
)

func getRawColumns(
	config *userconfig.Config,
	env *context.Environment,
) (context.RawColumns, error) {

	rawColumns := context.RawColumns{}

	for _, columnConfig := range config.RawColumns {
		var buf bytes.Buffer
		buf.WriteString(env.ID)
		buf.WriteString(columnConfig.GetName())
		buf.WriteString(columnConfig.GetType())

		var rawColumn context.RawColumn
		switch typedColumnConfig := columnConfig.(type) {
		case *userconfig.RawIntColumn:
			buf.WriteString(s.Bool(typedColumnConfig.Required))
			buf.WriteString(s.Obj(typedColumnConfig.Min))
			buf.WriteString(s.Obj(typedColumnConfig.Max))
			buf.WriteString(s.Obj(util.SortInt64sCopy(typedColumnConfig.Values)))
			id := util.HashBytes(buf.Bytes())
			idWithTags := util.HashStr(id + typedColumnConfig.Tags.ID())
			rawColumn = &context.RawIntColumn{
				ComputedResourceFields: &context.ComputedResourceFields{
					ResourceFields: &context.ResourceFields{
						ID:           id,
						IDWithTags:   idWithTags,
						ResourceType: resource.RawColumnType,
					},
				},
				RawIntColumn: typedColumnConfig,
			}
		case *userconfig.RawFloatColumn:
			buf.WriteString(s.Bool(typedColumnConfig.Required))
			buf.WriteString(s.Obj(typedColumnConfig.Min))
			buf.WriteString(s.Obj(typedColumnConfig.Max))
			buf.WriteString(s.Obj(util.SortFloat32sCopy(typedColumnConfig.Values)))
			id := util.HashBytes(buf.Bytes())
			idWithTags := util.HashStr(id + typedColumnConfig.Tags.ID())
			rawColumn = &context.RawFloatColumn{
				ComputedResourceFields: &context.ComputedResourceFields{
					ResourceFields: &context.ResourceFields{
						ID:           id,
						IDWithTags:   idWithTags,
						ResourceType: resource.RawColumnType,
					},
				},
				RawFloatColumn: typedColumnConfig,
			}
		case *userconfig.RawStringColumn:
			buf.WriteString(s.Bool(typedColumnConfig.Required))
			buf.WriteString(s.Obj(util.SortStrsCopy(typedColumnConfig.Values)))
			id := util.HashBytes(buf.Bytes())
			idWithTags := util.HashStr(id + typedColumnConfig.Tags.ID())
			rawColumn = &context.RawStringColumn{
				ComputedResourceFields: &context.ComputedResourceFields{
					ResourceFields: &context.ResourceFields{
						ID:           id,
						IDWithTags:   idWithTags,
						ResourceType: resource.RawColumnType,
					},
				},
				RawStringColumn: typedColumnConfig,
			}
		default:
			return nil, errors.New(userconfig.Identify(columnConfig), s.ErrInvalidStr(userconfig.TypeKey, userconfig.IntegerColumnType.String(), userconfig.FloatColumnType.String(), userconfig.StringColumnType.String())) // unexpected error
		}

		rawColumns[columnConfig.GetName()] = rawColumn
	}

	return rawColumns, nil
}