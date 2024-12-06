// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genlogic

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenLogicConfig = `gfcli.gen.logic`
	CGenLogicUsage  = `gf gen logic [OPTION]`
	CGenLogicBrief  = `automatically generate go files for logic/do/entity`
	CGenLogicEg     = `
gf gen logic
gf gen logic -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
gf gen logic -p ./model -g user-center -t user,user_detail,user_login
gf gen logic -r user_
`

	CGenLogicAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing.
    The configuration node name is "gfcli.gen.logic", which also supports multiple databases, for example(config.yaml):
	gfcli:
	  gen:
		logic:
		- link:     "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
		  tables:   "order,products"
		  jsonCase: "CamelLower"
		- link:   "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
		  path:   "./my-app"
		  prefix: "primary_"
		  tables: "user, userDetail"
		  typeMapping:
			decimal:
			  type:   decimal.Decimal
			  import: github.com/shopspring/decimal
			numeric:
			  type: string
		  fieldMapping:
			table_name.field_name:
			  type:   decimal.Decimal
			  import: github.com/shopspring/decimal
`
	CGenLogicBriefPath              = `directory path for generated files`
	CGenLogicBriefLink              = `database configuration, the same as the ORM configuration of GoFrame`
	CGenLogicBriefTables            = `generate models only for given tables, multiple table names separated with ','`
	CGenLogicBriefTablesEx          = `generate models excluding given tables, multiple table names separated with ','`
	CGenLogicBriefPrefix            = `add prefix for all table of specified link/database tables`
	CGenLogicBriefRemovePrefix      = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenLogicBriefRemoveFieldPrefix = `remove specified prefix of the field, multiple prefix separated with ','`
	CGenLogicBriefStdTime           = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	CGenLogicBriefWithTime          = `add created time for auto produced go files`
	CGenLogicBriefGJsonSupport      = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	CGenLogicBriefImportPrefix      = `custom import prefix for generated go files`
	CGenLogicBriefLogicPath         = `directory path for storing generated logic files under path`
	// CGenLogicBriefDoPath            = `directory path for storing generated do files under path`
	// CGenLogicBriefEntityPath        = `directory path for storing generated entity files under path`
	CGenLogicBriefOverwriteLogic    = `overwrite all logic files both inside/outside internal folder`
	CGenLogicBriefModelFile         = `custom file name for storing generated model content`
	CGenLogicBriefModelFileForLogic = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	CGenLogicBriefDescriptionTag    = `add comment to description tag for each field`
	CGenLogicBriefNoJsonTag         = `no json tag will be added for each field`
	CGenLogicBriefNoModelComment    = `no model comment will be added for each field`
	CGenLogicBriefClear             = `delete all generated go files that do not exist in database`
	CGenLogicBriefTypeMapping       = `custom local type mapping for generated struct attributes relevant to fields of table`
	CGenLogicBriefFieldMapping      = `custom local type mapping for generated struct attributes relevant to specific fields of table`
	CGenLogicBriefGroup             = `
specifying the configuration group name of database for generated ORM instance,
it's not necessary and the default value is "default"
`
	CGenLogicBriefJsonCase = `
