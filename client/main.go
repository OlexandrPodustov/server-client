package client

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"server-client/my_pool"
)

//func initFactory() func() interface{} {
//	fact := func() interface{} {
//		return &http.Client{}
//	}
//	return fact
//}

// To initialize a pool you need to pass expected capacity
// and instance with which you would like to operate inside of pool.
// Initial goal of this pool was to use it as the pool of http clients.
// To initialize such pool you can pass (3, &http.Client{})
func NewPool(capacity int, poolOfWhat interface{}) (*my_pool.Pool, error) {
	fact := func() interface{} {
		return poolOfWhat
	}
	pool, err := my_pool.NewPool(capacity, fact)
	//pool, err := my_pool.NewPool(capacity, initFactory())
	if err != nil {
		log.Fatal(err)
	}
	return pool, nil
}

func makeRequest(client interface{}, urlIn, contentType string, body io.Reader) {
	resp, err := client.(*http.Client).Post(urlIn, contentType, body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
}
func GetFromPoolAndMakeRequest(initialisedPool *my_pool.Pool, urlIn, contentType string, body io.Reader) {

	clientFromPool, err := initialisedPool.Get()
	defer initialisedPool.Put(clientFromPool)
	if err != nil {
		log.Fatal(err)
	}
	makeRequest(clientFromPool, urlIn, contentType, body)
}
