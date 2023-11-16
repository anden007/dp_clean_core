package extgen

import (
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/gobeam/stringy"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

func fieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}
	// 生成gorm标签
	if gormTag := fd.Directives.ForName("gormTag"); gormTag != nil {
		if gormTagVal := gormTag.Arguments.ForName("value"); gormTagVal != nil {
			f.Tag = fmt.Sprintf("%s gorm:%s", f.Tag, gormTagVal.Value.String())
		}
	}
	// 覆盖原json标签
	if coverTag := fd.Directives.ForName("coverTag"); coverTag != nil {
		if coverTagVal := coverTag.Arguments.ForName("value"); coverTagVal != nil {
			f.Tag = fmt.Sprintf("%s json:%s", f.Tag, coverTagVal.Value.String())
		}
	}
	// 生成caption标签
	if captionTag := fd.Directives.ForName("captionTag"); captionTag != nil {
		if captionTagVal := captionTag.Arguments.ForName("value"); captionTagVal != nil {
			f.Tag = fmt.Sprintf("%s caption:%s", f.Tag, captionTagVal.Value.String())
		}
	}
	return f, nil
}

func GenGraphQL(modelName string, overwrite bool) {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	p := modelgen.Plugin{
		FieldHook: fieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p), api.AddPlugin(NewExtGenPlugin(modelName, overwrite)))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

func GenModelCode(modelName string, catalog string, target string, overwrite bool) {
	catalog = strings.ToLower(catalog)
	target = strings.ToLower(target)
	switch target {
	case "all":
		genDelivery(modelName, catalog, overwrite)
		genRepository(modelName, catalog, overwrite)
		genUsecase(modelName, catalog, overwrite)
		genModule(modelName, catalog, overwrite)
	case "custom":
		genCustomDelivery(modelName, catalog, overwrite)
		genCustomRepository(modelName, catalog, overwrite)
		genCustomUsecase(modelName, catalog, overwrite)
		genCustomModule(modelName, catalog, overwrite)
	case "delivery":
		genModule(modelName, catalog, overwrite)
		genDelivery(modelName, catalog, overwrite)
	case "repository":
		genModule(modelName, catalog, overwrite)
		genRepository(modelName, catalog, overwrite)
	case "usecase":
		genModule(modelName, catalog, overwrite)
		genUsecase(modelName, catalog, overwrite)
	case "fsm":
		genFsm(modelName, catalog, overwrite)
	}
}

