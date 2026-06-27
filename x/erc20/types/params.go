package types

type Params struct {
	Authority string
}

var DefaultAuthority = ""

func DefaultParams() Params {
	return Params{
		Authority: DefaultAuthority,
	}
}

func (p Params) Validate() error {
	if p.Authority == "" {
		return ErrInvalidParams.Wrap("authority cannot be empty")
	}
	return nil
}
