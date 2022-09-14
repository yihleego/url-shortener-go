package main

import (
	"github.com/yihleego/base62"
	"github.com/yihleego/murmurhash3"
	"io"
	"log"
	"net/http"
	"sync"
)

var hashes []*murmur3.MurmurHash32
var cache *sync.Map // TODO Use redis instead

func init() {
	hashes = make([]*murmur3.MurmurHash32, 16)
	for i := 0; i < 16; i++ {
		hashes[i] = murmur3.New32WithSeed(i)
	}
	cache = &sync.Map{}
}

func main() {
	http.HandleFunc("/", dispatch)
	err := http.ListenAndServe(":18080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func dispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		shorten(w, r)
	} else if r.Method == http.MethodGet {
		redirect(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func shorten(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if b == nil || len(b) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(b)
	for i := range hashes {
		hashCode := hashes[i].HashString(url)
		// Base62 is better
		//key := make([]byte, base64.URLEncoding.EncodedLen(len(b)))
		//base64.URLEncoding.Encode(key, hashCode.AsBytes())
		key := base62.StdEncoding.Encode(hashCode.AsBytes())
		actual, loaded := cache.LoadOrStore(string(key), url)
		if !loaded || actual.(string) == url {
			w.Write(key)
			return
		}
	}
	w.WriteHeader(http.StatusConflict)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	v, ok := cache.Load(key)
	if ok {
		http.Redirect(w, r, v.(string), http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
