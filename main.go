package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// BidData : bid-data object structure
type BidData struct {
	ID          string `json:"id"`
	PlacementID string `json:"placementID"`
	BidPrice    int    `json:"bidPrice"`
	Currency    string `json:"currency"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/call-bidders", CallBidders)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// CallBidders : this is the handler for calling bidders.
// makes multipe calls to the bidders at the same time.
func CallBidders(w http.ResponseWriter, r *http.Request) {
	var AllBids []BidData
	ch := make(chan string)
	for i := 1; i <= 10; i++ {
		go MakeRequest("http://localhost:5000/make-bid?placement-id=123", ch)
	}
	// fmt.Println(reflect.TypeOf(<-ch))
	Bid := &BidData{}
	for i := 1; i <= 10; i++ {
		b := <-ch
		println(b)
		err := json.Unmarshal([]byte(b), Bid)
		if err != nil {
			panic(err)
		}
		AllBids = append(AllBids, *Bid)
	}

	fmt.Fprintln(w, AllBids)
}

// MakeRequest : function to actually call bidding service.
func MakeRequest(url string, ch chan string) {
	timeOver := time.Duration(200 * time.Millisecond)
	client := http.Client{
		Timeout: timeOver,
	}
	resp, err := client.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- string(body)
}
