package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"

	"../pool"
)

const (
	poolCapacity = 3
	queueSize    = 8000
)

var poolOfWhat = &http.Client{}

func makeGetRequest(client interface{}) {
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
func getFromPoolAndMakeRequest(initialisedPool *pool.Pool) {
	clientFromPool, err := initialisedPool.Get()
	defer initialisedPool.Put(clientFromPool)
	if err != nil {
		log.Fatal(err)
	}
	makeGetRequest(clientFromPool)
}

func TestPool(t *testing.T) {
	newPool, err := NewPool(poolCapacity, poolOfWhat)
	if err != nil {
		t.Fatal(err)
	}
	var waitGroupTesting sync.WaitGroup
	for i := 0; i < queueSize; i++ {
		waitGroupTesting.Add(1)
		go func() {
			defer waitGroupTesting.Done()
			getFromPoolAndMakeRequest(newPool)
		}()
	}
	waitGroupTesting.Wait()
	t.Log("test passed successfully")
}

func BenchmarkPool(b *testing.B) {
	newPool, err := NewPool(poolCapacity, poolOfWhat)
	if err != nil {
		log.Fatal(err)
	}

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
