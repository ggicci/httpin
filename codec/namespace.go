package codec

import (
	"github.com/ggicci/strconvx"
)

type Namespace strconvx.Namespace

func NewNamespace() *Namespace {
	return (*Namespace)(strconvx.NewNamespace())
}