func genDelivery(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	usecase := stringy.New(fmt.Sprintf("%sUsecase", mo.CamelCase())).LcFirst()
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/json-iterator/go", "jsoniter")
	content.ImportAlias("github.com/liamylian/jsontime/v2/v2", "jsonTime")
	content.ImportAlias("github.com/kataras/iris/v12", "iris")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/docs", "docs")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/part", "part")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/pkg/base", "base")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/graph/model", "model")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/misc", "misc")

	// ViewModel
	content.Type().Id(fmt.Sprintf("VM_%s", mo.CamelCase())).Struct(
		jen.Qual("github.com/anden007/af_dp_clean_core/pkg/base", "BaseViewModel"),
		jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
	)
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
	).Id("ToViewModel").Params().Params(jen.Err().Error()).Block(
		jen.Return(jen.Nil()),
	)
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
	).Id("ToDBModel").Params().Params(jen.Err().Error()).Block(
		jen.Return(jen.Nil()),
	).Line()
	// Struct
	content.Type().Id(fmt.Sprintf("%sHandler", mo.CamelCase())).Struct(
		jen.Id("jsonEncoder").Qual("github.com/json-iterator/go", "API"),
		jen.Id("jwt").Qual("github.com/anden007/af_dp_clean_core/part", "IJwtService"),
		jen.Id(usecase).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Line()
	// New
	content.Func().Id("NewHttpHandler").Params(
		jen.Id("party").Qual("github.com/kataras/iris/v12", "Party"),
		jen.Id("jwt").Qual("github.com/anden007/af_dp_clean_core/part", "IJwtService"),
		jen.Id(usecase).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Block(
		jen.Comment("根据当前模块挂载路径,修改Swagger中API请求路径"),
		jen.Qual("github.com/anden007/af_dp_clean_core/docs", "SwaggerInfo").Dot("SwaggerTemplate").Op("=").Qual("misc", "ProcessSwaggerTemplate").Params(
			jen.Qual("github.com/anden007/af_dp_clean_core/docs", "SwaggerInfo").Dot("SwaggerTemplate"),
			jen.Lit(fmt.Sprintf("/_%s_%s_path_", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower())),
			jen.Id("party").Dot("GetRelPath").Call(),
		),
		jen.Line(),
		jen.Id("handler").Op(":= &").Id(fmt.Sprintf("%sHandler", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id("jsonEncoder"): jen.Qual("github.com/liamylian/jsontime/v2/v2", "ConfigWithCustomTimeFormat"),
				jen.Id("jwt"):         jen.Id("jwt"),
				jen.Id(usecase):       jen.Id(usecase),
			},
		),
		jen.Id("party").Dot("Post").Params(
			jen.Lit(fmt.Sprintf("/%s/add", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("Add"),
		),
		jen.Id("party").Dot("Post").Params(
			jen.Lit(fmt.Sprintf("/%s/delByIds", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("DelByIds"),
		),
		jen.Id("party").Dot("Post").Params(
			jen.Lit(fmt.Sprintf("/%s/edit", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("Edit"),
		),
		jen.Id("party").Dot("Get").Params(
			jen.Lit(fmt.Sprintf("/%s/getById/{id:string}", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("GetById"),
		),
		jen.Id("party").Dot("Get").Params(
			jen.Lit(fmt.Sprintf("/%s/getAll", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("GetAll"),
		),
		jen.Id("party").Dot("Post").Params(
			jen.Lit(fmt.Sprintf("/%s/getByCondition", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("GetByCondition"),
		),
	).Line()
	// Add
	content.Comment(fmt.Sprintf("@Summary %s Add", modelName))
	content.Comment("@Description 添加数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept json")
	content.Comment("@Produce json")
	content.Comment("@Param JsonData body string true \"body参数(payload)\"")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/add [post]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("Add").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Var().Id("err").Error(),
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Id("vModel").Op(":=").Id(fmt.Sprintf("VM_%s", mo.CamelCase())).Values(),
		jen.If(
			jen.Err().Op("=").Qual("github.com/anden007/af_dp_clean_core/misc", "ReadBody").Call(
				jen.Id("ctx"), jen.Op("&").Id("vModel"),
			),
			jen.Err().Op("==").Nil(),
		).Block(
			jen.Id("eCtx").Op(":=").Id("m").Dot("jwt").Dot("GetExecutorContext").Call(jen.Id("ctx")),
			jen.Err().Op("=").Id("m").Dot(usecase).Dot("Add").Call(jen.Id("vModel").Dot(mo.CamelCase()), jen.Id("eCtx"), jen.Lit("")),
		),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
				},
			),
		),
	).Line()
	// DelByIds
	content.Comment(fmt.Sprintf("@Summary %s DelByIds", modelName))
	content.Comment("@Description 删除数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept x-www-form-urlencoded")
	content.Comment("@Produce json")
	content.Comment("@Param ids formData string true \"要删除的数据IDs, 半角逗号分隔\"")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/delByIds [post]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("DelByIds").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Var().Id("err").Error(),
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Id("toDelIdArray").Op(":=").Index().String().Values(),
		jen.Id("ids").Op(":=").Id("ctx").Dot("FormValue").Call(jen.Lit("ids")),
		jen.Id("toDelIds").Op(":=").Qual("strings", "Split").Call(
			jen.Id("ids"), jen.Lit(","),
		),
		jen.Id("toDelIdArray").Op("=").Append(jen.Id("toDelIdArray"), jen.Id("toDelIds").Op("...")),
		jen.Id("eCtx").Op(":=").Id("m").Dot("jwt").Dot("GetExecutorContext").Call(jen.Id("ctx")),
		jen.Err().Op("=").Id("m").Dot(usecase).Dot("DelByIds").Call(jen.Id("toDelIdArray"), jen.Id("eCtx"), jen.Lit("")),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
				},
			),
		),
	).Line()
	// Edit
	content.Comment(fmt.Sprintf("@Summary %s Edit", modelName))
	content.Comment("@Description 编辑数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept json")
	content.Comment("@Produce json")
	content.Comment("@Param JsonData body string true \"body参数(payload)\"")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/edit [post]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("Edit").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Id("vModel").Op(":=").Id(fmt.Sprintf("VM_%s", mo.CamelCase())).Values(),
		jen.Err().Op(":=").Qual("github.com/anden007/af_dp_clean_core/misc", "ReadBody").Call(
			jen.Id("ctx"), jen.Op("&").Id("vModel"),
		),
		jen.If(
			jen.Err().Op("==").Nil(),
		).Block(
			jen.Id("eCtx").Op(":=").Id("m").Dot("jwt").Dot("GetExecutorContext").Call(jen.Id("ctx")),
			jen.Err().Op("=").Id("m").Dot(usecase).Dot("Edit").Call(jen.Id("vModel").Dot(mo.CamelCase()), jen.Id("eCtx"), jen.Lit("")),
		),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
				},
			),
		),
	).Line()
	// GetById
	content.Comment(fmt.Sprintf("@Summary %s GetById", modelName))
	content.Comment("@Description 根据Id查询数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept x-www-form-urlencoded")
	content.Comment("@Produce json")
	content.Comment("@Param id path string true \"数据Id\"")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/getById/{id} [get]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("GetById").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Var().Id("result").Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
		jen.Id("id").Op(":=").Id("ctx").Dot("Params").Call().Dot("Get").Call(jen.Lit("id")),
		jen.List(jen.Id("dbResult"), jen.Err()).Op(":=").Id("m").Dot(usecase).Dot("GetById").Call(jen.Id("id")),
		jen.If(
			jen.Id("err").Op("==").Nil(),
		).Block(
			jen.List(jen.Id("result"), jen.Err()).Op("=").Qual("github.com/anden007/af_dp_clean_core/misc", "DBModel2View").Index(
				jen.List(jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
					jen.Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
				),
			).Params(jen.Id("dbResult")),
		),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
					jen.Lit("result"):  jen.Id("result"),
				},
			),
		),
	).Line()
	// GetAll
	content.Comment(fmt.Sprintf("@Summary %s GetAll", modelName))
	content.Comment("@Description 查询所有数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept x-www-form-urlencoded")
	content.Comment("@Produce json")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/getAll [get]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("GetAll").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Var().Id("result").Index().Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
		jen.List(jen.Id("dbResult"), jen.Err()).Op(":=").Id("m").Dot(usecase).Dot("GetAll").Call(),
		jen.If(
			jen.Id("err").Op("==").Nil(),
		).Block(
			jen.List(jen.Id("result"), jen.Err()).Op("=").Qual("github.com/anden007/af_dp_clean_core/misc", "DBModelList2ViewList").Index(
				jen.List(jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
					jen.Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
				),
			).Params(jen.Id("dbResult")),
		),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
					jen.Lit("result"):  jen.Id("result"),
				},
			),
		),
	).Line()
	// GetByCondition
	content.Comment(fmt.Sprintf("@Summary %s GetByCondition", modelName))
	content.Comment("@Description 根据条件查询数据")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept json")
	content.Comment("@Produce json")
	content.Comment("@Param JsonData body string true \"body参数(payload)\"")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/getByCondition [post]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("GetByCondition").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Var().Id("err").Error(),
		jen.Var().Id("total").Int64().Op("=").Lit(0),
		jen.Var().Id("result").Index().Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
		jen.Var().Id("dbResult").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.If(
			jen.List(
				jen.Id("payload"),
				jen.Id("bErr"),
			).Op(":=").Id("ctx").Dot("GetBody").Params(),
			jen.Id("bErr").Op("==").Nil(),
		).Block(
			jen.Id("condition").Op(":=").Map(jen.String()).Qual("github.com/anden007/af_dp_clean_core/pkg", "SearchCondition").Values(),
			jen.If(
				jen.Id("jErr").Op(":=").Id("m").Dot("jsonEncoder").Dot("Unmarshal").Call(jen.Id("payload"), jen.Op("&").Id("condition")),
				jen.Id("jErr").Op("==").Nil(),
			).Block(
				jen.If(
					jen.List(jen.Id("dbResult"), jen.Id("total"), jen.Err()).Op("=").Id("m").Dot(usecase).Dot("GetByCondition").Call(jen.Id("condition")),
					jen.Err().Op("==").Nil(),
				).Block(
					jen.List(jen.Id("result"), jen.Err()).Op("=").Qual("github.com/anden007/af_dp_clean_core/misc", "DBModelList2ViewList").Index(
						jen.List(jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
							jen.Op("*").Id(fmt.Sprintf("VM_%s", mo.CamelCase())),
						),
					).Params(jen.Id("dbResult")),
				),
			).Else().Block(
				jen.Id("err").Op("=").Id("jErr"),
			),
		),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
					jen.Lit("result"): jen.Qual("github.com/kataras/iris/v12", "Map").Values(
						jen.Dict{
							jen.Lit("content"):       jen.Id("result"),
							jen.Lit("totalElements"): jen.Id("total"),
						}),
				},
			),
		),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "delivery_http.go", content.GoString(), overwrite)
}

