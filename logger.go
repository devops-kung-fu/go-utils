package util

import "log"

func isError(err error, pre string) error {
	if err != nil {
		log.Printf("%v: %v", pre, err)
	}
	return err
}

func isErrorBool(err error, pre string) (b bool) {
	if err != nil {
		log.Printf("%v: %v", pre, err)
		b = true
	}
	return
}
