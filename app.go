package main

import (
	"net/http"
	"os"

	"github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/Financial-Times/public-people-api/people"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Financial-Times/up-headlines-api/headlines"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rcrowley/go-metrics"
)

func main() {
	app := cli.App("public-headlines-api", "A public RESTful API for accessing the Top Headlines")
	mongoURL := app.String(cli.StringOpt{
		Name:   "mongo-url",
		Value:  "localhost:27017",
		Desc:   "neo4j endpoint URL",
		EnvVar: "MONGO_URL",
	})
	logLevel := app.String(cli.StringOpt{
		Name:   "log-level",
		Value:  "INFO",
		Desc:   "Log level to use",
		EnvVar: "LOG_LEVEL",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})
	listURL := app.String(cli.StringOpt{
		Name:   "list-url",
		Value:  "",
		Desc:   "List URL",
		EnvVar: "LISTS_URL",
	})
	conceptURL := app.String(cli.StringOpt{
		Name:   "concept-url",
		Value:  "",
		Desc:   "Content By Concept URL",
		EnvVar: "CONCEPT_URL",
	})

	app.Action = func() {
		parsedLogLevel, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.WithFields(log.Fields{"logLevel": logLevel, "err": err}).Fatal("Incorrect log level")
		}
		log.SetLevel(parsedLogLevel)

		log.Infof("public-headlines-api will listen on port: %s, connecting to: %s", *port, *mongoURL)

		headlineService := headlines.NewHeadlineService(*mongoURL, *listURL, *conceptURL)

		handler := headlines.NewHeadlineHandler(headlineService)

		servicesRouter := mux.NewRouter()

		// Then API specific ones
		servicesRouter.HandleFunc("/headlines/flash/{uuid}", handler.GetFlashBriefing).Methods("GET")
		servicesRouter.HandleFunc("/headlines/list/{uuid}", handler.GetListHeadlines).Methods("GET")
		servicesRouter.HandleFunc("/headlines/concept/{uuid}", handler.GetConceptHeadlines).Methods("GET")
		servicesRouter.HandleFunc("/headlines", handler.GetHeadlinesByUUID).Methods("POST")

		var monitoringRouter http.Handler = servicesRouter
		monitoringRouter = httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), monitoringRouter)
		monitoringRouter = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, monitoringRouter)

		// The following endpoints should not be monitored or logged (varnish calls one of these every second, depending on config)
		// The top one of these build info endpoints feels more correct, but the lower one matches what we have in Dropwizard,
		// so it's what apps expect currently same as ping, the content of build-info needs more definition
		http.HandleFunc(status.PingPath, status.PingHandler)
		http.HandleFunc(status.PingPathDW, status.PingHandler)
		http.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
		http.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)
		http.HandleFunc("/__gtg", people.GoodToGo)
		http.Handle("/", monitoringRouter)

		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}
