package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	statfulHost = flag.String("statful.host", "api.statful.com", "Statful host that will be used to send metrics to.")
	statfulPort = flag.Int("statful.port", 443, "Statful port that will be used to send metrics to.")
	statfulApiToken = flag.String("statful.api-token", "", "Statful API Token for authentication purposes.")
	statfulProtocol = flag.String("statful.protocol", "api", "Protocol to be used to send metrics to statful. Can either be `UDP` or `API`")
	statfulTimeout = flag.Duration("statful.timeout", 2*time.Second, "Maximum time to wait for a response from the server. Defaults to 2s.")
	statfulDryRun = flag.Bool("statful.dry-run", false, "Don't send metrics to statful.")
	statfulNamespace = flag.String("statful.namespace", "", "")
	statfulBasePath = flag.String("statful.base-path", "/tel/v2.0/metrics", "Base Path of the api to be used for sending metrics to statful.")
	//statfulTags Tags

	prometheusHost = flag.String("prometheus.host", "prometheus", "Prometheus host to fetch metrics from.")
	prometheusPollingInterval = flag.Duration("prometheus.polling-interval", 10*time.Second, "How frequent to query data from prometheus.")
	prometheusRequestTimeout = flag.Duration("prometheus.request-timeout", 5*time.Second, "Maximum time for a request of data from prometheus.")
)

func main() {
	//flag.Var(&statfulTags, "statful.tags", "Tags to be appended for every metric before sending to statful.")
	flag.Parse()
	fmt.Println("Starting prometheus-statful-exporter")

	supplier, err := NewPrometheusSupplier(PrometheusSupplierConfig{
		host: *prometheusHost,
		requestTimeout: *prometheusRequestTimeout,
		pollingInterval: *prometheusPollingInterval,
	})

	if err != nil {
		log.Fatal("Failed to create prometheus supplier:", err)
	}

	consumer, err := NewLoggerConsumer()

	if err != nil {
		log.Fatal("Failed to create logger consumer:", err)
	}

	messages := make(chan Message)
	defer close(messages)

	go supplier.Supply(messages)
	consumer.Consume(messages)
}
