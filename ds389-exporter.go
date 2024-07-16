package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"main/src"
	"net/http"
	"time"
)

// Meta information
const exporterVersion = "1.0"

// Global logger
var logger = *logrus.New()

// Metrics description
var metrics = []*src.Metric{
	{
		Name:        "ds389_scape_duration_millis",
		Help:        "Duration of the request to obtain metrics in milliseconds",
		MetricType:  src.Gauge,
		FindKeyword: "version",
	},
	{
		Name:        "ds389_version",
		Help:        "Version 389 Directory Server",
		MetricType:  src.GaugeVec,
		FindKeyword: "version",
	},
	{
		Name:        "ds389_uptime",
		Help:        "Uptime in seconds",
		MetricType:  src.Gauge,
		FindKeyword: "starttime",
	},
	{
		Name:        "ds389_startup_timestamp",
		Help:        "Server startup unix timestamp",
		MetricType:  src.Gauge,
		FindKeyword: "starttime",
	},
	{
		Name:        "ds389_current_connections",
		Help:        "Number of active connections to the server",
		MetricType:  src.Gauge,
		FindKeyword: "currentconnections",
	},
	{
		Name:        "ds389_total_connections",
		Help:        "Total number of connections made to the server",
		MetricType:  src.Counter,
		FindKeyword: "totalconnections",
	},
	{
		Name:        "ds389_max_threads_per_conn_hits",
		Help:        "Total number of times the maximum threads per connection limit has been reached",
		MetricType:  src.Counter,
		FindKeyword: "maxthreadsperconnhits",
	},
	{
		Name:        "ds389_connections_in_max_threads",
		Help:        "Number of connections currently utilizing the maximum allowed threads",
		MetricType:  src.Counter,
		FindKeyword: "connectionsinmaxthreads",
	},
	{
		Name:        "ds389_connections_max_threads",
		Help:        "Total number of concurrent connections that the LDAP server can service simultaneously using the maximum number of threads",
		MetricType:  src.Counter,
		FindKeyword: "connectionsmaxthreadscount",
	},
	{
		Name:        "ds389_current_connections_at_max_threads",
		Help:        "Number of connections currently utilizing the maximum allowed threads per connection",
		MetricType:  src.Gauge,
		FindKeyword: "currentconnectionsatmaxthreads",
	},
	{
		Name:        "ds389_read_waiters",
		Help:        "Number of threads waiting for read operations to complete",
		MetricType:  src.Gauge,
		FindKeyword: "readwaiters",
	},
	{
		Name:        "ds389_threads",
		Help:        "Total number of threads currently active in the server",
		MetricType:  src.Gauge,
		FindKeyword: "threads",
	},
	{
		Name:        "ds389_operations",
		Help:        "Total number of entry operations by type processed by the server",
		MetricType:  src.CounterVec,
		FindKeyword: "addentryops",
		WithLabels: []src.WithLabels{
			{
				Key:         "type",
				Value:       "SEARCH",
				FindKeyword: "searchops",
			},
			{
				Key:         "type",
				Value:       "ADD",
				FindKeyword: "addentryops",
			},

			{
				Key:         "type",
				Value:       "MODIFY",
				FindKeyword: "modifyentryops",
			},
			{
				Key:         "type",
				Value:       "MODIFY_RND",
				FindKeyword: "modifyrdnops",
			},
			{
				Key:         "type",
				Value:       "COMPARE",
				FindKeyword: "compareops",
			},
			{
				Key:         "type",
				Value:       "DELETE",
				FindKeyword: "removeentryops",
			},
		},
	},
	{
		Name:        "ds389_bind_operations",
		Help:        "Total number of bind operations performed on the server",
		MetricType:  src.CounterVec,
		FindKeyword: "anonymousbinds",
		WithLabels: []src.WithLabels{
			{
				Key:         "type",
				Value:       "ANONYMOUS",
				FindKeyword: "anonymousbinds",
			},
			{
				Key:         "type",
				Value:       "UNAUTHORIZED",
				FindKeyword: "unauthbinds",
			},
			{
				Key:         "type",
				Value:       "SIMPLE",
				FindKeyword: "simpleauthbinds",
			},
			{
				Key:         "type",
				Value:       "STRONG",
				FindKeyword: "strongauthbinds",
			},
		},
	},
	{
		Name:        "ds389_bind_errors",
		Help:        "Total number of bind operations rejected due to security-related errors",
		MetricType:  src.Counter,
		FindKeyword: "bindsecurityerrors",
	},

	{
		Name:        "ds389_rx_bytes",
		Help:        "Total number of bytes received by the server",
		MetricType:  src.Counter,
		FindKeyword: "bytesrecv",
	},
	{
		Name:        "ds389_tx_bytes",
		Help:        "Total number of bytes transferred by the server",
		MetricType:  src.Counter,
		FindKeyword: "bytessent",
	},
	{
		Name:        "ds389_returned_entries",
		Help:        "Total number of entries returned by search operations on the server",
		MetricType:  src.Counter,
		FindKeyword: "entriesreturned",
	},
	{
		Name:        "ds389_completed_operations",
		Help:        "Total number of directory operations completed by the server",
		MetricType:  src.Counter,
		FindKeyword: "opscompleted",
	},
	{
		Name:        "ds389_initiated_operations",
		Help:        "Total number of directory operations initiated by clients and processed by the server",
		MetricType:  src.Counter,
		FindKeyword: "opsinitiated",
	},
	{
		Name:        "ds389_cache_entries",
		Help:        "Number of entries currently cached in the server's cache",
		MetricType:  src.Gauge,
		FindKeyword: "cacheentries",
	},
	{
		Name:        "ds389_cache_hits_count",
		Help:        "Total number of times entries have been found in the cache and returned without accessing the backend storage",
		MetricType:  src.Counter,
		FindKeyword: "cachehits",
	},
	{
		Name:        "ds389_dtable_size",
		Help:        "Size of the Directory Server descriptor table",
		MetricType:  src.Gauge,
		FindKeyword: "dtablesize",
	},
	{
		Name:        "ds389_search_operations_level",
		Help:        "Total number of search operations by level executed by the server",
		MetricType:  src.CounterVec,
		FindKeyword: "anonymousbinds",
		WithLabels: []src.WithLabels{
			{
				Key:         "type",
				Value:       "ONE",
				FindKeyword: "onelevelsearchops",
			},
			{
				Key:         "type",
				Value:       "SUBTREE",
				FindKeyword: "wholesubtreesearchops",
			},
		},
	},
	{
		Name:        "ds389_copy_entries",
		Help:        "Total number of operations to copy or move entries between different containers or subsections",
		MetricType:  src.Counter,
		FindKeyword: "copyentries",
	},
	{
		Name:        "ds389_errors",
		Help:        "Total number of errors that occurred on the server",
		MetricType:  src.Counter,
		FindKeyword: "errors",
	},
	{
		Name:        "ds389_number_backends",
		Help:        "Number of alternative backends or data stores in use that are integrated and used by the LDAP server",
		MetricType:  src.Gauge,
		FindKeyword: "nbackends",
	},
	{
		Name:        "ds389_returned_referrals",
		Help:        "Total number of referrals returned to the client in response to its request to the LDAP server",
		MetricType:  src.Counter,
		FindKeyword: "referralsreturned",
	},
	{
		Name:        "ds389_security_errors",
		Help:        "Total number of security errors that occurred on the server",
		MetricType:  src.Counter,
		FindKeyword: "securityerrors",
	},
	{
		Name:        "ds389_supplier_entries",
		Help:        "Total number of entries that were received from the data provider",
		MetricType:  src.Counter,
		FindKeyword: "supplierentries",
	},
	{
		Name:        "ds389_consumer_hits",
		Help:        "Total number of times entries from the supplier's database have been accessed by consumer servers",
		MetricType:  src.Counter,
		FindKeyword: "consumerhits",
	},
}