func genRepository(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/anden007/af_dp_clean_core/graph/model", "model")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/part", "part")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/pkg", "pkg")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/misc", "misc")
	// Struct
	content.Type().Id(fmt.Sprintf("%sRepository", mo.CamelCase())).Struct(
		jen.Id("db").Qual("github.com/anden007/af_dp_clean_core/part", "IDataBase"),
		jen.Id("dataLogger").Qual("github.com/anden007/af_dp_clean_core/part", "IDataLogger"),
	).Line()
	// New
	content.Func().Id("NewRepository").Params(
		jen.Id("db").Qual("github.com/anden007/af_dp_clean_core/part", "IDataBase"),
		jen.Id("dataLogger").Qual("github.com/anden007/af_dp_clean_core/part", "IDataLogger"),
	).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sRepository", mo.CamelCase())).Block(
		jen.If(jen.Qual("github.com/anden007/af_dp_clean_core/part", "ENV").Op("==").Qual("github.com/anden007/af_dp_clean_core/part", "ENUM_ENV_DEV").Op("&&").Qual("github.com/spf13/viper", "GetBool").Params(jen.Lit("mysql.auto_migrate"))).Block(
			jen.Id("db").Dot("GetDB").Call().Dot("AutoMigrate").Params(jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values()),
		),
		jen.Return().Op("&").Id(fmt.Sprintf("%sRepository", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id("db"):         jen.Id("db"),
				jen.Id("dataLogger"): jen.Id("dataLogger"),
			},
		),
	).Line()
	// Add
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("Add").Params(
		jen.Id("entity").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Comment("实体对象默认值在模型定义中的BeforeCreate钩子方法中处理,不建议在此处初始化"),
		jen.Id("transId").Op("=").Qual("strings", "TrimSpace").Call(jen.Id("transId")),
		jen.If(
			jen.Id("transId").Op("!=").Lit(""),
		).Block(
			jen.If(
				jen.List(jen.Id("tx"), jen.Id("gErr")).Op(":=").Id("m").Dot("db").Dot("GetTransaction").Call(jen.Id("transId")),
				jen.Id("gErr").Op("==").Nil(),
			).Block(
				jen.Err().Op("=").Id("tx").Dot("Create").Call(jen.Op("&").Id("entity")).Dot("Error"),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteAddLog").Call(jen.List(jen.Id("entity"), jen.Id("ctx"), jen.Id("transId"))),
				),
			).Else().Block(
				jen.Err().Op("=").Id("gErr"),
			),
		).Else().Block(
			jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Create").Call(jen.Op("&").Id("entity")).Dot("Error"),
			jen.If(jen.Err().Op("==").Nil()).Block(
				jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteAddLog").Call(jen.List(jen.Id("entity"), jen.Id("ctx"), jen.Lit(""))),
			),
		),
		jen.Return(),
	).Line()
	// DelByIds
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("DelByIds").Params(
		jen.Id("ids").Index().String(),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Id("transId").Op("=").Qual("strings", "TrimSpace").Call(jen.Id("transId")),
		jen.If(
			jen.List(jen.Id("oldEntityArray"), jen.Id("oErr")).Op(":=").Id("m").Dot("GetByIds").Call(jen.Id("ids")),
			jen.Id("oErr").Op("==").Nil(),
		).Block(
			jen.If(
				jen.Id("transId").Op("!=").Lit(""),
			).Block(
				jen.If(
					jen.List(jen.Id("tx"), jen.Id("gErr")).Op(":=").Id("m").Dot("db").Dot("GetTransaction").Call(jen.Id("transId")),
					jen.Id("gErr").Op("==").Nil(),
				).Block(
					jen.Err().Op("=").Id("tx").Dot("Delete").Call(
						jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
						jen.Lit("id in (?)"),
						jen.Id("ids"),
					).Dot("Error"),
					jen.If(jen.Err().Op("==").Nil()).Block(
						jen.For(jen.List(jen.Id("_"), jen.Id("entity")).Op(":=").Range().Id("oldEntityArray")).Block(
							jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteDelLog").Call(jen.List(jen.Id("entity"), jen.Id("ctx"), jen.Id("transId"))),
						),
					),
				).Else().Block(
					jen.Err().Op("=").Id("gErr"),
				),
			).Else().Block(
				jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Delete").Call(
					jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
					jen.Lit("id in (?)"),
					jen.Id("ids"),
				).Dot("Error"),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.For(jen.List(jen.Id("_"), jen.Id("entity")).Op(":=").Range().Id("oldEntityArray")).Block(
						jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteDelLog").Call(jen.List(jen.Id("entity"), jen.Id("ctx"), jen.Lit(""))),
					),
				),
			),
		).Else().Block(
			jen.Err().Op("=").Id("oErr"),
		),
		jen.Return(),
	).Line()
	// Edit
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("Edit").Params(
		jen.Id("entity").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Id("transId").Op("=").Qual("strings", "TrimSpace").Call(jen.Id("transId")),
		jen.If(
			jen.List(jen.Id("oldEntity"), jen.Id("oErr")).Op(":=").Id("m").Dot("GetById").Call(jen.Id("entity").Dot("ID")),
			jen.Id("oErr").Op("==").Nil(),
		).Block(
			jen.If(
				jen.Id("transId").Op("!=").Lit(""),
			).Block(
				jen.If(
					jen.List(jen.Id("tx"), jen.Id("gErr")).Op(":=").Id("m").Dot("db").Dot("GetTransaction").Call(jen.Id("transId")),
					jen.Id("gErr").Op("==").Nil(),
				).Block(
					jen.Err().Op("=").Id("tx").Dot("Save").Call(jen.Op("&").Id("entity")).Dot("Error"),
					jen.If(jen.Err().Op("==").Nil()).Block(
						jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteEditLog").Call(jen.List(jen.Id("oldEntity"), jen.Id("entity"), jen.Id("ctx"), jen.Id("transId"))),
					),
				).Else().Block(
					jen.Err().Op("=").Id("gErr"),
				),
			).Else().Block(
				jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Save").Call(jen.Op("&").Id("entity")).Dot("Error"),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteEditLog").Call(jen.List(jen.Id("oldEntity"), jen.Id("entity"), jen.Id("ctx"), jen.Lit(""))),
				),
			),
		).Else().Block(
			jen.Err().Op("=").Id("oErr"),
		),
		jen.Return(),
	).Line()
	// Updates
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("Updates").Params(
		jen.Id("id").String(),
		jen.Id("fieldValues").Map(jen.String()).Interface(),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Id("transId").Op("=").Qual("strings", "TrimSpace").Call(jen.Id("transId")),
		jen.If(
			jen.List(jen.Id("oldEntity"), jen.Id("oErr")).Op(":=").Id("m").Dot("GetById").Call(jen.Id("id")),
			jen.Id("oErr").Op("==").Nil(),
		).Block(
			jen.If(
				jen.Id("transId").Op("!=").Lit(""),
			).Block(
				jen.If(
					jen.List(jen.Id("tx"), jen.Id("gErr")).Op(":=").Id("m").Dot("db").Dot("GetTransaction").Call(jen.Id("transId")),
					jen.Id("gErr").Op("==").Nil(),
				).Block(
					jen.Err().Op("=").Id("tx").Dot("Model").Call(
						jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
					).Dot("Where").Call(
						jen.Lit("id = ?"),
						jen.Id("id"),
					).Dot("Updates").Call(jen.Id("fieldValues")).Dot("Error"),
					jen.If(jen.Err().Op("==").Nil()).Block(
						jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteUpdateLog").Call(jen.List(jen.Id("oldEntity"), jen.Id("fieldValues"), jen.Id("ctx"), jen.Id("transId"))),
					),
				).Else().Block(
					jen.Err().Op("=").Id("gErr"),
				),
			).Else().Block(
				jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
					jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
				).Dot("Where").Call(
					jen.Lit("id = ?"),
					jen.Id("id"),
				).Dot("Updates").Call(jen.Id("fieldValues")).Dot("Error"),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Err().Op("=").Id("m").Dot("dataLogger").Dot("WriteUpdateLog").Call(jen.List(jen.Id("oldEntity"), jen.Id("fieldValues"), jen.Id("ctx"), jen.Lit(""))),
				),
			),
		).Else().Block(
			jen.Err().Op("=").Id("oErr"),
		),
		jen.Return(),
	).Line()
	// GetAll
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("GetAll").Params().Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
			jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
		).Dot("Find").Call(jen.Op("&").Id("result")).Dot("Error"),
		jen.Return(),
	).Line()
	// GetById
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("GetById").Params(
		jen.Id("id").String(),
	).Parens(
		jen.List(
			jen.Id("result").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
			jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
		).Dot("Where").Call(
			jen.Lit("id = ?"),
			jen.Id("id"),
		).Dot("First").Call(jen.Op("&").Id("result")).Dot("Error"),
		jen.Return(),
	).Line()
	// GetByIds
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("GetByIds").Params(
		jen.Id("ids").Index().String(),
	).Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Err().Op("=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
			jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
		).Dot("Where").Call(
			jen.Lit("id in (?)"),
			jen.Id("ids"),
		).Dot("Find").Call(jen.Op("&").Id("result")).Dot("Error"),
		jen.Return(),
	).Line()
	// GetByCondition
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("GetByCondition").Params(
		jen.Id("condition").Map(jen.String()).Qual("github.com/anden007/af_dp_clean_core/pkg", "SearchCondition"),
	).Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Id("total").Int64(),
			jen.Err().Error(),
		),
	).Block(
		jen.Id("query").Op(":=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
			jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
		),
		jen.Id("countQuery").Op(":=").Id("m").Dot("db").Dot("GetDB").Call().Dot("Model").Call(
			jen.Op("&").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values(),
		),
		jen.Err().Op("=").Qual("github.com/anden007/af_dp_clean_core/pkg", "NewQueryCondition").Call(
			jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", "ModelDBFields"),
		).Dot("GetQuery").Call(
			jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()).Values().Dot("TableName").Call(),
			jen.Id("condition"),
			jen.Id("query"),
			jen.Id("countQuery"),
		),
		jen.If(jen.Err().Op("==").Nil()).Block(
			jen.If(jen.Id("cErr").Op(":=").Id("countQuery").Dot("Count").Call(jen.Op("&").Id("total")).Dot("Error").Op(";").Id("cErr").Op("==").Nil()).Block(
				jen.Err().Op("=").Id("cErr"),
			),
			jen.If(jen.Id("fErr").Op(":=").Id("query").Dot("Find").Call(jen.Op("&").Id("result")).Dot("Error").Op(";").Id("fErr").Op("==").Nil()).Block(
				jen.Err().Op("=").Id("fErr"),
			),
		),
		jen.Return(),
	).Line()
	// Commit
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("Commit").Params(
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Err().Op("=").Id("m").Dot("db").Dot("Commit").Call(jen.Id("transId")),
		jen.Return(),
	).Line()
	// Rollback
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id("Rollback").Params(
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Err().Op("=").Id("m").Dot("db").Dot("Rollback").Call(jen.Id("transId")),
		jen.Return(),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "repository_mysql.go", content.GoString(), overwrite)
}

