package appinsights_wrapper

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
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

func (tc *TelemetryClient) StartOperation(name string) {
	tc.Context().Tags.Operation().SetId(newUUID().String())
	tc.Context().Tags.Operation().SetName(name)
	fmt.Printf("\nSTART OPERATION | ID:%s\n", tc.Context().Tags.Operation().GetId())
}

func (tc *TelemetryClient) EndOperation() {
	fmt.Printf("\nEND OPERATION | ID:%s\n", tc.Context().Tags.Operation().GetId())
	for k := range tc.Context().Tags.Operation() {
		delete(tc.Context().Tags.Operation(), k)
	}
}

func (tc *TelemetryClient) Log(telemetry appinsights.Telemetry) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	tc.Context().CommonProperties["file"] = file
	tc.Context().CommonProperties["line"] = fmt.Sprint(line)

	tc.Track(telemetry)
}

func (tc *TelemetryClient) LogEvent(name string) {
	tc.Log(appinsights.NewEventTelemetry(name))
}

func (tc *TelemetryClient) LogMetric(name string, value float64) {
	tc.Log(appinsights.NewMetricTelemetry(name, value))
}

func (tc *TelemetryClient) LogTrace(message string, severity contracts.SeverityLevel) {
	tc.Log(appinsights.NewTraceTelemetry(message, severity))
}

func (tc *TelemetryClient) LogRequest(method, url string, duration time.Duration, responseCode string) {
	tc.Log(appinsights.NewRequestTelemetry(method, url, duration, responseCode))
}

func (tc *TelemetryClient) LogRemoteDependency(name, dependencyType, target string, success bool) {
	tc.Log(appinsights.NewRemoteDependencyTelemetry(name, dependencyType, target, success))
}

func (tc *TelemetryClient) LogAvailability(name string, duration time.Duration, success bool) {
	tc.Log(appinsights.NewAvailabilityTelemetry(name, duration, success))
}

func (tc *TelemetryClient) LogException(err interface{}) {
	tc.Log(appinsights.NewExceptionTelemetry(err))
}
