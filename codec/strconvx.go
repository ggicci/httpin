package codec

import (
	"github.com/ggicci/strconvx"
)

// strconvx is a tiny Go package that defines a unified interface for
// converting values to and from strings using `ToString` and `FromString`.

type (
	StringCodec        = strconvx.StringCodec
	StringCodecAdaptor = strconvx.AnyAdaptor
)

var (
	ErrFieldTypeMismatch    = strconvx.ErrTypeMismatch
	ErrUnsupportedFieldType = strconvx.ErrUnsupportedType
)