func genUsecase(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	repo := stringy.New(fmt.Sprintf("%sRepo", mo.CamelCase())).LcFirst()
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/anden007/af_dp_clean_core/graph/model", "model")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/pkg", "pkg")
	// Struct
	content.Type().Id(fmt.Sprintf("%sUsecase", mo.CamelCase())).Struct(
		jen.Id(repo).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sRepository", mo.CamelCase())),
	).Line()
	// New
	content.Func().Id("NewUsecase").Params(
		jen.Id(repo).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sRepository", mo.CamelCase())),
	).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sUsecase", mo.CamelCase())).Block(
		jen.Return().Op("&").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id(repo): jen.Id(repo),
			},
		),
	).Line()
	// Add
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Add").Params(
		jen.Id("entity").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("Add").Call(
			jen.Id("entity"), jen.Id("ctx"), jen.Id("transId"),
		),
	).Line()
	// DelByIds
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("DelByIds").Params(
		jen.Id("ids").Index().String(),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("DelByIds").Call(
			jen.Id("ids"), jen.Id("ctx"), jen.Id("transId"),
		),
	).Line()
	// Edit
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Edit").Params(
		jen.Id("entity").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("Edit").Call(
			jen.Id("entity"), jen.Id("ctx"), jen.Id("transId"),
		),
	).Line()
	// Updates
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Updates").Params(
		jen.Id("id").String(),
		jen.Id("fieldValues").Map(jen.String()).Interface(),
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("transId").String(),
	).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("Updates").Call(
			jen.Id("id"), jen.Id("fieldValues"), jen.Id("ctx"), jen.Id("transId"),
		),
	).Line()
	// GetAll
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("GetAll").Params().Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("GetAll").Call(),
	).Line()
	// GetById
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("GetById").Params(jen.Id("id").String()).Parens(
		jen.List(
			jen.Id("result").Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("GetById").Call(jen.Id("id")),
	).Line()
	// GetByIds
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("GetByIds").Params(jen.Id("ids").Index().String()).Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Err().Error(),
		),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("GetByIds").Call(jen.Id("ids")),
	).Line()
	// GetByCondition
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("GetByCondition").Params(
		jen.Id("condition").Map(jen.String()).Qual("github.com/anden007/af_dp_clean_core/pkg", "SearchCondition"),
	).Parens(
		jen.List(
			jen.Id("result").Index().Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase()),
			jen.Id("total").Int64(),
			jen.Err().Error(),
		),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("GetByCondition").Call(jen.Id("condition")),
	).Line()
	// Commit
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Commit").Params(jen.Id("transId").String()).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("Commit").Call(jen.Id("transId")),
	).Line()
	// Rollback
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Rollback").Params(jen.Id("transId").String()).Parens(
		jen.Err().Error(),
	).Block(
		jen.Return().Id("m").Dot(repo).Dot("Rollback").Call(jen.Id("transId")),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "usecase.go", content.GoString(), overwrite)
}

