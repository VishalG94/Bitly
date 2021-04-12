package main

import (
	//"bytes"
	"context"
	"encoding/json"
	"fmt"
	//uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	//"github.com/catinello/base62"
	"database/sql"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
	"github.com/rs/cors"
	"strings"
	"time"
)

// MongoDB Config
//var mongodb_server = "mongodb"
//var mongodb_database = "cmpe281"
//var mongodb_collection = "gumball"

// RabbitMQ Config
var rabbitmq_server = "10.0.3.80"
var rabbitmq_port = "5672"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"
var createShortLink_queue = "shortlinks"
var hitShortLink_queue = "hitshortlinks"

var nosql_loadbalancer = "http://NoSqlNLB-d36ca151fc19703a.elb.us-west-2.amazonaws.com:9000/api/"

var (
	username = "cmpe281"
	password = "cmpe281"
	hostname = "10.0.1.76"
	port 	 = "4306"
	dbname   = "bitly"
)

//var (
	//currCntrRange Range 	= Range{}
	//count         uint64   	= 0
	//gatewayIp           	= "http://34.221.223.237:8000"
	//redirectPath        	= "/lr/"
	//counterServerIp 		= "http://10.0.1.199:3002/basecount"
//)

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
	mx.HandleFunc("/{shortLink}", shortlinkRedirect(formatter)).Methods("GET")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

//// Get the new range from counter server
//func getCounterRange() error {
//	w, err := http.Get(counterServerIp)
//	if err != nil {
//		failOnError(err, "Counter server connection failed")
//	}
//	decoder := json.NewDecoder(w.Body)
//	var cntrRange Range
//	err = decoder.Decode(&cntrRange)
//	if err != nil {
//		log.Fatal(err)
//	}
//	count = cntrRange.Min
//	return err
//}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Link Redirect API version 1.0 alive!"})
	}
}

// API Gumball Machine Handler
func shortlinkRedirect(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var shortlnk ShortLinks
		shortlnk.ShortLink = "http://" + req.Host + req.URL.Path
		log.Println("short Link: ", shortlnk.ShortLink)
		splits := strings.Split(shortlnk.ShortLink, "/")
		if( len(splits)<1 ){
			log.Fatal("Short Link Dose not exist!")
		}
		docKey := splits[len(splits)-1]
		log.Println("key: " + docKey)
		res, err := http.Get(nosql_loadbalancer+docKey)
		if err!=nil {
			log.Println(err, "Failed to fetch data form NoSQL cache; ", err.Error())
		}
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			failOnError(err, "Failed while reading response body")
		}
		err = json.Unmarshal(responseBody, &shortlnk)
		if err!=nil {
			//failOnError(err, "Unmarshal failed for responseBody")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error: failed to get shortlink details"}`))
			return
		}
		queue_send(shortlnk, hitShortLink_queue)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, []byte(`{"error": "Error: failed to get shortlink details"}`))
		}
		log.Println("Redirect link: ",shortlnk.URL)
		http.Redirect(w, req, shortlnk.URL, http.StatusSeeOther)
	}
}

// API Create New short Link
// func shortlinkCreateHandler(formatter *render.Render) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 			uuid := uuid.NewV4()
// 			decoder := json.NewDecoder(req.Body)
// 			var shortlnk ShortLinks
// 			var (
// 				id     		string
// 				shortlink  	string
// 				url  		string
// 			)
// 			err := decoder.Decode(&shortlnk)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			url = shortlnk.URL
// 			id = uuid.String()
// 			shortlnk.Id = id
// 			shortlnk.Count=0

// 			// Hashing
// 			var uriHash string
// 			if url == "" {
// 				log.Fatal("URL not present in the request!")
// 			}
// 			if count >= currCntrRange.Max {
// 				err := getCounterRange()
// 				if err != nil {
// 					failOnError(err, "Error encountered in counter server")
// 				}
// 			}
// 			uriHash = base62.Encode(int(count))
// 			count++
// 			shortlink = gatewayIp + redirectPath + uriHash
// 			shortlnk.ShortLink = shortlink

// 			// Publishing in queue
// 			queue_send(shortlnk, createShortLink_queue)

// 			formatter.JSON(w, http.StatusOK, shortlnk)
// 	}
// }

func queue_send(message ShortLinks, queueName string) {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//body := message
	body, err := json.Marshal(message)
	if err != nil {
		failOnError(err, "Error encoding JSON")
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
			//ContentType: "text/plain",
			//Body:        []byte(body),
		})
	if err != nil {
		failOnError(err, "Error encoding JSON")
	}
}

//Receive Order from Queue to Process
func queue_receive() []string {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		createShortLink_queue, // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"orders",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	shortlink_ids := make(chan string)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			shortlink_ids <- string(d.Body)
		}
		close(shortlink_ids)
	}()

	err = ch.Cancel("orders", false)
	if err != nil {
	    log.Fatalf("basic.cancel: %v", err)
	}

	var shortlink_ids_array []string
	for n := range shortlink_ids {
 	shortlink_ids_array = append(shortlink_ids_array, n)
 }

 return shortlink_ids_array
}
