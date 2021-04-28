package httpin

type Directive struct {
	Key string // e.g. query.page, header.x-api-token
}

func BuildDirective(key string) (*Directive, error) {
	return &Directive{Key: key}, nil
}
