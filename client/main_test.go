package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"server-client/my_pool"
	"sync"
	"testing"
)

const (
	poolCapacity int = 3
	queueSize    int = 8000
)

var poolOfWhat = &http.Client{}

func makeRequest(client interface{}) {
	resp, err := client.(*http.Client).Get("http://127.0.0.1:8082/publish")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//ioutil.ReadAll(resp.Body)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println("Resp body: ", string(responseBody))
}
func getFromPoolAndMakeRequest(initialisedPool *my_pool.Pool) {
	clientFromPool, err := initialisedPool.Get()
	defer initialisedPool.Put(clientFromPool)
	if err != nil {
		log.Fatal(err)
	}
	makeRequest(clientFromPool)
}

func TestPool(t *testing.T) {
	pool, err := NewPool(poolCapacity, poolOfWhat)
	if err != nil {
		t.Fatal(err)
	}
	var waitGroupTesting sync.WaitGroup
	for i := 0; i < queueSize; i++ {
		waitGroupTesting.Add(1)
		go func() {
			defer waitGroupTesting.Done()
			getFromPoolAndMakeRequest(pool)
		}()
	}
	waitGroupTesting.Wait()
	t.Log("test passed successfully")
}

func BenchmarkPool(b *testing.B) {
	pool, err := NewPool(poolCapacity, poolOfWhat)
	if err != nil {
		log.Fatal(err)
	}

	var waitGroupBench sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		waitGroupBench.Add(1)
		go func() {
			defer waitGroupBench.Done()
			getFromPoolAndMakeRequest(pool)
		}()
	}
	waitGroupBench.Wait()
}
