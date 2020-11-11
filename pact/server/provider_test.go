package main

import (
	"context"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"log"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

const pactServerPortAddr = 2378

var (
	repository      TodoRepository
	serverAddr      = fmt.Sprintf("localhost:%d", pactServerPortAddr)
	providerBaseURL = fmt.Sprintf("http://localhost:%d", pactServerPortAddr)
	pactUrls        = []string{
		filepath.FromSlash("pacts/jvm_todo_client-todo_api.json"),
	}
)

func startServer(repository TodoRepository, done chan bool) {
	mux := NewTodoAPI(repository)
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
		Host: "localhost",
		Provider: "TodoApi",
	}

	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: providerBaseURL,
		PactURLs:        pactUrls,
		StateHandlers: types.StateHandlers{
			// Setup any state required by the test
			// in this case, we ensure there is a "user" in the system
			"An empty repository": func() error {
				return nil
			},
		},
	})
}