func genFsm(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	usecase := stringy.New(fmt.Sprintf("%sUsecase", mo.CamelCase())).LcFirst()
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/smallnest/gofsm", "fsm")
	// New
	content.Func().Id("NewFsm").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(usecase).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Op("*").Qual("github.com/smallnest/gofsm", "StateMachine").Block(
		jen.Id("delegate").Op(":=").Op("&").Qual("github.com/smallnest/gofsm", "DefaultDelegate").Values(
			jen.Dict{
				jen.Id("P"): jen.Op("&").Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())).Values(
					jen.Dict{
						jen.Id("ctx"):   jen.Id("ctx"),
						jen.Id(usecase): jen.Id(usecase),
					},
				),
			},
		),
		jen.Comment("Form:起始状态 To:下一个状态 Event:事件名,用于Trigger方法 Action:状态变更时需要执行的动作,在Action方法中执行"),
		jen.Id("transitions").Op(":=").Index().Qual("github.com/smallnest/gofsm", "Transition").Values(
			jen.Values(
				jen.Dict{
					jen.Id("From"):   jen.Lit(""),
					jen.Id("Event"):  jen.Lit("Agree"),
					jen.Id("To"):     jen.Lit("passed"),
					jen.Id("Action"): jen.Lit("Check"),
				},
			),
			jen.Values(
				jen.Dict{
					jen.Id("From"):   jen.Lit(""),
					jen.Id("Event"):  jen.Lit("Disagree"),
					jen.Id("To"):     jen.Lit("refused"),
					jen.Id("Action"): jen.Lit("Check"),
				}),
		),
		jen.Return(jen.Qual("github.com/smallnest/gofsm", "NewStateMachine").Call(jen.Id("delegate"), jen.Id("transitions").Op("..."))),
	).Line()
	// Struct
	content.Type().Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())).Struct(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(usecase).Qual("github.com/anden007/af_dp_clean_core/graph/model", fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Line()
	// Action
	content.Comment("当尝试执行Action时执行,可在此方法中对数据进行检查,检查失败则toState状态的OnEnter事件不会执行,而是触发执行OnActionFailure方法")
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())),
	).Id("Action").Params(jen.Id("action").String(), jen.Id("fromState").String(), jen.Id("toState").String(), jen.Id("args").Index().Interface()).Error().Block(
		jen.If(jen.Id("action").Op("==").Lit("Check")).Block(
			jen.Id("entity").Op(":=").Id("args").Index(jen.Id("0")).Op(".").Params(jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase())),
			jen.If(jen.Id("entity").Dot("Title").Op("==").Lit("")).Block(
				jen.Return(jen.Qual("errors", "New").Params(jen.Lit("必须输入标题"))),
			),
		),
		jen.Return().Nil(),
	).Line()
	// OnExit
	content.Comment("当退出状态时执行")
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())),
	).Id("OnExit").Params(jen.Id("fromState").String(), jen.Id("args").Index().Interface()).Block(
		jen.Qual("fmt", "Printf").Params(jen.Lit("OnExit: %s\n"), jen.Id("fromState")),
	).Line()
	// OnEnter
	content.Comment("当进入状态时执行，可用于数据库状态修改、短信通知等")
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())),
	).Id("OnEnter").Params(jen.Id("toState").String(), jen.Id("args").Index().Interface()).Block(
		jen.Qual("fmt", "Printf").Params(jen.Lit("OnEnter: %s\n"), jen.Id("toState")),
		jen.If(
			jen.Len(jen.Id("args")).Op(">").Id("0"),
		).Block(
			jen.Id("entity").Op(":=").Id("args").Index(jen.Id("0")).Op(".").Params(jen.Qual("github.com/anden007/af_dp_clean_core/graph/model", mo.CamelCase())),
			jen.Id("m").Dot(usecase).Dot("Updates").Call(jen.Id("entity").Dot("Id"), jen.Map(jen.String()).Interface().Values(
				jen.Dict{
					jen.Lit("state"): jen.Id("toState"),
				},
			), jen.Id("m").Dot("ctx"), jen.Lit("")),
		),
	).Line()
	// OnActionFailure
	content.Comment("当Action发生错误时执行")
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sEventProcessor", mo.CamelCase())),
	).Id("OnActionFailure").Params(jen.Id("action").String(), jen.Id("fromState").String(), jen.Id("toState").String(), jen.Id("args").Index().Interface(), jen.Err().Error()).Block(
		//fmt.Printf("OnActionFailure: %s|%s -> %s|%s\n", action, fromState, toState, err.Error())
		jen.Qual("fmt", "Printf").Params(jen.Lit("OnActionFailure: %s|%s -> %s|%s\n"), jen.Id("action"), jen.Id("fromState"), jen.Id("toState"), jen.Id("err").Dot("Error").Call()),
	).Line()

	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "fsm.go", content.GoString(), overwrite)
}

