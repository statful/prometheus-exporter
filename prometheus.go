package main

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
	"time"
)

type Message struct {
	content []interface{}
}

type Supplier interface {
	Supply(channel chan<- Message)
}

func NewPrometheusSupplier(config PrometheusSupplierConfig) (Supplier, error) {
	prometheusClient, err := api.NewClient(api.Config{
		Address: config.host,
	})

	if err != nil {
		return nil, err
	}

	prometheusApiClient := v1.NewAPI(prometheusClient)

	return &prometheusSupplier{
		api:             prometheusApiClient,
		pollingInterval: config.pollingInterval,
		requestTimeout:  config.requestTimeout,
	}, nil
}

type prometheusSupplier struct {
	api             v1.API
	requestTimeout  time.Duration
	pollingInterval time.Duration
}

type PrometheusSupplierConfig struct {
	host            string
	pollingInterval time.Duration
	requestTimeout  time.Duration
}

func (ps *prometheusSupplier) Supply(mc chan<- Message) {
	startTime := time.Now()
	for {
		time.Sleep(ps.pollingInterval)

		// get all labels
		timeoutCtx, cancelContext := context.WithTimeout(context.Background(), ps.requestTimeout)
		//labelNames, _, err := ps.api.LabelNames(timeoutCtx)
		//cancelContext()
		//if err != nil {
		//	log.Println("Could not obtain label names from prometheus:", err)
		//	continue
		//}

		//log.Println("Found the following labelNames:", strings.Join(labelNames, "\n"))

		// get metric names
		timeoutCtx, cancelContext = context.WithTimeout(context.Background(), ps.requestTimeout)
		metricNames, _, err := ps.api.LabelValues(timeoutCtx, "__name__")
		cancelContext()
		if err != nil {
			log.Println("Could not obtain metric names from prometheus:", err)
			continue
		}
		//log.Println("Found the following metric names:", metricNames)

		// get metric values
		endTime := time.Now()
		for _, metricName := range metricNames {
			//go func() {
			//	log.Println("Requesting values for ", string(metricName))
				timeoutCtx, cancelContext = context.WithTimeout(context.Background(), ps.requestTimeout)
				metricValues, _, err := ps.api.QueryRange(timeoutCtx, string(metricName), v1.Range{
					Start: startTime,
					End:   endTime,
					Step:  10 * time.Second,
				})
				cancelContext()

				if err != nil {
					log.Println("Could not obtain metric values from prometheus:", err)
					return
				}

				metricValues.Type()
				mc <- Message{content: []interface{}{metricName, ":", metricValues.Type()},}
			//}()
		}

		startTime = endTime
	}
}

type Consumer interface {
	Consume(<-chan Message)
}

func NewLoggerConsumer() (Consumer, error) {
	return &loggerConsumer{}, nil
}

type loggerConsumer struct {
}

func (lc *loggerConsumer) Consume(mc <-chan Message) {
	for {
		select {
		case m, ok := <-mc:
			if !ok {
				log.Println("Events channel closed.")
				return
			}
			log.Println(m.content...)
		}
	}
}

type Mapper interface {
	Map()
}
