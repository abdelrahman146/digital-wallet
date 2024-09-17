package types

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

type Interval time.Duration

func (i Interval) Duration() time.Duration {
	return time.Duration(i)
}

func (i Interval) Value() (driver.Value, error) {
	return time.Duration(i).String(), nil
}

func (i *Interval) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}
	// Convert format of hh:mm:ss into format parseable by time.ParseDuration()
	str = strings.Replace(str, ":", "h", 1)
	str = strings.Replace(str, ":", "m", 1)
	str += "s"
	dur, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*i = Interval(dur)
	return nil
}
