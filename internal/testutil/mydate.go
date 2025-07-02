package testutil

import "time"

type MyDate time.Time // adapted from time.Time

func (t MyDate) ToString() (string, error) {
	return time.Time(t).Format("2006-01-02"), nil
}

func (t *MyDate) FromString(value string) error {
	v, err := time.Parse("2006-01-02", value)
	if err != nil {
		return &InvalidDate{Value: value, Err: err}
	}
	*t = MyDate(v)
	return nil
}
