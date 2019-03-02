package client

import (
	"context"
	"fmt"
	"time"
)

//ErrRateLimited is returned when the RateLimit is being enforced
var ErrRateLimited = fmt.Errorf("not enough time has passed since last action")

//RateLimit is used to implement a simple time based rate limit
type RateLimit struct {
	blocked AtomicBool
	wait    time.Duration
}

//NewRateLimit returns a RateLimit with the passed in wait
func NewRateLimit(wait time.Duration) RateLimit {
	return RateLimit{wait: wait}
}

//EnforceLimit will return ErrRateLimited if the limit is already in effect.  If not it will start the limit and return nil.
func EnforceLimit(ctx context.Context, limit RateLimit) error {
	if LoadBool(limit.blocked) {
		return ErrRateLimited
	}

	StoreBool(limit.blocked, true)
	go func() {
		select {
		case <-time.After(limit.wait):
		case <-ctx.Done():
		}

		StoreBool(limit.blocked, false)
	}()

	return nil
}

//Limited wraps the provided action with a RateLimit.  If the RateLimit is in effect the action is not called and ErrRateLimited is returned
//otherwise the action is called and its error returned.
func Limited(ctx context.Context, limit RateLimit, action func() error) func() error {
	return func() error {
		if err := EnforceLimit(ctx, limit); err != nil {
			return err
		}

		return action()
	}
}