func genModule(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("go.uber.org/fx", "fx")
	// MockModule
	content.Var().Id("MockModule").Op("=").Qual("go.uber.org/fx", "Module").Call(
		jen.Lit(mo.SnakeCase().ToLower()),
		jen.Qual("go.uber.org/fx", "Provide").Call(jen.Id("NewRepository"), jen.Id("NewUsecase")),
	)

	// Module
	content.Func().Id("Module").Params(jen.Id("partyTag").String()).Qual("go.uber.org/fx", "Option").Block(
		jen.Return().Qual("go.uber.org/fx", "Module").Call(
			jen.Lit(mo.SnakeCase().ToLower()),
			jen.Qual("go.uber.org/fx", "Provide").Call(jen.Id("NewRepository"), jen.Id("NewUsecase")),
			jen.Qual("go.uber.org/fx", "Invoke").Call(
				jen.Qual("go.uber.org/fx", "Annotate").Call(
					jen.Id("NewHttpHandler"),
					jen.Qual("go.uber.org/fx", "ParamTags").Call(jen.Id("partyTag")),
				),
			),
		),
	)

	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "module.go", content.GoString(), overwrite)
}

func genCustomModule(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("go.uber.org/fx", "fx")
	// MockModule
	content.Var().Id("MockModule").Op("=").Qual("go.uber.org/fx", "Module").Call(
		jen.Lit(mo.SnakeCase().ToLower()),
		jen.Qual("go.uber.org/fx", "Provide").Call(jen.Id("NewRepository"), jen.Id("NewUsecase")),
	)

	// Module
	content.Func().Id("Module").Params(jen.Id("partyTag").String()).Qual("go.uber.org/fx", "Option").Block(
		jen.Return().Qual("go.uber.org/fx", "Module").Call(
			jen.Lit(mo.SnakeCase().ToLower()),
			jen.Qual("go.uber.org/fx", "Provide").Call(jen.Id("NewRepository"), jen.Id("NewUsecase")),
			jen.Qual("go.uber.org/fx", "Invoke").Call(
				jen.Qual("go.uber.org/fx", "Annotate").Call(
					jen.Id("NewHttpHandler"),
					jen.Qual("go.uber.org/fx", "ParamTags").Call(jen.Id("partyTag")),
				),
			),
		),
	)

	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "module.go", content.GoString(), overwrite)
}

