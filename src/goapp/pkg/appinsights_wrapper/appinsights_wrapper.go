package appinsights_wrapper

import (
	"fmt"
	"log"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var (
	tc *appinsights.TelemetryConfiguration
)

type TelemetryClient struct {
	appinsights.TelemetryClient
}

func Init(instrumentationKey string) {
	tc = appinsights.NewTelemetryConfiguration(instrumentationKey)

	/*turn on diagnostics to help troubleshoot problems with telemetry submission. */
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})
}

func NewClient() *TelemetryClient {
	return &TelemetryClient{
		TelemetryClient: appinsights.NewTelemetryClientFromConfig(tc),
	}
}

func (c *TelemetryClient) StartOperation(name string) {
	c.Context().Tags.Operation().SetId(newUUID().String())
	c.Context().Tags.Operation().SetName(name)
	fmt.Printf("\nSTART OPERATION | ID:%s\n", c.Context().Tags.Operation().GetId())
}

func (c *TelemetryClient) EndOperation() {
	fmt.Printf("\nEND OPERATION | ID:%s\n", c.Context().Tags.Operation().GetId())
	for k := range c.Context().Tags.Operation() {
		delete(c.Context().Tags.Operation(), k)
	}
}
