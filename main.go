package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"math/rand"
	"time"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"logged_users": "true",
	},
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "goapp_http_requests_total",
	Help: "Total Http requests",
}, []string{})


var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "goapp_http_request_duration",
	Help: "Duration in seconds of all HTTP requests",
}, []string{"handler"})

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)

	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(2000)))
		}
	}()

	home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(8))*time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello Full Cycle"))
	})
	
	contact := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(18))*time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Contact"))
	})

	duration := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, home),
	)

	duration_contact := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "contact"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, contact),
	)

	http.Handle("/", duration)
	http.Handle("/contact", duration_contact)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8181", nil))
}