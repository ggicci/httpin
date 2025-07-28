package codec

import (
	"github.com/ggicci/strconvx"
)

type Namespace struct {
	*strconvx.Namespace
}

func NewNamespace() *Namespace {
	return &Namespace{
		Namespace: strconvx.NewNamespace(),
	}
}
