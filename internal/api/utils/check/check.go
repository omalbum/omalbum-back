package check

import "log"

func Check(check func() error) {
	err := check()
	if err != nil {
		log.Print(err)
	}
}
