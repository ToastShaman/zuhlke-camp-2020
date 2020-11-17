package main

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const pactServerPortAddr = 2378

var (
	repository      TodoRepository
	serverAddr      = fmt.Sprintf("localhost:%d", pactServerPortAddr)
	providerBaseURL = fmt.Sprintf("http://localhost:%d", pactServerPortAddr)
	pactUrls        = findPactFiles("pacts")
)

func findPactFiles(root string) []string {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	var filtered = make([]string, 0)
	for _, fileInfo := range files {
		if strings.HasSuffix(fileInfo.Name(), ".json") {
			filtered = append(filtered, filepath.Join(root, fileInfo.Name()))
		}
	}

	return filtered
}

func startServer(repository TodoRepository, done chan bool) {
	mux := NewTodoAPI(repository, Options{Cors: false})
	server := &http.Server{Addr: serverAddr, Handler: mux}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	go func() {
		<-done // block until we receive a value

		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		log.Fatal(server.Shutdown(ctx))
	}()
}

func stopServer(done chan<- bool) {
	done <- true
}

func TestMain(m *testing.M) {
	repository = NewInMemoryTodoRepository()

	done := make(chan bool, 1)

	startServer(repository, done)

	m.Run()

	stopServer(done)
}

func TestProvider(t *testing.T) {
	pact := &dsl.Pact{
		Host:     "localhost",
		Provider: "TodoApi",
	}

	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: providerBaseURL,
		PactURLs:        pactUrls,
		StateHandlers: types.StateHandlers{
			// Setup any state required by the test
			// in this case, we ensure there is a "user" in the system
			"an empty repository": func() error {
				return nil
			},

			"an existing todo with id=--MLqrG6LkLkkKc1iMLBt": func() error {
				repository.Persist(Todo{
					Id:       "-MLqrG6LkLkkKc1iMLBt",
					Revision: "-MLivp1BrS59mMbSN7Jr",
					Text:     "Don't forget the milk",
					Status:   "TODO",
					Category: "shopping",
					Tags:     []string{"groceries", "food"},
				})
				return nil
			},
		},
	})
}