func genCustomDelivery(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	usecase := stringy.New(fmt.Sprintf("%sUsecase", mo.CamelCase())).LcFirst()
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/json-iterator/go", "jsoniter")
	content.ImportAlias("github.com/liamylian/jsontime/v2/v2", "jsonTime")
	content.ImportAlias("github.com/kataras/iris/v12", "iris")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/docs", "docs")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/part", "part")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/pkg/base", "base")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/graph/model", "model")
	content.ImportAlias("github.com/anden007/af_dp_clean_core/misc", "misc")

	// Struct
	content.Type().Id(fmt.Sprintf("%sHandler", mo.CamelCase())).Struct(
		jen.Id("jsonEncoder").Qual("github.com/json-iterator/go", "API"),
		jen.Id("jwt").Qual("github.com/anden007/af_dp_clean_core/part", "IJwtService"),
		jen.Id(usecase).Id(fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Line()
	// New
	content.Func().Id("NewHttpHandler").Params(
		jen.Id("party").Qual("github.com/kataras/iris/v12", "Party"),
		jen.Id("jwt").Qual("github.com/anden007/af_dp_clean_core/part", "IJwtService"),
		jen.Id(usecase).Id(fmt.Sprintf("I%sUsecase", mo.CamelCase())),
	).Block(
		jen.Comment("根据当前模块挂载路径,修改Swagger中API请求路径"),
		jen.Qual("github.com/anden007/af_dp_clean_core/docs", "SwaggerInfo").Dot("SwaggerTemplate").Op("=").Qual("misc", "ProcessSwaggerTemplate").Params(
			jen.Qual("github.com/anden007/af_dp_clean_core/docs", "SwaggerInfo").Dot("SwaggerTemplate"),
			jen.Lit(fmt.Sprintf("/_%s_%s_path_", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower())),
			jen.Id("party").Dot("GetRelPath").Call(),
		),
		jen.Line(),
		jen.Id("handler").Op(":= &").Id(fmt.Sprintf("%sHandler", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id("jsonEncoder"): jen.Qual("github.com/liamylian/jsontime/v2/v2", "ConfigWithCustomTimeFormat"),
				jen.Id("jwt"):         jen.Id("jwt"),
				jen.Id(usecase):       jen.Id(usecase),
			},
		),
		jen.Id("party").Dot("Get").Params(
			jen.Lit(fmt.Sprintf("/%s/custom", mo.SnakeCase().ToLower())),
			jen.Id("handler").Dot("Custom"),
		),
	).Line()
	// Custom
	content.Comment(fmt.Sprintf("@Summary %s Custom", modelName))
	content.Comment("@Description 请修改")
	content.Comment(fmt.Sprintf("@Tags %s", ca.CamelCase()))
	content.Comment("@Accept json")
	content.Comment("@Produce json")
	content.Comment("@Success 200 {object} pkg.APIResult")
	content.Comment("@Failure 500 {object} pkg.APIErrorResult")
	content.Comment(fmt.Sprintf("@Router /_%s_%s_path_/%s/custom [get]", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sHandler", mo.CamelCase())),
	).Id("Custom").Params(
		jen.Id("ctx").Qual("github.com/kataras/iris/v12", "Context"),
	).Block(
		jen.Var().Id("err").Error(),
		jen.Id("success").Op(":=").True(),
		jen.Id("message").Op(":=").Lit(""),
		jen.Id("result").Op(":=").Lit(""),
		jen.Id("eCtx").Op(":=").Id("m").Dot("jwt").Dot("GetExecutorContext").Call(jen.Id("ctx")),
		jen.List(jen.Id("result"), jen.Err()).Op("=").Id("m").Dot(usecase).Dot("Custom").Call(jen.Id("eCtx")),
		jen.If(
			jen.Id("err").Op("!=").Nil(),
		).Block(
			jen.Id("success").Op("=").False(),
			jen.Id("message").Op("=").Id("err").Dot("Error").Call(),
			jen.Id("result").Op("=").Lit(""),
		),
		jen.Qual("github.com/anden007/af_dp_clean_core/misc", "WriteJson").Call(
			jen.Id("ctx"),
			jen.Qual("github.com/kataras/iris/v12", "Map").Values(
				jen.Dict{
					jen.Lit("success"): jen.Id("success"),
					jen.Lit("message"): jen.Id("message"),
					jen.Lit("result"):  jen.Id("result"),
				},
			),
		),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "delivery_http.go", content.GoString(), overwrite)
}

func genCustomUsecase(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	repo := stringy.New(fmt.Sprintf("%sRepo", mo.CamelCase())).LcFirst()
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/anden007/af_dp_clean_core/pkg", "pkg")
	// Interface
	content.Type().Id(fmt.Sprintf("I%sUsecase", mo.CamelCase())).Interface(
		jen.Id("Custom").Params(
			jen.Id("ctx").Qual("context", "Context"),
		).Parens(
			jen.List(
				jen.Id("result").String(),
				jen.Err().Error(),
			),
		),
	).Line()
	// Struct
	content.Type().Id(fmt.Sprintf("%sUsecase", mo.CamelCase())).Struct(
		jen.Id("logCenter").Qual("github.com/anden007/af_dp_clean_core/part", "ILogCenter"),
		jen.Id(repo).Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Line()
	// New
	content.Func().Id("NewUsecase").Params(
		jen.Id("logCenter").Qual("github.com/anden007/af_dp_clean_core/part", "ILogCenter"),
		jen.Id(repo).Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())),
	).Id(fmt.Sprintf("I%sUsecase", mo.CamelCase())).Block(
		jen.Return().Op("&").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id("logCenter"): jen.Id("logCenter"),
				jen.Id(repo):        jen.Id(repo),
			},
		),
	).Line()
	// Custom
	content.Func().Params(
		jen.Id("m").Op("*").Id(fmt.Sprintf("%sUsecase", mo.CamelCase())),
	).Id("Custom").Params(
		jen.Id("ctx").Qual("context", "Context"),
	).Parens(
		jen.List(
			jen.Id("result").String(),
			jen.Err().Error(),
		),
	).Block(
		jen.Id("result").Op("=").Lit("ok!"),
		jen.Return(),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "usecase.go", content.GoString(), overwrite)
}

