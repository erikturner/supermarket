package main

import (
	"bytes"
	"encoding/json"
	"superMarketAPI/bsapi"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetProduce(t *testing.T) {

	tests := map[string]struct {
		testID       string
		expectedResp int
	}{
		"Get Produce - Pass: happy path":					{"1", http.StatusOK},
	}

	for _, test := range tests {
		println("EXECUTING Test_GetProduce: ", test.testID)
		router := mux.NewRouter()
		router.HandleFunc("/produce/fetch", getProduceHandler)

		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "http://localhost:3000/produce/fetch", nil)
		router.ServeHTTP(w, r)
		resp := w.Result()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if resp.StatusCode != test.expectedResp {
			println("FAILED FOR test Test_GetProduce: ", test.testID)
			t.Errorf("expecting %v response code. got code [%d] description [%s].",
				test.expectedResp, resp.StatusCode, resp.Status)
			t.FailNow()
		}
		println("EXECUTED Test_GetProduce: ", test.testID)
	}
}

func Test_AddProduce(t *testing.T) {

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

	tests := map[string]struct {
		testID       string
		expectedResp int
	}{
		"Add Produce - Pass: happy path":									{"1", http.StatusOK},
		"Add Produce - Fail: error unmarshalling":							{"2", http.StatusInternalServerError},
		"Add Produce - Fail: error missing required payload data":			{"3", http.StatusBadRequest},
		"Add Produce - Fail: error produce code invalid format":			{"4", http.StatusBadRequest},
		"Add Produce - Fail: error produce code already exists":			{"5", http.StatusBadRequest},
	}

	for _, test := range tests {
		println("EXECUTING Test_AddProduce: ", test.testID)
		router := mux.NewRouter()
		router.HandleFunc("/produce/add", addProduceHandler)

		w := httptest.NewRecorder()
		
		var requests []bsapi.Produce
		var request bsapi.Produce
		
		request = bsapi.Produce{
			ProduceCode: stringRef("TR32-YUT7-93WE-290K"),
			Name:        stringRef("Grapes"),
			UnitPrice:   floatRef(1.00),
		}

		if test.testID == "3" {
			request = bsapi.Produce{
				ProduceCode: stringRef("TR32-YUT7-93WE-290K"),
				Name:        stringRef(""),
				UnitPrice:   floatRef(1.00),
			}
		}

		if test.testID == "4" {
			request = bsapi.Produce{
				ProduceCode: stringRef("TR32-YUT7-93WE-290"),
				Name:        stringRef("Grapes"),
				UnitPrice:   floatRef(1.00),
			}
		}

		if test.testID == "5" {
			request = bsapi.Produce{
				ProduceCode: stringRef("A12T-4GH7-QPL9-3N4M"),
				Name:        stringRef("Lettuce"),
				UnitPrice:   floatRef(3.46),
			}
		}

		requests = append(requests, request)

		payload, _ := json.Marshal(requests)
		if test.testID == "2" {
			payload = []byte(`{`)
		}
		body := bytes.NewReader(payload)

		r := httptest.NewRequest(http.MethodPost, "http://localhost:3000/produce/add", body)
		router.ServeHTTP(w, r)
		resp := w.Result()

		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if resp.StatusCode != test.expectedResp {
			println("FAILED FOR test Test_AddProduce: ", test.testID)
			t.Errorf("expecting %v response code. got code [%d] description [%s].",
				test.expectedResp, resp.StatusCode, resp.Status)
			t.FailNow()
		}
		println("EXECUTED Test_AddProduce: ", test.testID)
	}
}

func Test_RemoveProduce(t *testing.T) {

	tests := map[string]struct {
		testID       string
		expectedResp int
	}{
		"Remove Produce - Pass: happy path":									{"1", http.StatusOK},
		"Remove Produce - Fail: error missing url parameter":					{"2", http.StatusBadRequest},
	}

	for _, test := range tests {
		println("EXECUTING Test_RemoveProduce: ", test.testID)
		router := mux.NewRouter()

		if test.testID == "2" {
			router.HandleFunc("/produce/{'code'}/remove", removeProduceHandler)
		} else {
			router.HandleFunc("/produce/{code}/remove", removeProduceHandler)
		}

		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/produce/TR32-YUT7-93WE-290K/remove", nil)
		router.ServeHTTP(w, r)
		resp := w.Result()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if resp.StatusCode != test.expectedResp {
			println("FAILED FOR test Test_RemoveProduce: ", test.testID)
			t.Errorf("expecting %v response code. got code [%d] description [%s].",
				test.expectedResp, resp.StatusCode, resp.Status)
			t.FailNow()
		}
		println("EXECUTED Test_RemoveProduce: ", test.testID)
	}
}