generated json tag case for model struct, cases are as follows:
| Case            | Example            |
|---------------- |--------------------|
| Camel           | AnyKindOfString    |
| CamelLower      | anyKindOfString    | default
| Snake           | any_kind_of_string |
| SnakeScreaming  | ANY_KIND_OF_STRING |
| SnakeFirstUpper | rgb_code_md5       |
| Kebab           | any-kind-of-string |
| KebabScreaming  | ANY-KIND-OF-STRING |
`
	CGenLogicBriefTplLogicIndexPath    = `template file path for logic index file`
	CGenLogicBriefTplLogicInternalPath = `template file path for logic internal file`
	CGenLogicBriefTplLogicDoPathPath   = `template file path for logic do file`
	CGenLogicBriefTplLogicEntityPath   = `template file path for logic entity file`

	tplVarTableName               = `{TplTableName}`
	tplVarBasePath                = `{TplBasePath}`
	tplVarTableNameCamelCase      = `{TplTableNameCamelCase}`
	tplVarTableNameCamelLowerCase = `{TplTableNameCamelLowerCase}`
	tplVarPackageImports          = `{TplPackageImports}`
	tplVarImportPrefix            = `{TplImportPrefix}`
	tplVarStructDefine            = `{TplStructDefine}`
	tplVarColumnDefine            = `{TplColumnDefine}`
	tplVarColumnNames             = `{TplColumnNames}`
	tplVarGroupName               = `{TplGroupName}`
	tplVarDatetimeStr             = `{TplDatetimeStr}`
	tplVarCreatedAtDatetimeStr    = `{TplCreatedAtDatetimeStr}`
	tplVarPackageName             = `{TplPackageName}`
)

var (
	createdAt          = gtime.Now()
	defaultTypeMapping = map[DBFieldTypeName]CustomAttributeType{
		"decimal": {
			Type: "float64",
		},
		"money": {
			Type: "float64",
		},
		"numeric": {
			Type: "float64",
		},
		"smallmoney": {
			Type: "float64",
		},
	}
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenLogicConfig`:                 CGenLogicConfig,
		`CGenLogicUsage`:                  CGenLogicUsage,
		`CGenLogicBrief`:                  CGenLogicBrief,
		`CGenLogicEg`:                     CGenLogicEg,
		`CGenLogicAd`:                     CGenLogicAd,
		`CGenLogicBriefPath`:              CGenLogicBriefPath,
		`CGenLogicBriefLink`:              CGenLogicBriefLink,
		`CGenLogicBriefTables`:            CGenLogicBriefTables,
		`CGenLogicBriefTablesEx`:          CGenLogicBriefTablesEx,
		`CGenLogicBriefPrefix`:            CGenLogicBriefPrefix,
		`CGenLogicBriefRemovePrefix`:      CGenLogicBriefRemovePrefix,
		`CGenLogicBriefRemoveFieldPrefix`: CGenLogicBriefRemoveFieldPrefix,
		`CGenLogicBriefStdTime`:           CGenLogicBriefStdTime,
		`CGenLogicBriefWithTime`:          CGenLogicBriefWithTime,
		`CGenLogicBriefLogicPath`:         CGenLogicBriefLogicPath,
		// `CGenLogicBriefDoPath`:               CGenLogicBriefDoPath,
		// `CGenLogicBriefEntityPath`:           CGenLogicBriefEntityPath,
		`CGenLogicBriefGJsonSupport`:         CGenLogicBriefGJsonSupport,
		`CGenLogicBriefImportPrefix`:         CGenLogicBriefImportPrefix,
		`CGenLogicBriefOverwriteLogic`:       CGenLogicBriefOverwriteLogic,
		`CGenLogicBriefModelFile`:            CGenLogicBriefModelFile,
		`CGenLogicBriefModelFileForLogic`:    CGenLogicBriefModelFileForLogic,
		`CGenLogicBriefDescriptionTag`:       CGenLogicBriefDescriptionTag,
		`CGenLogicBriefNoJsonTag`:            CGenLogicBriefNoJsonTag,
		`CGenLogicBriefNoModelComment`:       CGenLogicBriefNoModelComment,
		`CGenLogicBriefClear`:                CGenLogicBriefClear,
		`CGenLogicBriefTypeMapping`:          CGenLogicBriefTypeMapping,
		`CGenLogicBriefFieldMapping`:         CGenLogicBriefFieldMapping,
		`CGenLogicBriefGroup`:                CGenLogicBriefGroup,
		`CGenLogicBriefJsonCase`:             CGenLogicBriefJsonCase,
		`CGenLogicBriefTplLogicIndexPath`:    CGenLogicBriefTplLogicIndexPath,
		`CGenLogicBriefTplLogicInternalPath`: CGenLogicBriefTplLogicInternalPath,
		`CGenLogicBriefTplLogicDoPathPath`:   CGenLogicBriefTplLogicDoPathPath,
		`CGenLogicBriefTplLogicEntityPath`:   CGenLogicBriefTplLogicEntityPath,
	})
}

