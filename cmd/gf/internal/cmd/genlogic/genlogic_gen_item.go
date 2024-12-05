// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genlogic

type (
	CGenLogicInternalGenItems struct {
		index int
		Items []CGenLogicInternalGenItem
	}
	CGenLogicInternalGenItem struct {
		Clear              bool
		StorageDirPaths    []string
		GeneratedFilePaths []string
	}
)

func newCGenLogicInternalGenItems() *CGenLogicInternalGenItems {
	return &CGenLogicInternalGenItems{
		index: -1,
		Items: make([]CGenLogicInternalGenItem, 0),
	}
}

func (i *CGenLogicInternalGenItems) Scale() {
	i.Items = append(i.Items, CGenLogicInternalGenItem{
		StorageDirPaths:    make([]string, 0),
		GeneratedFilePaths: make([]string, 0),
		Clear:              false,
	})
	i.index++
}

func (i *CGenLogicInternalGenItems) SetClear(clear bool) {
	i.Items[i.index].Clear = clear
}

func (i CGenLogicInternalGenItems) AppendDirPath(storageDirPath string) {
	i.Items[i.index].StorageDirPaths = append(
		i.Items[i.index].StorageDirPaths,
		storageDirPath,
	)
}

func (i CGenLogicInternalGenItems) AppendGeneratedFilePath(generatedFilePath string) {
	i.Items[i.index].GeneratedFilePaths = append(
		i.Items[i.index].GeneratedFilePaths,
		generatedFilePath,
	)
}
