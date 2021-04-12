package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/catinello/base62"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"time"
	"github.com/rs/cors"
)

// RabbitMQ Config
var rabbitmq_server = "10.0.3.80"
var rabbitmq_port = "5672"
var createShortLink_queue = "shortlinks"
var rabbitmq_user = "guest"
var rabbitmq_pass = "guest"

var (
	username = "cmpe281"
	password = "cmpe281"
	hostname = "10.0.1.76"
	port 	 = "4306"
	dbname   = "bitly"
)

var (
	currCntrRange Range 	= Range{}
	currCount     uint64   	= 0
	gatewayIp           	= "http://34.221.223.237:8000"
	redirectPath        	= "/lr/"
	counterServerIp 		= "http://10.0.1.199:3002/basecount"
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
	mx.HandleFunc("/shortlink", gumballHandler(formatter)).Methods("GET")
	mx.HandleFunc("/shortlink", shortlinkCreateHandler(formatter)).Methods("POST")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Get the new range from counter server
func getCounterRange() error {
	w, err := http.Get(counterServerIp)
	if err != nil {
		failOnError(err, "Counter server connection failed")
	}
	decoder := json.NewDecoder(w.Body)
	
	err = decoder.Decode(&currCntrRange)
	if err != nil {
		log.Fatal(err)
	}
	currCount = currCntrRange.Min
	return err
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Control Panel API version 1.0 alive!"})
	}
}

// API Gumball Machine Handler
func gumballHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		db, err := GetDBInstance()
		if err != nil {
			log.Print("DB error, fatal handler.")
			log.Fatal(err)
		} else {

			var (
				id     		string
				shortlink  	string
				url  		string
				count 		int
			)
			rows, err := db.Query("select id, shortlink, url from shortlinks where id = ?", 1)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&id, &url, &shortlink, &count)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(id, shortlink, url)
			}
			err = rows.Err()
			if err != nil {
				log.Fatal(err)
			}
			formatter.JSON(w, http.StatusOK, rows)
		}
	}
}

// API Create New short Link
func shortlinkCreateHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		uuid := uuid.NewV4()
		decoder := json.NewDecoder(req.Body)
		var shortlnk ShortLinks
		var (
			id     		string
			shortlink  	string
			url  		string
			count 		int
		)
		err := decoder.Decode(&shortlnk)
		if err != nil {
			log.Fatal(err)
		}
		
		url = shortlnk.URL
		id = uuid.String()
		shortlnk.Id = id
		shortlnk.Count=0


		db, err := GetDBInstance()
		if err != nil {
			log.Print("DB error, fatal handler.")
			log.Fatal(err)
		} else {
			rows, err := db.Query("select * from shortlinks where url =  ?", shortlnk.URL)
			if err != nil 	{
				log.Print("DB error, Could not check records.")
				log.Fatal(err)	
			}
			// var count int
			var rowCount=0

			for rows.Next() {
		    	rowCount++
		    	err = rows.Scan(&id, &url, &shortlink, &count)
				if err != nil {
					log.Fatal(err)
				}
		    }   
		    if(rowCount>=1){
		    	formatter.JSON(w, http.StatusConflict, shortlink)
		    	return
		    }
		}
			

		// Hashing
		var uriHash string
		if url == "" {
			log.Fatal("URL not present in the request!")
		}
		// uriHash = hasher.Encode(10000)
		//fmt.Println("count: %d rangeMax: %d", count, currCntrRange.Max)
		if currCount >= currCntrRange.Max {
			err := getCounterRange()
			log.Println(currCount, "max: ", currCntrRange.Max)
			if err != nil {
				failOnError(err, "Error encountered in counter server")
			}
		}
		uriHash = base62.Encode(int(currCount))
		currCount++
		shortlink = gatewayIp + redirectPath + uriHash
		shortlnk.ShortLink = shortlink

		log.Println("ShotLink: %s, URL: %s, id: %d", shortlnk.ShortLink, shortlnk.URL, shortlnk.Id, shortlnk.Count)
		queue_send(shortlnk)
		formatter.JSON(w, http.StatusOK, shortlnk)
	}
}

func queue_send(message ShortLinks) {
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

/*
CREATE TABLE shortlinks (
id varchar(255) NOT NULL UNIQUE,
url varchar(255) NOT NULL UNIQUE,
shortlink varchar(255) NOT NULL UNIQUE,
count bigint(20) NOT NULL,
PRIMARY KEY (id)
) ;


insert into shortlinks ( id, url, shortlink) values ( 1, 'https://www.youtube.com/watch?v=S4ugBZmctKA', 'abckded' ) ;
*/