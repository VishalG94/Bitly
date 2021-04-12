package main

import (
	// "bytes"
	// "encoding/json"
	// "os"
	// "strings"
	"context"
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	// "github.com/streadway/amqp"
	"github.com/unrolled/render"
	"github.com/rs/cors"
	"time"
)

// RabbitMQ Config
var rabbitmq_server = "10.0.3.80"
var rabbitmq_port = "5672"
var createShortLink_queue = "shortlinks"
var hitShortLink_queue = "hitshortlinks"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

var nosql_loadbalancer = "http://NoSqlNLB-d36ca151fc19703a.elb.us-west-2.amazonaws.com:9000/api/"

var (
	username = "cmpe281"
	password = "cmpe281"
	hostname = "10.0.1.76"
	port 	 = "4306"
	dbname   = "bitly"
)
var linksCount = 20
// Singleton MySQ:L Object
var dbConn *sql.DB = nil


// Returns data source name (dsn)
func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbName)
}

// Initialized the database connection
func GetDBInstance() (*sql.DB, error){
	if dbConn == nil {
		db, err := sql.Open("mysql", dsn(dbname))
		if err != nil {
			log.Fatal(err)
		} else {
			db.SetMaxOpenConns(20)
			db.SetMaxIdleConns(20)
			db.SetConnMaxLifetime(time.Minute * 5)
			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			err = db.PingContext(ctx)
			if err != nil {
				log.Printf("Errors %s pinging DB", err)
				return nil, err
			}
			log.Printf("Connected to DB %s successfully\n", dbname)
			dbConn = db
		}
	}
	return dbConn, nil
}

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.Use(c)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/shortlinktrend", trendingLinks(formatter)).Methods("GET")
	// mx.HandleFunc("/shortlinktrend", shortlinkCreateHandler(formatter)).Methods("POST")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Trend Server API version 1.0 alive!"})
	}
}

// API Gumball Machine Handler
func trendingLinks(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		db, err := GetDBInstance()
		if err != nil {
			log.Print("DB error, fatal handler.")
			log.Fatal(err)
		} else {

			// var (
			// 	id     		string
			// 	shortlink  	string
			// 	url  		string
			// 	count 		int
			// )
			
			rows, err := db.Query("SELECT * FROM shortlinks ORDER BY count DESC LIMIT ?;", linksCount)
			log.Println("Get Trend Links!")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			var shortlnks []ShortLinks
			for rows.Next() {
				var shortlnk = ShortLinks{}
				err := rows.Scan(&shortlnk.Id, &shortlnk.URL, &shortlnk.ShortLink , &shortlnk.Count)
				if err != nil {
					log.Fatal(err)
				}
				shortlnks = append(shortlnks, shortlnk)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			formatter.JSON(w, http.StatusOK, shortlnks)
		}
	}
}

