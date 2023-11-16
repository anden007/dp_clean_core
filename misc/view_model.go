package misc

import (
	"github.com/anden007/af_dp_clean_core/pkg/base"

	"github.com/gotidy/copy"
	"github.com/kataras/iris/v12"
	jsonTime "github.com/liamylian/jsontime/v2/v2"
)

func ReadBody(ctx iris.Context, vModel base.IViewModel) (err error) {
	jsonEncoder := jsonTime.ConfigWithCustomTimeFormat
	if ctx != nil {
		if payload, bodyErr := ctx.GetBody(); bodyErr == nil {
			err = jsonEncoder.Unmarshal(payload, &vModel)
			if err == nil {
				err = vModel.ToDBModel()
			}
		} else {
			err = bodyErr
		}
	}
	return
}

func DBModelList2ViewList[Source base.IDBModel, DEST base.IViewModel](models []Source) (result []DEST, err error) {
	copiers := copy.New()
	copier := copiers.Get(new(DEST), new(Source))
	for i := 0; i < len(models); i++ {
		var dst DEST
		copier.Copy(&dst, &(models[i]))
		if err == nil {
			err = dst.ToViewModel()
		}
		result = append(result, dst)
	}
	return
}

func DBModel2View[Source base.IDBModel, DEST base.IViewModel](model Source) (result DEST, err error) {
	copiers := copy.New()
	copier := copiers.Get(new(DEST), new(Source))
	copier.Copy(&result, &(model))
	if err == nil {
		err = result.ToViewModel()
	}
	return
}
