// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package countconnector // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/countconnector"

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.uber.org/zap"
)

// Default metrics are emitted if no conditions are specified.
const (
	defaultMetricNameSpans      = "trace.span.count"
	defaultMetricDescSpans      = "The number of spans observed."
	defaultMetricNameSpanEvents = "trace.span.event.count"
	defaultMetricDescSpanEvents = "The number of span events observed."

	defaultMetricNameMetrics    = "metric.count"
	defaultMetricDescMetrics    = "The number of metrics observed."
	defaultMetricNameDataPoints = "metric.datapoint.count"
	defaultMetricDescDataPoints = "The number of data points observed."

	defaultMetricNameLogs = "log.record.count"
	defaultMetricDescLogs = "The number of log records observed."
)

// Config for the connector
type Config struct {
	Spans      map[string]MetricInfo `mapstructure:"spans"`
	SpanEvents map[string]MetricInfo `mapstructure:"spanevents"`
	Metrics    map[string]MetricInfo `mapstructure:"metrics"`
	DataPoints map[string]MetricInfo `mapstructure:"datapoints"`
	Logs       map[string]MetricInfo `mapstructure:"logs"`
}

// MetricInfo for a data type
type MetricInfo struct {
	Description string   `mapstructure:"description"`
	Conditions  []string `mapstructure:"conditions"`
}

func (c *Config) Validate() error {
	for name, info := range c.Spans {
		if name == "" {
			return fmt.Errorf("spans: metric name missing")
		}
		parser, err := newSpanParser(zap.NewNop())
		if err != nil {
			return err
		}
		if _, err = parseConditions(parser, info.Conditions); err != nil {
			return fmt.Errorf("spans condition: metric %q: %w", name, err)
		}
	}
	for name, info := range c.SpanEvents {
		if name == "" {
			return fmt.Errorf("spanevents: metric name missing")
		}
		parser, err := newSpanEventParser(zap.NewNop())
		if err != nil {
			return err
		}
		if _, err = parseConditions(parser, info.Conditions); err != nil {
			return fmt.Errorf("spanevents condition: metric %q: %w", name, err)
		}
	}
	for name, info := range c.Metrics {
		if name == "" {
			return fmt.Errorf("metrics: metric name missing")
		}
		parser, err := newMetricParser(zap.NewNop())
		if err != nil {
			return err
		}
		if _, err = parseConditions(parser, info.Conditions); err != nil {
			return fmt.Errorf("metrics condition: metric %q: %w", name, err)
		}
	}

	for name, info := range c.DataPoints {
		if name == "" {
			return fmt.Errorf("datapoints: metric name missing")
		}
		parser, err := newDataPointParser(zap.NewNop())
		if err != nil {
			return err
		}
		if _, err = parseConditions(parser, info.Conditions); err != nil {
			return fmt.Errorf("datapoints condition: metric %q: %w", name, err)
		}
	}
	for name, info := range c.Logs {
		if name == "" {
			return fmt.Errorf("logs: metric name missing")
		}
		parser, err := newLogParser(zap.NewNop())
		if err != nil {
			return err
		}
		if _, err = parseConditions(parser, info.Conditions); err != nil {
			return fmt.Errorf("logs condition: metric %q: %w", name, err)
		}
	}
	return nil
}

var _ confmap.Unmarshaler = (*Config)(nil)

// Unmarshal with custom logic to set default values.
// This is necessary to ensure that default metrics are
// not configured if the user has specified any custom metrics.
func (c *Config) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		// Nothing to do if there is no config given.
		return nil
	}
	if err := componentParser.Unmarshal(c); err != nil {
		return err
	}
	if !componentParser.IsSet("spans") {
		c.Spans = defaultSpansConfig()
	}
	if !componentParser.IsSet("spanevents") {
		c.SpanEvents = defaultSpanEventsConfig()
	}
	if !componentParser.IsSet("metrics") {
		c.Metrics = defaultMetricsConfig()
	}
	if !componentParser.IsSet("datapoints") {
		c.DataPoints = defaultDataPointsConfig()
	}
	if !componentParser.IsSet("logs") {
		c.Logs = defaultLogsConfig()
	}
	return nil
}

func defaultSpansConfig() map[string]MetricInfo {
	return map[string]MetricInfo{
		defaultMetricNameSpans: {
			Description: defaultMetricDescSpans,
		},
	}
}

func defaultSpanEventsConfig() map[string]MetricInfo {
	return map[string]MetricInfo{
		defaultMetricNameSpanEvents: {
			Description: defaultMetricDescSpanEvents,
		},
	}
}

func defaultMetricsConfig() map[string]MetricInfo {
	return map[string]MetricInfo{
		defaultMetricNameMetrics: {
			Description: defaultMetricDescMetrics,
		},
	}
}

func defaultDataPointsConfig() map[string]MetricInfo {
	return map[string]MetricInfo{
		defaultMetricNameDataPoints: {
			Description: defaultMetricDescDataPoints,
		},
	}
}

func defaultLogsConfig() map[string]MetricInfo {
	return map[string]MetricInfo{
		defaultMetricNameLogs: {
			Description: defaultMetricDescLogs,
		},
	}
}
