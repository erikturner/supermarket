package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"test.com/supermarket/bsapi/bsapi"
)

func stringRef(s string) *string {
	return &s
}

func floatRef(s float64) *float64 {
	return &s
}

var marketInv []bsapi.Produce

// Produce handler endpoints
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/produce/fetch", getProduceHandler).Methods(http.MethodGet)

	router.HandleFunc("/produce/add", addProduceHandler).Methods(http.MethodPost)

	router.HandleFunc("/produce/{code}/remove", removeProduceHandler).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// getProduceHandler gets all produce items
func getProduceHandler(w http.ResponseWriter, r *http.Request) {

	/*if len(marketInv) < 1 {
		log.Println("no rows found")
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}*/

	b, err := json.Marshal(marketInv)
	if err != nil {
		log.Println("error marshalling")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// addProduceHandler adds one or more produce items
func addProduceHandler(w http.ResponseWriter, r *http.Request) {

	inventory := []bsapi.Produce{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading body")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &inventory)
	if err != nil {
		log.Println("error unmarshalling")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, produce := range inventory {

		/* -----------------------------------------
		| check for required fields
		-------------------------------------------*/
		if produce.ProduceCode == nil || *produce.ProduceCode == "" || produce.Name == nil ||
			*produce.Name == "" || produce.UnitPrice == nil {
			log.Println("error missing required data in payload (produce code, name, unit price)")
			http.Error(w, "error missing required data in payload (produce code, name, unit price)", http.StatusBadRequest)
			return
		}

		/* -----------------------------------------
		| check for valid format of produce code
		-------------------------------------------*/
		match, _ := regexp.MatchString("[^A-Za-z0-9]", *produce.ProduceCode)
		str := strings.Split(*produce.ProduceCode, "-")
		for i := range str {
			size := []rune(str[i])
			if !match  || len(str) != 4 || (match && len(size) != 4) || len(size) != 4 {
				log.Println("error produce code invalid format")
				http.Error(w, "error produce code invalid format", http.StatusBadRequest)
				return
			}
		}

		/* -----------------------------------------
		| check if produce code already exists
		-------------------------------------------*/
		for _, market := range marketInv {
			if strings.ToUpper(*produce.ProduceCode) == strings.ToUpper(*market.ProduceCode) {
				log.Println("error produce code already exists")
				http.Error(w, "error produce code already exists", http.StatusBadRequest)
				return
			}
		}

		/* -----------------------------------------
		| add produce
		-------------------------------------------*/
		*produce.ProduceCode = strings.ToUpper(*produce.ProduceCode)
		marketInv = append(marketInv, produce)
	}

	b, err := json.Marshal(marketInv)
	if err != nil {
		log.Println("error marshalling")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// removeProduceHandler removes a produce item based on produce code
func removeProduceHandler(w http.ResponseWriter, r *http.Request) {

	/* -----------------------------------------
	| check url parameter exists
	-------------------------------------------*/
	code, ok := mux.Vars(r)["code"]
	if !ok {
		log.Println("error missing url parameter")
		http.Error(w, "error missing url parameter", http.StatusBadRequest)
		return
	}

	/* -----------------------------------------
	| remove produce item
	-------------------------------------------*/
	for j, inv := range marketInv {
		if strings.ToUpper(*inv.ProduceCode) == strings.ToUpper(code) {
			marketInv = append(marketInv[:j])
		}
	}

	b, err := json.Marshal(marketInv)
	if err != nil {
		log.Println("error marshalling")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func main () {

	produce1 := bsapi.Produce{
		ProduceCode: stringRef("A12T-4GH7-QPL9-3N4M"),
		Name:        stringRef("Lettuce"),
		UnitPrice:   floatRef(3.46),
	}

	produce2 := bsapi.Produce{
		ProduceCode: stringRef("E5T6-9UI3-TH15-QR88"),
		Name:        stringRef("Peach"),
		UnitPrice:   floatRef(2.99),
	}

	produce3 := bsapi.Produce{
		ProduceCode: stringRef("YRT6-72AS-K736-L4AR"),
		Name:        stringRef("Green Pepper"),
		UnitPrice:   floatRef(0.79),
	}

	produce4 := bsapi.Produce{
		ProduceCode: stringRef("TQ4C-VV6T-75ZX-1RMR"),
		Name:        stringRef("Gala Apple"),
		UnitPrice:   floatRef(3.59),
	}
	marketInv = append(marketInv, produce1, produce2, produce3, produce4)

	handleRequests()
}