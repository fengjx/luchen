package calcpub

import "context"

var CalcAPI calcAPI

type calcAPI interface {
	Add(context.Context, int32, int32) (int32, error)
}

func SetCalcAPI(impl calcAPI) {
	CalcAPI = impl
}