func main() {

	// Parse arguments
	configFile := flag.String("config.file", "./config.yml", "Path to configuration file")
	webTelemetryPath := flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	webListenAddress := flag.String("web.listen.address", "0.0.0.0:9389", "Address to listen on for get telemetry")
	logLevel := flag.String("log.level", "info", "Only log messages with the given severity or above. One of: [debug, info, warn, error, fatal]")
	flag.Parse()

	// Init logger
	initLogger(*logLevel)

	logger.Infof("Starting ds389-exporter %s", exporterVersion)

	// Read configuration
	config, err := src.GetConfig(*configFile)
	if err != nil {
		logger.Fatalf("Failed to read config file: %v", err)
	}
	logger.Infof("Using config file: %s", *configFile)

	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, config, *webTelemetryPath)
	})

	logger.Infof("Metrics are available at http://%s%s", *webListenAddress, *webTelemetryPath)
	if err := http.ListenAndServe(*webListenAddress, server); err != nil {
		logger.Fatal(err)
	}
}

// Function to initialize the logger with a given level
func initLogger(level string) {
	logger.SetFormatter(&logrus.TextFormatter{DisableColors: true, FullTimestamp: true})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(lvl)
}

// Function for processing an incoming HTTP request
func handleRequest(w http.ResponseWriter, r *http.Request, config *src.Config, webTelemetryPath string) {
	if r.URL.Path != webTelemetryPath {
		notFoundHandler(w, r, webTelemetryPath)
		return
	}

	registry := prometheus.NewRegistry()
	collector := src.LDAPCollector{
		LdapURL:      "ldap://" + config.Server + ":" + config.Port,
		BindDN:       config.BindDN,
		BindPassword: config.BindPassword,
	}

	start := time.Now()
	result, err := collector.Collect()
	duration := time.Since(start).Milliseconds()
	if err != nil {
		logger.Errorf("Failed to get metrics from 389 Directory Server: %v", err)
		http.Error(w, "Failed to get metrics from 389 Directory Server", http.StatusInternalServerError)
		return
	}

	for _, metric := range metrics {
		processMetric(registry, metric, result, duration)
	}

	logger.WithFields(logrus.Fields{"remote_addr": r.RemoteAddr, "method": r.Method, "path": r.URL.Path}).Info("Success")
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

// Function for processing a non-existent page
func notFoundHandler(w http.ResponseWriter, r *http.Request, webTelemetryPath string) {
	logger.WithFields(logrus.Fields{"remote_addr": r.RemoteAddr, "method": r.Method, "path": r.URL.Path}).Info("Page not found")
	body := `<html>
             <head><title>389 Directory Server Exporter</title></head>
             <body>
             <h1>389 Directory Server Exporter</h1>
             <p>For the metrics visit: <a href='` + webTelemetryPath + `'>` + webTelemetryPath + `</a></p>
             </body>
             </html>`
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(body))
	if err != nil {
		logger.Error("Failed to respond to HTTP request")
		return
	}
}

