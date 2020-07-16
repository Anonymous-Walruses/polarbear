package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/Anonymous-Walruses/polarbear/fb"
	"github.com/Anonymous-Walruses/polarbear/gh"
)

func main() {

	// Grab the tokens.
	firebaseToken := os.Getenv("FIREBASE_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")
	if firebaseToken == "" {
		log.Fatalln("could not find valid firebase token")
	}
	if githubToken == "" {
		log.Fatalln("could not find valid github token")
	}

	// Make a context.
	rootContext := context.Background()

	// Make the clients.
	firebaseClient := fb.NewClient(firebaseToken)
	githubClient := gh.NewClient(githubToken, rootContext)

	// Grab the users.
	users, err := firebaseClient.Users(rootContext)
	if err != nil {
		log.Fatalln(err)
	}

	// Then update all the users.
	wg := &sync.WaitGroup{}
	wg.Add(len(users))
	for _, user := range users {
		go func(user string, wg *sync.WaitGroup) {
			defer wg.Done()
			ok, err := gh.IsMLHFellow(user, githubClient, rootContext)
			if err != nil {
				log.Printf("failed to verify status for %s: %s\n", user, err)
			}
			if !ok {
				log.Printf("user %s is not a MLH Fellow, skipping\n", user)
				return
			}
		}(user, wg)

	}

	wg.Wait()
}
