package utils

import "time"

type ExecFunc func() error

func Retry(exec ExecFunc, delay time.Duration, count int) (err error)  {
	for i := 0; i < count; i++ {
		if err = exec(); err != nil {
			time.Sleep(delay)
		}
	}
	return
}
