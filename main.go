package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var bucket Bucket

type Bucket struct {
	MaxToken                int64
	Tokens                  int64
	Tag                     string
	RefreshTokenTimeMiliSec int64
}

var tokensConsumed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_conumed_tokens", // metric name
		Help: "Tokens Consumed.",
	},
	[]string{"endpoint", "bucket_tag"}, // labels
)

var deniedRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_denied_requests", // metric name
		Help: "Denied .",
	},
	[]string{"endpoint", "bucket_tag"}, // labels
)

func (b *Bucket) AddTokens() {
	for true {
		time.Sleep(time.Duration(b.RefreshTokenTimeMiliSec * time.Hour.Milliseconds()))
		b.Tokens = b.Tokens + 1
		if b.Tokens > b.MaxToken {
			b.Tokens = b.MaxToken
		}
	}
}

func (b *Bucket) TakeTokens() bool {
	b.Tokens = b.Tokens - 1
	if b.Tokens < 0 {
		deniedRequest.WithLabelValues("/", bucket.Tag).Inc()
		b.Tokens = 0
		return false
	}
	return true

}

func (b *Bucket) Init(maxToken int64, tag string, refreshTokenTimeMili int64) {
	b.MaxToken = maxToken
	b.Tokens = 0
	b.Tag = tag
	b.RefreshTokenTimeMiliSec = refreshTokenTimeMili
	go b.AddTokens()
}

func takeTokensController(w http.ResponseWriter, r *http.Request) {
	tokensConsumed.WithLabelValues("/", bucket.Tag).Inc()
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	if bucket.TakeTokens() {
		w.WriteHeader(http.StatusOK)
		resp["message"] = "Status Created"

	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		resp["message"] = "Too Many Requests"
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func handler() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", takeTokensController)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	prometheus.MustRegister(tokensConsumed)
	prometheus.MustRegister(deniedRequest)
	bucket.Init(10, "teste", 100)
	handler()

}
