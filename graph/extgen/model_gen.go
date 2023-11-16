package extgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/anden007/dp_clean_core/graph/extgen/template"
	"github.com/anden007/dp_clean_core/misc"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/plugin"
	"github.com/gobeam/stringy"
)

func NewExtGenPlugin(objectName string, overwrite bool) plugin.Plugin {
	return &Plugin{
		objectName: objectName,
		overwrite:  overwrite,
	}
}

type Plugin struct {
	objectName string
	overwrite  bool
}

func (m *Plugin) Name() string {
	return "extgen"
}

// func (m *Plugin) MutateConfig(cfg *config.Config) error {
// 	return nil
// }

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	genObjects := codegen.Objects{}
	for _, object := range data.Objects {
		if !object.BuiltIn {
			if object.Name == "Query" {
				genObjects = append(genObjects, object)
			} else if object.Name == "Mutation" {
				genObjects = append(genObjects, object)
			} else {
				objectNameList := strings.Split(m.objectName, ";")
				for _, objectName := range objectNameList {
					if strings.EqualFold(objectName, object.Name) {
						gqlfileName := path.Base(object.Position.Src.Name)
						fileType := path.Ext(gqlfileName)
						catalog := strings.TrimSuffix(gqlfileName, fileType)
						fileName := stringy.New(object.Name)
						fileFullPath := fmt.Sprintf("./graph/model/%s_ext.go", fileName.SnakeCase().ToLower())
						if fileExists, _ := misc.PathExists(fileFullPath); fileExists && !m.overwrite {
							fileFullPath = fmt.Sprintf("./graph/model/%s_ext.go.new", fileName.SnakeCase().ToLower())
						}
						if fileObj, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
							defer fileObj.Close()
							buffer := new(bytes.Buffer)
							hasCreateTime := 0
							hasUpdateTime := 0
							for _, field := range object.Fields {
								if field.FieldDefinition.Type.NonNull {
									if strings.EqualFold(field.GoFieldName, "CreateTime") {
										hasCreateTime++
									}
									if strings.EqualFold(field.GoFieldName, "UpdateTime") {
										hasUpdateTime++
									}
								}
							}
							template.ModelExt(object, template.ModelExtOption{
								PackageName:   "model",
								HasCreateTime: hasCreateTime > 0,
								HasUpdateTime: hasUpdateTime > 0,
							}, buffer)
							buffer.WriteTo(fileObj)
						}
						GenModelCode(object.Name, catalog, "all", false)
					}
				}
			}
		}
	}
	if genObjects != nil {
		fileFullPath := "./graph/test/gql_api.go"
		if fileExists, _ := misc.PathExists(fileFullPath); fileExists && !m.overwrite {
			fileFullPath = "./graph/test/gql_api.go.new"
		}
		if fileObj, err := os.OpenFile(fileFullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
			defer fileObj.Close()
			buffer := new(bytes.Buffer)
			template.GQLTestFunc(genObjects, buffer)
			buffer.WriteTo(fileObj)
		}
	}
	return nil
}
