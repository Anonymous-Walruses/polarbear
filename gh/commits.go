package gh

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/go-github/github"
)

func GetCommitCountInRange(user string, start, end time.Time, client *github.Client) (n uint64, err error) {
	ctx := context.Background()

	// Grab all the repos.
	repos, _, err := client.Repositories.List(ctx, user, &github.RepositoryListOptions{Type: "all"})
	if err != nil {
		return n, err
	}

	// Define a function that checks if a commit was made by the author and falls between the date range.
	checkFunc := func(repo *github.Repository, wg *sync.WaitGroup, n *uint64) {
		defer wg.Done()

		// Grab all the commits.
		tCommits, _, err := client.Repositories.ListCommits(ctx, user, repo.GetName(), &github.CommitsListOptions{})
		if err != nil {
			return
		}

		for _, commit := range tCommits {
			// Grab the author and the date.
			// And check for valid dates.
			commitAuthor := strings.ToLower(commit.GetCommitter().GetLogin())
			if strings.ToLower(user) != commitAuthor {
				continue
			}

			commitDate := commit.Commit.Committer.Date
			if commitDate == nil {
				continue
			}
			if !inTimeSpan(start, end, *commitDate) {
				continue
			}

			// If we're here we atomically increase the counter.
			atomic.AddUint64(n, 1)
		}
	}

	// Initialize a wait group.
	wg := &sync.WaitGroup{}
	wg.Add(len(repos))

	for _, repo := range repos {
		go checkFunc(repo, wg, &n)
	}

	wg.Wait()

	return n, err
}