// Function converting a set into metrics OpenMetrics format
func processMetric(registry *prometheus.Registry, metric *src.Metric, result map[string]string, duration int64) {
	stringValue, value, _ := src.FindValue(result, metric.FindKeyword)
	labelNames := []string{}
	labels := []map[string]string{}
	values := []float64{}

	if metric.Name == "ds389_scape_duration_millis" {
		value = float64(duration)
	}

	if metric.Name == "ds389_version" {
		labelNames = []string{"version"}
		labels = append(labels, map[string]string{"version": stringValue})
	}

	if metric.Name == "ds389_uptime" || metric.Name == "ds389_startup_timestamp" {
		pattern := "20060102150405Z"
		t, err := time.Parse(pattern, stringValue)
		if err != nil {
			logger.Error(err)
			return
		}

		if metric.Name == "ds389_uptime" {
			value = float64(time.Now().Unix() - t.Unix())
		}

		if metric.Name == "ds389_startup_timestamp" {
			value = float64(t.Unix())
		}
	}

	if metric.WithLabels != nil {
		labelNames = []string{"type"}
		for _, label := range metric.WithLabels {
			stringValue, value, _ = src.FindValue(result, label.FindKeyword)
			labels = append(labels, map[string]string{label.Key: label.Value})
			values = append(values, value)
		}
	}

	switch metric.MetricType {
	case src.Gauge:
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: metric.Name, Help: metric.Help})
		gauge.Set(value)
		registry.MustRegister(gauge)
	case src.GaugeVec:
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metric.Name, Help: metric.Help}, labelNames)
		gauge.With(labels[0]).Set(1)
		registry.MustRegister(gauge)
	case src.Counter:
		counter := prometheus.NewCounter(prometheus.CounterOpts{Name: metric.Name, Help: metric.Help})
		counter.Add(value)
		registry.MustRegister(counter)
	case src.CounterVec:
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: metric.Name, Help: metric.Help}, labelNames)
		for i, label := range labels {
			counter.With(label).Add(values[i])
		}
		registry.MustRegister(counter)
	}
}
