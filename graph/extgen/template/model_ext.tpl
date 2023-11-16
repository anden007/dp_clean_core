<%!
import (
    "strings"
    "github.com/99designs/gqlgen/codegen"
	"github.com/gobeam/stringy"
)

type ModelExtOption struct{
    PackageName string
	HasCreateTime bool
	HasUpdateTime bool
}
%>
<%: func ModelExt(object *codegen.Object, options ModelExtOption, buffer *bytes.Buffer) %>
package <%==s options.PackageName %>

import (
	misc "github.com/anden007/dp_clean_core/misc"
	base "github.com/anden007/dp_clean_core/pkg/base"
	"reflect"
	"strings"
<%if options.HasCreateTime || options.HasUpdateTime {%>
	"time"
<%}%>
	"gorm.io/gorm"
)

func (<%==s object.Name%>) TableName() string {
	return "t_<%==s stringy.New(object.Name).SnakeCase().ToLower()%>"
}

func (<%==s object.Name%>) ModelName() string {
	return "<%==s object.Name%>" // 请修改表名，将在数据日志中使用
}

// BeforeCreate gorm钩子，在创建对象之前执行，可用于默认值初始化（正常默认值建议用default关键字）
func (m *<%==s object.Name%>) BeforeCreate(tx *gorm.DB) (err error) {
	// 默认生成Id
	if strings.TrimSpace(m.ID) == "" {
		m.ID = misc.NewHexId()
	}
	<%if options.HasCreateTime {%>
	// 初始化CreateTime
	if m.CreateTime.IsZero() {
		m.CreateTime = time.Now()
	}<%}%>
	return
}

// BeforeUpdate gorm钩子,在修改对象之前执行,可用于默认值初始化(正常默认值建议用default关键字)
func (m *<%==s object.Name%>) BeforeUpdate(tx *gorm.DB) (err error) {
	<%if options.HasUpdateTime {%>
	// 自动变更UpdateTime字段初始值
	m.UpdateTime = time.Now()<%}%>
	return
}

type I<%==s object.Name%>Repository interface {
	base.ICRUDRepository[<%==s object.Name%>]
}

type I<%==s object.Name%>Usecase interface {
	base.ICRUDUsecase[<%==s object.Name%>]
}

func init() {
	entity := <%==s object.Name%>{}
	tableName := entity.TableName()
	DomainModels["<%==s object.Name%>"] = reflect.TypeOf((*<%==s object.Name%>)(nil)).Elem()
	if fields, err := misc.GetModelDBFieldNames(entity); err == nil {
		ModelDBFields[tableName] = make(map[string]string)
		for fieldName, field := range fields {
			ModelDBFields[tableName][fieldName] = field.DBFieldName
		}
	} else {
		panic(err)
	}
}
