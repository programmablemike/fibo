package http

import (
	"net/http"
	"io"
	log "github.com/sirupsen/logrus"
)

type FiboHandler struct {}

func (*FiboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL:" + r.URL.String())
}

func Tmp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "version 2")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &FiboHandler{})
	mux.HandleFunc("/tmp", Tmp)

	//htt
}
