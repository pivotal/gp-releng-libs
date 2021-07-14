package retry

import (
	"errors"
	"time"

	"github.com/pivotal/gp-releng-libs/vlogs"

	"github.com/avast/retry-go"
)

func Retry(f func() error, attempts uint, delay time.Duration) error {
	funcWithRecover := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				vlogs.Info("panic detected, and recovered, retrying ...")
				err = errors.New("panic error")
			}
		}()
		err = f()
		return
	}

	return retry.Do(funcWithRecover,
		retry.DelayType(retry.FixedDelay),
		retry.Delay(delay),
		retry.Attempts(attempts),

		retry.OnRetry(func(n uint, err error) {
			vlogs.Info("retry %d, %s", n, err.Error())
		}))
}