func genCustomRepository(modelName string, catalog string, overwrite bool) {
	ca := stringy.New(catalog)
	mo := stringy.New(modelName)
	content := jen.NewFile(fmt.Sprintf("%s_%s", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()))
	// Import
	content.ImportAlias("github.com/anden007/af_dp_clean_core/part", "part")
	// Struct
	content.Type().Id(fmt.Sprintf("%sRepository", mo.CamelCase())).Struct(
		jen.Id("db").Qual("github.com/anden007/af_dp_clean_core/part", "IDataBase"),
		jen.Id("dataLogger").Qual("github.com/anden007/af_dp_clean_core/part", "IDataLogger"),
	).Line()
	// New
	content.Func().Id("NewRepository").Params(
		jen.Id("db").Qual("github.com/anden007/af_dp_clean_core/part", "IDataBase"),
		jen.Id("dataLogger").Qual("github.com/anden007/af_dp_clean_core/part", "IDataLogger"),
	).Op("*").Id(fmt.Sprintf("%sRepository", mo.CamelCase())).Block(
		jen.Return().Op("&").Id(fmt.Sprintf("%sRepository", mo.CamelCase())).Values(
			jen.Dict{
				jen.Id("db"):         jen.Id("db"),
				jen.Id("dataLogger"): jen.Id("dataLogger"),
			},
		),
	).Line()
	writeFile(fmt.Sprintf("./modules/%s_%s/", ca.SnakeCase().ToLower(), mo.SnakeCase().ToLower()), "repository_mysql.go", content.GoString(), overwrite)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeFile(filePath, fileName, content string, overwrite bool) {
	targetFile := filePath + fileName
	if exist, _ := pathExists(filePath); !exist {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			fmt.Println("× 生成目录失败，请检查文件路径", err.Error())
			return
		}
	}
	if existOldFile, _ := pathExists(targetFile); existOldFile && !overwrite {
		targetFile += ".new"
	}

	if fileObj, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
		defer fileObj.Close()
		if fileObj != nil {
			if _, err := fileObj.WriteString(content); err == nil {
				fmt.Printf("√ %s  --生成成功\n", targetFile)
				return
			}
		}
	} else {
		fmt.Println("× 生成文件失败，请检查文件路径", err.Error())
	}
	fmt.Printf("× %s  --生成失败\n", targetFile)
}