type (
	CGenLogic      struct{}
	CGenLogicInput struct {
		g.Meta            `name:"logic" config:"{CGenLogicConfig}" usage:"{CGenLogicUsage}" brief:"{CGenLogicBrief}" eg:"{CGenLogicEg}" ad:"{CGenLogicAd}"`
		Path              string `name:"path"                short:"p"  brief:"{CGenLogicBriefPath}" d:"internal"`
		Link              string `name:"link"                short:"l"  brief:"{CGenLogicBriefLink}"`
		Tables            string `name:"tables"              short:"t"  brief:"{CGenLogicBriefTables}"`
		TablesEx          string `name:"tablesEx"            short:"x"  brief:"{CGenLogicBriefTablesEx}"`
		Group             string `name:"group"               short:"g"  brief:"{CGenLogicBriefGroup}" d:"default"`
		Prefix            string `name:"prefix"              short:"f"  brief:"{CGenLogicBriefPrefix}"`
		RemovePrefix      string `name:"removePrefix"        short:"r"  brief:"{CGenLogicBriefRemovePrefix}"`
		RemoveFieldPrefix string `name:"removeFieldPrefix"   short:"rf" brief:"{CGenLogicBriefRemoveFieldPrefix}"`
		JsonCase          string `name:"jsonCase"            short:"j"  brief:"{CGenLogicBriefJsonCase}" d:"CamelLower"`
		ImportPrefix      string `name:"importPrefix"        short:"i"  brief:"{CGenLogicBriefImportPrefix}"`
		LogicPath         string `name:"logicPath"             short:"d"  brief:"{CGenLogicBriefLogicPath}" d:"logic"`
		// DoPath               string `name:"doPath"              short:"o"  brief:"{CGenLogicBriefDoPath}" d:"model/do"`
		// EntityPath           string `name:"entityPath"          short:"e"  brief:"{CGenLogicBriefEntityPath}" d:"model/entity"`
		TplLogicIndexPath    string `name:"tplLogicIndexPath"     short:"t1" brief:"{CGenLogicBriefTplLogicIndexPath}"`
		TplLogicInternalPath string `name:"tplLogicInternalPath"  short:"t2" brief:"{CGenLogicBriefTplLogicInternalPath}"`
		// TplLogicDoPath       string `name:"tplLogicDoPath"        short:"t3" brief:"{CGenLogicBriefTplLogicDoPathPath}"`
		// TplLogicEntityPath   string `name:"tplLogicEntityPath"    short:"t4" brief:"{CGenLogicBriefTplLogicEntityPath}"`
		StdTime        bool `name:"stdTime"             short:"s"  brief:"{CGenLogicBriefStdTime}" orphan:"true"`
		WithTime       bool `name:"withTime"            short:"w"  brief:"{CGenLogicBriefWithTime}" orphan:"true"`
		GJsonSupport   bool `name:"gJsonSupport"        short:"n"  brief:"{CGenLogicBriefGJsonSupport}" orphan:"true"`
		OverwriteLogic bool `name:"overwriteLogic"        short:"v"  brief:"{CGenLogicBriefOverwriteLogic}" orphan:"true"`
		DescriptionTag bool `name:"descriptionTag"      short:"c"  brief:"{CGenLogicBriefDescriptionTag}" orphan:"true"`
		NoJsonTag      bool `name:"noJsonTag"           short:"k"  brief:"{CGenLogicBriefNoJsonTag}" orphan:"true"`
		NoModelComment bool `name:"noModelComment"      short:"m"  brief:"{CGenLogicBriefNoModelComment}" orphan:"true"`
		Clear          bool `name:"clear"               short:"a"  brief:"{CGenLogicBriefClear}" orphan:"true"`

		TypeMapping  map[DBFieldTypeName]CustomAttributeType  `name:"typeMapping"  short:"y"  brief:"{CGenLogicBriefTypeMapping}"  orphan:"true"`
		FieldMapping map[DBTableFieldName]CustomAttributeType `name:"fieldMapping" short:"fm" brief:"{CGenLogicBriefFieldMapping}" orphan:"true"`

		// internal usage purpose.
		genItems *CGenLogicInternalGenItems
	}
	CGenLogicOutput struct{}

	CGenLogicInternalInput struct {
		CGenLogicInput
		DB            gdb.DB
		TableNames    []string
		NewTableNames []string
	}
	DBTableFieldName    = string
	DBFieldTypeName     = string
	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"`
		Import string `brief:"custom import for this type"`
	}
)

func (c CGenLogic) Logic(ctx context.Context, in CGenLogicInput) (out *CGenLogicOutput, err error) {
	in.genItems = newCGenLogicInternalGenItems()
	if in.Link != "" {
		doGenLogicForArray(ctx, -1, in)
	} else if g.Cfg().Available(ctx) {
		v := g.Cfg().MustGet(ctx, CGenLogicConfig)
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				doGenLogicForArray(ctx, i, in)
			}
		} else {
			doGenLogicForArray(ctx, -1, in)
		}
	} else {
		doGenLogicForArray(ctx, -1, in)
	}
	doClear(in.genItems)
	mlog.Print("done!")
	return
}

