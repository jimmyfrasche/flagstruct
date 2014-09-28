package flagstruct

import (
	"errors"
	"fmt"
)

func errf(spec string, v ...interface{}) error {
	return fmt.Errorf("flagstruct: "+spec, v...)
}

func errw(err error) error {
	return errors.New("flagstruct: " + err.Error())
}
