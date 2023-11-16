package misc

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/oleiade/reflections"
	"gorm.io/gorm/schema"
)

type Field struct {
	Type           string
	Caption        string
	ModelFieldName string
	DBFieldName    string
}

// 通过反射转换，反向通过字段Tag信息查找字段名，然后转换成实际数据库字段名称，此方法在系统启动时运行，不影响运行效率
func GetModelDBFieldNames(module interface{}) (result map[string]*Field, err error) {
	result = make(map[string]*Field)
	err = nil
	// 反射模型字段，生成数据库字段
	typ := reflect.TypeOf(module)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("misc.GetModelDBFieldNames执行错误：参数module不是Struct类型")
	}
	jsonTagFields, jErr := reflections.TagsDeep(module, "json")
	captionTagFields, cErr := reflections.TagsDeep(module, "caption")
	if jErr == nil && cErr == nil {
		for fieldName, jsonTag := range jsonTagFields {
			if jsonTag != "-" {
				if strings.Contains(jsonTag, "omitempty") {
					re1, _ := regexp.Compile(`\,(\s+)?omitempty`)
					jsonTag = re1.ReplaceAllString(jsonTag, "")
				}
				fieldTyp, _ := reflections.GetFieldType(module, fieldName)
				result[jsonTag] = &Field{
					ModelFieldName: fieldName,
					DBFieldName:    schema.NamingStrategy{}.ColumnName("", fieldName),
					Type:           fieldTyp,
					Caption:        captionTagFields[fieldName],
				}
			}
		}
	}
	return
}

func GetObjProperty(obj interface{}, fieldName string) (interface{}, error) {
	fields, err := reflections.Fields(obj)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		if strings.EqualFold(field, fieldName) {
			return reflections.GetField(obj, field)
		}
	}
	return nil, errors.New("Object doesn't have a property named: " + fieldName)
}
