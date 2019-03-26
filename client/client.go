/*
Package client exists to protect the Rainforest Eagle from excessive calls and ensure that concurrency is easy to reason about.  This is in contrast to the local package which is a bare transformation of rest calls into go objects and does not concern itself with concurrency and protection.

Concurrency

On creation the client starts a go routine that listens to a request channel.  Then all client requests functions are wrappers around a send of a request and a wait for a reply on that channel.  This means in effect that all requests are linearized.

Rate Limit

Given the Eagle is a fairly small server, it seems prudent to ensure that it is not aggressively called against.  To that end a simple time based rate limit is enforced.  Any calls that come in faster than the time limit will be responded to with an error instead of being forwarded to teh Eagle.

Hardware Address

The Eagle will read all of the devices on the zigbee network but in most cases we only care about the smart meter.  The client attempts to find the expected smart meter and then caches the hardware address for that meter as it should not change over the lifecycle of the client.

*/
package client

import (
	"context"
	"sync"
	"time"

	"github.com/kklipsch/reagle/local"
)

var localclient Local
var load sync.Once

//Get returns the Local
func Get(ctx context.Context, api local.API, wait time.Duration) Local {
	load.Do(func() {
		initMetricsForAllTypes()
		localclient = NewDangerous(ctx, api, wait)
	})

	return localclient
}

//NewDangerous creates a new Local, if multiple Locals are created during a single session you've lost the concurrency protections
//this client provides so you probably shouldn't use this.  Instead use Get.
func NewDangerous(ctx context.Context, api local.API, wait time.Duration) Local {
	l := make(chan Request)
	mediator := newMediator(api, wait)
	go mediator.mediate(ctx, l)

	return Local(l)
}

//Local wraps the local api with the client potections
type Local chan<- Request

//Request sends the Request to the Local
func (l Local) Request(ctx context.Context, request Request) (interface{}, error) {
	select {
	case l <- request:
	case <-ctx.Done():
		requestCancelled.WithLabelValues(typeName(request.typ)).Inc()
		return nil, ctx.Err()
	}

	return awaitResult(ctx, request)
}
