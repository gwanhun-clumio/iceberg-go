// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gwanhun-clumio/iceberg-go"
	"github.com/gwanhun-clumio/iceberg-go/io"
	"github.com/gwanhun-clumio/iceberg-go/table"
)

type CreateTableCfg struct {
	Location      string
	PartitionSpec *iceberg.PartitionSpec
	SortOrder     table.SortOrder
	Properties    iceberg.Properties
}

func GetMetadataLoc(location string, newVersion uint) string {
	return fmt.Sprintf("%s/metadata/%05d-%s.metadata.json",
		location, newVersion, uuid.New().String())
}

func WriteMetadata(ctx context.Context, metadata table.Metadata, loc string, props iceberg.Properties) error {
	fs, err := io.LoadFS(ctx, props, loc)
	if err != nil {
		return err
	}

	wfs, ok := fs.(io.WriteFileIO)
	if !ok {
		return errors.New("filesystem IO does not support writing")
	}

	out, err := wfs.Create(loc)
	if err != nil {
		return nil
	}

	defer out.Close()

	return json.NewEncoder(out).Encode(metadata)
}
