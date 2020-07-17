package gh

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

const mlhFellowID int64 = 65834464

func IsMLHFellow(username string, cli *github.Client, ctx context.Context) (bool, error) {
	// Make a struct to wrap the resulting call.
	type orgResult struct {
		orgs []*github.Organization
		err  error
	}
	c := make(chan orgResult)

	// Make the call concurrently with the context.
	go func() {
		orgs, _, err := cli.Organizations.List(ctx, username, nil)
		c <- orgResult{orgs: orgs, err: err}
	}()

	// Define a function to determine if a result is
	// actually a fellow.
	filter := func(o orgResult) (bool, error) {

		if o.err != nil {
			return false, o.err
		}

		for _, org := range o.orgs {
			if org.GetID() == mlhFellowID {
				return true, nil
			}
		}

		return false, nil
	}

	// Then we have a result from the channel,
	// we need to distinguish between a context error and
	// an actual error.
	select {
	case <-ctx.Done():
		<-c
		return false, ctx.Err()
	case result := <-c:
		return filter(result)
	}

}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
