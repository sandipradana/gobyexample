package main

import (
	"context"
	"net/http"
	"time"
)

const limaDetik = time.Second * time.Duration(5)

func timeout(h CustomHandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
	
		respChannel := make(chan Response)
		
		ctx, cancel := context.WithTimeout(r.Context(), limaDetik)
		
		defer cancel()
		
		r = r.WithContext(ctx)
		
		go func() {
		
			handlerResp, err := h(w, r)
			
			if err != nil {
				// do something with error
			}
			
			respChannel <- handlerResp
			
		}()

		select {
			// do something when context done happened
			case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("Request time out"))
			return
			case resp := <-respChannel:
			// do something with response
			w.Write(resp)
		}
	}
}

type Response []byte

type CustomHandlerFunc func(http.ResponseWriter, *http.Request) (Response, error)

const enamDetik = time.Second * time.Duration(6)
const sepuluhDetik = time.Second * time.Duration(10)

func customHandlerFunc1(w http.ResponseWriter, r *http.Request) (resp Response, err error) {
	/**/
	//misal proses eksekusi (di contoh ini sleep 6 detik) yang lamanya melebihi
	//batas waktu timeout (5 detik) yang sudah ditentukan
	time.Sleep(enamDetik)
	/**/
	resp, err = []byte("time out"), nil
	return
}

func customHandlerFunc2(w http.ResponseWriter, r *http.Request) (resp Response, err error) {
	resp, err = []byte("tidak timeout"), nil
	return
}

func customHandlerFunc3(w http.ResponseWriter, r *http.Request) {
	time.Sleep(sepuluhDetik)
	w.Write([]byte("Sepuluh detik"))
}

func main() {

	http.HandleFunc("/timeout", timeout(customHandlerFunc1))
	http.HandleFunc("/tidak-timeout", timeout(customHandlerFunc2))
	http.HandleFunc("/tidak-ada-timeout", customHandlerFunc3)
	
	http.ListenAndServe(":8080", nil)
	
}
