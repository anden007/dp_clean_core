<%!
import (
    "strings"
    "github.com/99designs/gqlgen/codegen"
)

%>
<%: func GQLTestFunc(objects codegen.Objects, buffer *bytes.Buffer) %>
package test

import (
	"github.com/anden007/dp_clean_core/graph/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

<% if objects != nil {
	for _, object := range objects {
		for _, field := range object.Fields {
			if field.IsResolver {
				if object.Name == "Query" {%>
func Query_<%==s field.GoFieldName%>(t *testing.T) {
	_, err := mock.MockResolver().Query().<%==s field.GoFieldName%>(mock.MockContext, "请自行修改参数")
	assert.Nil(t, err)
}
<%				} else if object.Name == "Mutation"{ 
%>
func Mutation_<%==s field.GoFieldName%>(t *testing.T) {
	_, err := mock.MockResolver().Mutation().<%==s field.GoFieldName%>(mock.MockContext, "请自行修改参数")
	assert.Nil(t, err)
}
<%
				}
			}
		}
	}
}%>