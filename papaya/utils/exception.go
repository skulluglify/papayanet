package utils

type PnBaseException struct {
	Message string
}

type PnGroupException struct {
	PnBaseException
	Name string
}

type PnException struct{}
