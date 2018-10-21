package main

import "context"

func dataGatherer(ctx context.Context, cfg Config, errors chan error) {
	<-ctx.Done()
	errors <- nil
}
