package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/OlexandrPodustov/server-client/pool"
)

const (
	poolCapacity = 3
	queueSize    = 8000
)

func makeGetRequest(client interface{}) {
	resp, err := client.(*http.Client).Get("http://127.0.0.1:8082/publish")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func getFromPoolAndMakeRequest(initialisedPool *pool.Pool) {
	clientFromPool, err := initialisedPool.Get()
	if err != nil {
		log.Fatal(err)
	}

	defer initialisedPool.Put(clientFromPool)

	makeGetRequest(clientFromPool)
}

func TestPool(t *testing.T) {
	newPool := newPool(poolCapacity, &http.Client{})

	var waitGroupTesting sync.WaitGroup

	for i := 0; i < queueSize; i++ {
		waitGroupTesting.Add(1)

		go func() {
			defer waitGroupTesting.Done()
			getFromPoolAndMakeRequest(newPool)
		}()
	}
	waitGroupTesting.Wait()
}

func BenchmarkPool(b *testing.B) {
	newPool := newPool(poolCapacity, &http.Client{})

	var waitGroupBench sync.WaitGroup

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		waitGroupBench.Add(1)

		go func() {
			defer waitGroupBench.Done()
			getFromPoolAndMakeRequest(newPool)
		}()
	}
	waitGroupBench.Wait()
}

func newPool(capacity int, poolOfWhat interface{}) *pool.Pool {
	fact := func() interface{} {
		return poolOfWhat
	}

	newPool, err := pool.New(capacity, fact)
	if err != nil {
		log.Fatal(err)
	}

	return newPool
}
