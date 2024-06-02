// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/handlers"
	"github.com/scraly/http-go-server/pkg/swagger/server/restapi/operations"
)

var (
	handlerDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration in seconds",
	}, []string{"path", "method", "status"})
	handlerCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "http request count",
	}, []string{"path", "method", "status"})
)

func init() {
	prometheus.MustRegister(handlerDuration)
}

//go:generate swagger generate server --target ../../server --name Hello --spec ../../swagger.yml --exclude-main

func configureFlags(api *operations.HelloAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.HelloAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	if api.GetHelloUserHandler == nil {
		api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(func(params operations.GetHelloUserParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetHelloUser has not yet been implemented")
		})
	}

	if api.CheckHealthHandler == nil {
		api.CheckHealthHandler = operations.CheckHealthHandlerFunc(func(params operations.CheckHealthParams) middleware.Responder {
			return middleware.NotImplemented("operation .CheckHealth has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	promHandler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
	err := prometheus.Register(handlerDuration)
	if err != nil {
		if err.Error() != "duplicate metrics collector registration attempted" {
			log.Fatalf("Failed to register handler duration: %s", err)
		}
	}
	err = prometheus.Register(handlerCount)
	if err != nil {
		if err.Error() != "duplicate metrics collector registration attempted" {
			log.Fatalf("Failed to register handler count: %s", err)
		}
	}
	return &metricsHandler{
		handler:        handler,
		promHandler:    promHandler,
		loggingHandler: handlers.LoggingHandler(os.Stdout, handler),
	}
}

type metricsHandler struct {
	handler        http.Handler
	promHandler    http.Handler
	loggingHandler http.Handler
}

func (m *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (strings.ToUpper(r.Method) == "GET") && (strings.ToLower(r.URL.Path) == "/metrics") {

		// Let Prometheus have it
		m.promHandler.ServeHTTP(w, r)
	} else {
		if r.URL.Path == "/healthz" {
			m.handler.ServeHTTP(w, r)
		} else {
			start := time.Now()
			lrw := newMetricResponseWriter(w)
			// Pass it on to the original handler
			m.loggingHandler.ServeHTTP(lrw, r)
			statusCode := lrw.statusCode
			duration := time.Since(start)
			handlerCount.WithLabelValues("/hello/{user}", strings.ToLower(r.Method), fmt.Sprintf("%d", statusCode)).Inc()
			handlerDuration.WithLabelValues("/hello/{user}", strings.ToLower(r.Method), fmt.Sprintf("%d", statusCode)).Observe(duration.Seconds())
		}
	}
}

type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newMetricResponseWriter(w http.ResponseWriter) *metricResponseWriter {
	return &metricResponseWriter{w, http.StatusOK}
}

func (lrw *metricResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
