package responses

import "PapayaNet/papaya/koala"

type PnRespBody struct {
	Map koala.KMap
}
type PnRespBodyImpl interface {
	Get(name string) any
}

func (param *PnRespBody) Get(name string) any {

	// not implemented
	return nil
}
