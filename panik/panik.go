package panik

import (
	"fmt"
)

func On(err error) {
	if err != nil {
		panic(err)
	}
}

func If(condition bool, message string, params ...interface{}) {
	if condition {
		s := fmt.Sprintf(message, params...)
		panic(s)
	}
}

func Do(message string, params ...interface{}) {
	If(true, message, params...)
}
