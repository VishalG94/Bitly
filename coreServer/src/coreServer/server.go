package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"context"
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
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
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Core Server API version 1.0 alive!"})
	}
}

// Send Order to Queue for Processing
func queue_send(message string) {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		createShortLink_queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}


func saveShortLinkToCache(shortlnk ShortLinks) {
	splits := strings.Split(shortlnk.ShortLink, "/")
	if( len(splits)<1 ){
		log.Fatal("Short Link Dose not exist!")
	}
	docKey := splits[len(splits)-1]
	log.Println("key: " + docKey)

	responseBody, err := json.Marshal(shortlnk)
	w, err := http.Post(nosql_loadbalancer+docKey, "application/json", bytes.NewBuffer(responseBody))
	if (err != nil) || (w.StatusCode != http.StatusOK) {
		failOnError(err, "Error encountered when saving to NoSQL database.")
		if w != nil {
			fmt.Println("Error encountered when saving to NoSQL database.", w.Body)
		}
	}
	log.Println("Successfully saved to NoSQL database,s", docKey)
}

// Save ShortLink To MySQL
func saveShortLinkToMySQL(shortlnk ShortLinks) {
	db, err := GetDBInstance()
	if err != nil {
		log.Print("DB error, fatal handler.")
		log.Fatal(err)
	} else {
		// Insert into the database
		insShtLnk, err := db.Prepare("insert into shortlinks(id, url, shortlink, count) VALUES(?,?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		insShtLnk.Exec(shortlnk.Id, shortlnk.URL, shortlnk.ShortLink, shortlnk.Count)
		log.Println("Successfully saved to MySQL database.")
	}
}

func create_link_queue_receive() {
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

	if err!=nil {
		failOnError(err, "Failed to register a consumer")	
	}

	var goroutineDelta = make(chan int)
	go func() {
		for d := range msgs {
			var shortlnk ShortLinks
			err = json.Unmarshal(d.Body, &shortlnk)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
		    saveShortLinkToCache(shortlnk)
			saveShortLinkToMySQL(shortlnk)	
			}
		// }
	}()

	// Avoid the go routine from exiting
	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		if numGoroutines == 0 { os.Exit(0) }
	}
}

func updateShortLinkToCache(shortlnk ShortLinks) {
	splits := strings.Split(shortlnk.ShortLink, "/")
	if( len(splits)<1 ){
		log.Fatal("Short Link Dose not exist!")
	}
	docKey := splits[len(splits)-1]
	log.Println("key: " + docKey)
	client := &http.Client{}
	responseBody, err := json.Marshal(shortlnk)
	req, err := http.NewRequest(http.MethodPut, nosql_loadbalancer+docKey, bytes.NewBuffer(responseBody))
	if (err != nil) {
		failOnError(err, "Error encountered when saving to NoSQL database.")
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(req)
	if (err != nil) || (res.StatusCode != http.StatusOK) {
		if res != nil {
			res.Body.Close()
		}
		return
	}
	log.Println("Successfully saved to NoSQL database,s", shortlnk.Id)
}

// Update ShortLink To MySQL
func updateShortLinkToMySQL(shortlnk ShortLinks) {
	db, err := GetDBInstance()
	if err != nil {
		log.Print("DB error, fatal handler.")
		log.Fatal(err)
	} else {
		// Insert into the database
		insShtLnk, err := db.Prepare(`update shortlinks set count = ? where id = ?`)
		if err != nil {
			log.Fatal(err)
		}
		insShtLnk.Exec(shortlnk.Count, shortlnk.Id )
		log.Println("Successfully saved to MySQL database.")
	}
}

func hit_link_queue_receive() {
	conn, err := amqp.Dial("amqp://"+rabbitmq_user+":"+rabbitmq_pass+"@"+rabbitmq_server+":"+rabbitmq_port+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		hitShortLink_queue, // name
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
	if err!=nil {
		failOnError(err, "Failed to register a consumer")	
	}
	
	var goroutineDelta = make(chan int)
	go func() {
		for d := range msgs {
			var shortlnk ShortLinks
			err = json.Unmarshal(d.Body, &shortlnk)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			shortlnk.Count = shortlnk.Count+1
			updateShortLinkToMySQL(shortlnk)
			updateShortLinkToCache(shortlnk)
		}
	}()

	// Avoid the go routine from exiting
	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		if numGoroutines == 0 { os.Exit(0) }
	}
}
