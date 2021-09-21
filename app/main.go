package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"test-mongo/app/config"
	"test-mongo/app/process"
	"test-mongo/app/storage/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dbMongo              *mongo.Mongo
	corsAllowHeaders     = "Content-Type,Bearer,content-type,Origin,Accept,Access-Control-Allow-Headers,Access-Control-Allow-Origin,Authorization,X-Requested-With,WebviewToken"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,PATCH,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func init() {
	dbMongo = getMongoConnection()
}

func getMongoConnection() *mongo.Mongo {
	var mongoURL = func() string {
		return fmt.Sprintf("%s://%s:%s@%s/%s?%s", config.MongoDBProtocol, config.MongoDBUsername, config.MongoDBPassword, config.MongoDBHost, config.MongoDBName, config.MongoDBOptions)
	}()
	m, err := mongo.New(mongoURL)
	if err != nil {
		log.Println("Connect database failed::", err.Error())
		os.Exit(1)
	}

	log.Println("Connect to mongo successfully")
	return m
}

func main() {
	defer dbMongo.Close(context.TODO())
	p, err := process.New(dbMongo)
	if err != nil {
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(cors)
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)

	router.HandleFunc("/healthz", health)

	router.With(middleware.Logger).Get("/pages/{page-id}/templates", p.GetAllTemplate)
	router.With(middleware.Logger).Post("/pages/{page-id}/templates", p.CreateTemplate)
	router.With(middleware.Logger).Put("/pages/{page-id}/templates/{template-id}", p.UpdateTemplate)
	router.With(middleware.Logger).Delete("/pages/{page-id}/templates/{template-id}", p.DeleteTemplate)

	port := ":" + config.Port
	http.ListenAndServe(port, router)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
		w.Header().Set("Content-Type", "application/json")
		if len(w.Header().Get("Ct")) > 0 {
			w.Header().Add("Content-Type", w.Header().Get("Ct"))
			w.Header().Del("Ct")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Max-Age", "5000")
	w.Header().Set("Content-Type", "text/plain charset=UTF-8")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusNoContent)
}

func health(w http.ResponseWriter, r *http.Request) {
}
