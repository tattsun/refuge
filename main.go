package main

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

func handleProxy(w http.ResponseWriter, r *http.Request) {
	hj, _ := w.(http.Hijacker)
	log.Println(r.Host)

	proxy, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer proxy.Close()

	client, _, err := hj.Hijack()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	auth := "test:pass"
	authb64 := base64.StdEncoding.EncodeToString([]byte(auth))
	r.Header.Set("Proxy-Authorization", "Basic "+authb64)
	if err = r.Write(proxy); err != nil {
		log.Fatal(errors.Wrap(err, "failed to write request "))
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		io.Copy(client, proxy)
		wg.Done()
	}()
	go func() {
		io.Copy(proxy, client)
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	log.Fatal(http.ListenAndServe(":86", http.HandlerFunc(handleProxy)))
}
