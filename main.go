package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/XciD/loxone-prometheus-exporter/config"

	loxone "github.com/XciD/loxone-ws"
	"github.com/bep/debounce"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	changes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loxone_changes",
			Help: "Number of changes",
		},
		[]string{"control", "room", "type", "cat", "state"},
	)
	values = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "loxone_values",
			Help: "Current Value of changes",
		},
		[]string{"control", "room", "type", "cat", "state"},
	)
)

func main() {
	ctx := context.Background()
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Read config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(err)
		return
	}

	// Start prometheus server
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)
	prometheus.MustRegister(changes)
	prometheus.MustRegister(values)

	// Open socket
	lox, err := loxone.New(cfg.Host, cfg.Port, cfg.User, cfg.Password)

	if err != nil {
		log.Error(err)
		return
	}

	// Get config
	loxoneConfig, err := lox.GetConfig()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Get Config OK")

	// Register events
	err = lox.RegisterEvents()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("RegisterEvents OK")

	// Build Control Map by states
	globalStates := make(map[string]*eventMetric)

	for _, control := range loxoneConfig.Controls {

		labels := map[string]string{
			"control": control.Name,
			"room":    loxoneConfig.RoomName(control.Room),
			"type":    control.Type,
			"cat":     loxoneConfig.CatName(control.Cat),
			"state":   "",
		}

		for stateName, stateValue := range control.States {
			// Can be a string or a float...
			switch stateValue := stateValue.(type) {
			case string:
				// Create the target map
				currentLabel := prometheus.Labels{}
				for key, value := range labels {
					currentLabel[key] = value
				}
				currentLabel["state"] = stateName
				globalStates[stateValue] = newEventMetric(&currentLabel)
			case []string:
				for index, childStateValue := range stateValue {
					// Create the target map
					currentLabel := prometheus.Labels{}
					for key, value := range labels {
						currentLabel[key] = value
					}
					currentLabel["state"] = stateName + "-" + string(index)
					globalStates[childStateValue] = newEventMetric(&currentLabel)
				}
			}
		}
	}

	for stateName, stateValue := range loxoneConfig.GlobalStates {
		currentLabel := prometheus.Labels{
			"control": "global",
			"room":    "global",
			"type":    "global",
			"cat":     "global",
			"state":   stateName,
		}
		globalStates[stateValue] = newEventMetric(&currentLabel)
	}

	log.Info("Start reading events")
	for {
		select {
		case <-ctx.Done():
			log.Infof("Shutting Down")
		case event := <-lox.GetEvents():
			if eventMetric, ok := globalStates[event.UUID]; ok {
				eventMetric.update(event.Value)
			} else {
				log.Debugf("event unknown: %+v\n", event)
			}
		}
	}
}

type eventMetric struct {
	labels           *prometheus.Labels
	initialized      bool
	debounceFunction func(f func())
}

func newEventMetric(labels *prometheus.Labels) *eventMetric {
	return &eventMetric{
		initialized:      false,
		labels:           labels,
		debounceFunction: debounce.New(500 * time.Millisecond),
	}
}

func (e *eventMetric) update(value float64) {
	values.With(*e.labels).Set(value)

	if !e.initialized {
		e.initialized = true
		return
	}

	log.Infof("New event %+v with value %f", e.labels, value)

	e.debounceFunction(func() {
		changes.With(*e.labels).Inc()
	})
}
