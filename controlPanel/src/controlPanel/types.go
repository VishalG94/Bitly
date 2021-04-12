package main

type ShortLinks struct {
	Id        string    `json: "id"`
	ShortLink string `json: "shortlink"`
	URL       string `json: "uri"`
	Count     int    `json: "hits"`
}

type Range struct {
	Min uint64
	Max uint64
}