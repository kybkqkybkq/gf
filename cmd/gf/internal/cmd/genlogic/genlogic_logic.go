// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genlogic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func generateLogic(ctx context.Context, in CGenLogicInternalInput) {
	var (
		dirPathLogic         = gfile.Join(in.Path, in.LogicPath)
		dirPathLogicInternal = gfile.Join(dirPathLogic, "internal")
	)
	in.genItems.AppendDirPath(dirPathLogic)
	for i := 0; i < len(in.TableNames); i++ {
		generateLogicSingle(ctx, generateLogicSingleInput{
			CGenLogicInternalInput: in,
			TableName:              in.TableNames[i],
			NewTableName:           in.NewTableNames[i],
			DirPathLogic:           dirPathLogic,
			DirPathLogicInternal:   dirPathLogicInternal,
		})
	}
}

type generateLogicSingleInput struct {
	CGenLogicInternalInput
	TableName            string // TableName specifies the table name of the table.
	NewTableName         string // NewTableName specifies the prefix-stripped name of the table.
	DirPathLogic         string
	DirPathLogicInternal string
}

// generateLogicSingle generates the logic and model content of given table.
func generateLogicSingle(ctx context.Context, in generateLogicSingleInput) {
	// Generating table data preparing.
	fieldMap, err := in.DB.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf(`fetching tables fields failed for table "%s": %+v`, in.TableName, err)
	}
	var (
		tableNameCamelCase      = formatFieldName(in.NewTableName, FieldNameCaseCamel)
		tableNameCamelLowerCase = formatFieldName(in.NewTableName, FieldNameCaseCamelLower)
		tableNameSnakeCase      = gstr.CaseSnake(in.NewTableName)
		importPrefix            = in.ImportPrefix
	)
	if importPrefix == "" {
		importPrefix = utils.GetImportPath(gfile.Join(in.Path, in.LogicPath))
	} else {
		importPrefix = gstr.Join(g.SliceStr{importPrefix, in.LogicPath}, "/")
	}

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}

	// // logic - index
	// generateLogicIndex(generateLogicIndexInput{
	// 	generateLogicSingleInput: in,
	// 	TableNameCamelCase:       tableNameCamelCase,
	// 	TableNameCamelLowerCase:  tableNameCamelLowerCase,
	// 	ImportPrefix:             importPrefix,
	// 	FileName:                 fileName,
	// })

	// logic - internal
	generateLogicInternal(generateLogicInternalInput{
		generateLogicSingleInput: in,
		TableNameCamelCase:       tableNameCamelCase,
		TableNameCamelLowerCase:  tableNameCamelLowerCase,
		ImportPrefix:             importPrefix,
		FileName:                 fileName,
		FieldMap:                 fieldMap,
	})
}

type generateLogicIndexInput struct {
	generateLogicSingleInput
	TableNameCamelCase      string
	TableNameCamelLowerCase string
	ImportPrefix            string
	FileName                string
}

func generateLogicIndex(in generateLogicIndexInput) {
	path := filepath.FromSlash(gfile.Join(in.DirPathLogic, in.FileName+".go"))
	// It should add path to result slice whenever it would generate the path file or not.
	in.genItems.AppendGeneratedFilePath(path)
	if in.OverwriteLogic || !gfile.Exists(path) {
		indexContent := gstr.ReplaceByMap(
			getTemplateFromPathOrDefault(in.TplLogicIndexPath, consts.TemplateGenLogicIndexContent),
			g.MapStrStr{
				tplVarImportPrefix:            in.ImportPrefix,
				tplVarTableName:               in.TableName,
				tplVarTableNameCamelCase:      in.TableNameCamelCase,
				tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
				tplVarPackageName:             filepath.Base(in.LogicPath),
			})
		indexContent = replaceDefaultVar(in.CGenLogicInternalInput, indexContent)
		if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", gfile.RealPath(path))
		}
	}
}

type generateLogicInternalInput struct {
	generateLogicSingleInput
	TableNameCamelCase      string
	TableNameCamelLowerCase string
	ImportPrefix            string
	FileName                string
	FieldMap                map[string]*gdb.TableField
	BasePath                string
}

func generateLogicInternal(in generateLogicInternalInput) {
	in.BasePath, _ = os.Getwd()
	tmplist := strings.Split(in.BasePath, "/")
	in.BasePath = tmplist[len(tmplist)-1]

	in.DirPathLogicInternal = in.DirPathLogic + "/" + in.FileName
	path := filepath.FromSlash(gfile.Join(in.DirPathLogicInternal, in.FileName+".go"))
	removeFieldPrefixArray := gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
	modelContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplLogicInternalPath, consts.TemplateGenLogicInternalContent),
		g.MapStrStr{
			tplVarImportPrefix:            in.ImportPrefix,
			tplVarTableName:               in.TableName,
			tplVarBasePath:                in.BasePath,
			tplVarGroupName:               in.Group,
			tplVarTableNameCamelCase:      in.TableNameCamelCase,
			tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
			tplVarColumnDefine:            gstr.Trim(generateColumnDefinitionForLogic(in.FieldMap, removeFieldPrefixArray)),
			tplVarColumnNames:             gstr.Trim(generateColumnNamesForLogic(in.FieldMap, removeFieldPrefixArray)),
		})
	modelContent = replaceDefaultVar(in.CGenLogicInternalInput, modelContent)
	in.genItems.AppendGeneratedFilePath(path)
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", gfile.RealPath(path))
	}
}

// generateColumnNamesForLogic generates and returns the column names assignment content of column struct
// for specified table.
func generateColumnNamesForLogic(fieldMap map[string]*gdb.TableField, removeFieldPrefixArray []string) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForLogic(fieldMap)
	)

	for index, name := range names {
		field := fieldMap[name]

		newFiledName := field.Name
		for _, v := range removeFieldPrefixArray {
			newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
		}

		array[index] = []string{
			"            #" + formatFieldName(newFiledName, FieldNameCaseCamel) + ":",
			fmt.Sprintf(` #"%s",`, field.Name),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	namesContent := buffer.String()
	// Let's do this hack of table writer for indent!
	namesContent = gstr.Replace(namesContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(namesContent)
	return buffer.String()
}

// generateColumnDefinitionForLogic generates and returns the column names definition for specified table.
func generateColumnDefinitionForLogic(fieldMap map[string]*gdb.TableField, removeFieldPrefixArray []string) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForLogic(fieldMap)
	)

	for index, name := range names {
		var (
			field   = fieldMap[name]
			comment = gstr.Trim(gstr.ReplaceByArray(field.Comment, g.SliceStr{
				"\n", " ",
				"\r", " ",
			}))
		)
		newFiledName := field.Name
		for _, v := range removeFieldPrefixArray {
			newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
		}
		array[index] = []string{
			"    #" + formatFieldName(newFiledName, FieldNameCaseCamel),
			" # " + "string",
			" #" + fmt.Sprintf(`// %s`, comment),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	defineContent := buffer.String()
	// Let's do this hack of table writer for indent!
	defineContent = gstr.Replace(defineContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(defineContent)
	return buffer.String()
}
