package httpin

type UnsupportedType string

func (e UnsupportedType) Error() string {
	return "unsupported type: " + string(e)
}
