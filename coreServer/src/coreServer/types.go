package main

type gumballMachine struct {
	Id             	int
	CountGumballs   int
	ModelNumber 	string
	SerialNumber 	string
}

type order struct {
	Id             	string 	
	OrderStatus 	string	
}

var orders map[string] order

var ShortLinksMap map[string] order


//type ShortLinks struct {
//	Id             	string  `json:"id,omitempty"`
//	ShortLink   	string  `json:"id,omitempty"`
//	URL			 	string  `json:"id,omitempty"`
//}
type ShortLinks struct {
	Id        string    `json: "id"`
	ShortLink string `json: "shortlink"`
	URL       string `json: "uri"`
	Count     int    `json: "hits"`
}