package responses

import "PapayaNet/papaya/koala"

type PnRespParams struct {
	Map koala.KMap
}
type PnRespParamsImpl interface {
	Get(name string) any
}

func (param *PnRespParams) Get(name string) any {

	// not implemented
	return nil
}
