package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cloudfoundry/dropsonde"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/st3v/arp-watch/observer"
)

var (
	config *observer.Config
	state  *observer.State

	infoLog *log.Logger
	errLog  *log.Logger
)

func init() {
	infoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	configPath := flag.String("configPath", "", "Path to config file. Optional.")
	stateFilePath := flag.String("stateFilePath", "", "Path to state file. Optional.")
	flag.Parse()

	config = new(observer.Config)
	if *configPath != "" {
		if err := config.Load(*configPath); err != nil {
			errLog.Fatalf("Error loading config file: %s", err.Error())
		}
	}

	state = new(observer.State)
	if *stateFilePath != "" {
		if _, err := os.Stat(*stateFilePath); os.IsNotExist(err) {
			state.Write(*stateFilePath)
		}

		if err := state.Load(*stateFilePath); err != nil {
			errLog.Fatalf("Error loading state file: %s", err.Error())
		}

		setupExitHandler(*stateFilePath)

		defer state.Write(*stateFilePath)
	}

	handleFn := logEvent
	if config.Metron.Endpoint != "" && config.Metron.Origin != "" {
		err := dropsonde.Initialize(
			config.Metron.Endpoint,
			config.Metron.Origin,
		)
		if err != nil {
			errLog.Fatalf("Dropsonde failed to initialize", err)
		}

		handleFn = func(event observer.AddressChange) {
			emitMetric(event)
			logEvent(event)
		}
	}

	addrChanges := make(chan observer.AddressChange)
	done := make(chan struct{})
	go handleObservations(handleFn, addrChanges, done)

	observer.Observe(*config, state, addrChanges)
	<-done
}

func getAction(event observer.AddressChange) string {
	action := "changed"
	if event.Old == "" {
		action = "set"
	}

	if event.New == "" {
		action = "unset"
	}

	return action
}

func getKey(event observer.AddressChange, action string) string {
	return fmt.Sprintf(
		"net.arp.%s.%s",
		getName(event),
		action,
	)
}

func getName(event observer.AddressChange) string {
	return strings.Replace(event.Name, ".", "-", -1)
}

func emitMetric(event observer.AddressChange) {
	if event.New == "" {
		return
	}

	key := getKey(event, "changed")
	metrics.SendValue(key, 1, "count")
}

func logEvent(event observer.AddressChange) {
	infoLog.Printf("%s: '%s' -> '%s'\n", getKey(event, getAction(event)), event.Old, event.New)
}

func handleObservations(handleFn func(event observer.AddressChange), events chan observer.AddressChange, done chan struct{}) {
	for event := range events {
		handleFn(event)
	}

	close(done)
}

func saveState(path string) {
	state.Write(path)
}

func setupExitHandler(stateFilePath string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		state.Write(stateFilePath)
		os.Exit(1)
	}()
}