// doGenLogicForArray implements the "gen logic" command for configuration array.
func doGenLogicForArray(ctx context.Context, index int, in CGenLogicInput) {
	var (
		err error
		db  gdb.DB
	)
	if index >= 0 {
		err = g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.%d`, CGenLogicConfig, index),
		).Scan(&in)
		if err != nil {
			mlog.Fatalf(`invalid configuration of "%s": %+v`, CGenLogicConfig, err)
		}
	}
	if dirRealPath := gfile.RealPath(in.Path); dirRealPath == "" {
		mlog.Fatalf(`path "%s" does not exist`, in.Path)
	}
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")

	// It uses user passed database configuration.
	if in.Link != "" {
		var tempGroup = gtime.TimestampNanoStr()
		gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
			Link: in.Link,
		})
		if db, err = gdb.Instance(tempGroup); err != nil {
			mlog.Fatalf(`database initialization failed: %+v`, err)
		}
	} else {
		db = g.DB(in.Group)
	}
	if db == nil {
		mlog.Fatal(`database initialization failed, may be invalid database configuration`)
	}

	var tableNames []string
	if in.Tables != "" {
		tableNames = gstr.SplitAndTrim(in.Tables, ",")
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("fetching tables failed: %+v", err)
		}
	}
	// Table excluding.
	if in.TablesEx != "" {
		array := garray.NewStrArrayFrom(tableNames)
		for _, v := range gstr.SplitAndTrim(in.TablesEx, ",") {
			array.RemoveValue(v)
		}
		tableNames = array.Slice()
	}

	// merge default typeMapping to input typeMapping.
	if in.TypeMapping == nil {
		in.TypeMapping = defaultTypeMapping
	} else {
		for key, typeMapping := range defaultTypeMapping {
			if _, ok := in.TypeMapping[key]; !ok {
				in.TypeMapping[key] = typeMapping
			}
		}
	}

	// Generating logic & model go files one by one according to given table name.
	newTableNames := make([]string, len(tableNames))
	for i, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		newTableName = in.Prefix + newTableName
		newTableNames[i] = newTableName
	}

	in.genItems.Scale()

	// Logic: index and internal.
	generateLogic(ctx, CGenLogicInternalInput{
		CGenLogicInput: in,
		DB:             db,
		TableNames:     tableNames,
		NewTableNames:  newTableNames,
	})
	// // Do.
	// generateDo(ctx, CGenLogicInternalInput{
	// 	CGenLogicInput: in,
	// 	DB:             db,
	// 	TableNames:     tableNames,
	// 	NewTableNames:  newTableNames,
	// })
	// // Entity.
	// generateEntity(ctx, CGenLogicInternalInput{
	// 	CGenLogicInput: in,
	// 	DB:             db,
	// 	TableNames:     tableNames,
	// 	NewTableNames:  newTableNames,
	// })

	in.genItems.SetClear(in.Clear)
}

func getImportPartContent(ctx context.Context, source string, isDo bool, appendImports []string) string {
	var packageImportsArray = garray.NewStrArray()
	if isDo {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/frame/g"`)
	}

	// Time package recognition.
	if strings.Contains(source, "gtime.Time") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/os/gtime"`)
	} else if strings.Contains(source, "time.Time") {
		packageImportsArray.Append(`"time"`)
	}

	// Json type.
	if strings.Contains(source, "gjson.Json") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/encoding/gjson"`)
	}

	// Check and update imports in go.mod
	if len(appendImports) > 0 {
		goModPath := utils.GetModPath()
		if goModPath == "" {
			mlog.Fatal("go.mod not found in current project")
		}
		mod, err := modfile.Parse(goModPath, gfile.GetBytes(goModPath), nil)
		if err != nil {
			mlog.Fatalf("parse go.mod failed: %+v", err)
		}
		for _, appendImport := range appendImports {
			found := false
			for _, require := range mod.Require {
				if gstr.Contains(appendImport, require.Mod.Path) {
					found = true
					break
				}
			}
			if !found {
				if err = gproc.ShellRun(ctx, `go get `+appendImport); err != nil {
					mlog.Fatalf(`%+v`, err)
				}
			}
			packageImportsArray.Append(fmt.Sprintf(`"%s"`, appendImport))
		}
	}

	// Generate and write content to golang file.
	packageImportsStr := ""
	if packageImportsArray.Len() > 0 {
		packageImportsStr = fmt.Sprintf("import(\n%s\n)", packageImportsArray.Join("\n"))
	}
	return packageImportsStr
}

func replaceDefaultVar(in CGenLogicInternalInput, origin string) string {
	var tplCreatedAtDatetimeStr string
	var tplDatetimeStr string = createdAt.String()
	if in.WithTime {
		tplCreatedAtDatetimeStr = fmt.Sprintf(`Created at %s`, tplDatetimeStr)
	}
	return gstr.ReplaceByMap(origin, g.MapStrStr{
		tplVarDatetimeStr:          tplDatetimeStr,
		tplVarCreatedAtDatetimeStr: tplCreatedAtDatetimeStr,
	})
}

func sortFieldKeyForLogic(fieldMap map[string]*gdb.TableField) []string {
	names := make(map[int]string)
	for _, field := range fieldMap {
		names[field.Index] = field.Name
	}
	var (
		i      = 0
		j      = 0
		result = make([]string, len(names))
	)
	for {
		if len(names) == 0 {
			break
		}
		if val, ok := names[i]; ok {
			result[j] = val
			j++
			delete(names, i)
		}
		i++
	}
	return result
}

func getTemplateFromPathOrDefault(filePath string, def string) string {
	if filePath != "" {
		if contents := gfile.GetContents(filePath); contents != "" {
			return contents
		}
	}
	return def
}
