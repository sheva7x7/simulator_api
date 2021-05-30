package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	ESClient "wenle/elasticsearch/esclient"
	"wenle/elasticsearch/vehicles"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func HandleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/vehicle", UpdateVehicle).Methods("POST")
	myRouter.HandleFunc("/boundary", UpdateBoundary).Methods("POST")
	myRouter.HandleFunc("/boundary", GetBoundary).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/vehicles/bound", GetVehiclesWithinBoundary).Methods("GET", "OPTIONS")
	// myRouter.HandleFunc("/vehicle/{id}", GetVehicle)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(myRouter)
	esclient, err := ESClient.GetESClient()
	if err != nil {
		fmt.Println("error", err)
	}
	if esclient != nil {
		fmt.Println("esclient started")
	}
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var v vehicles.Vehicle
	json.Unmarshal(reqBody, &v)
	res, err := vehicles.UpdateVehicle(v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad Request"}`))
	} else {
		resJson, _ := json.Marshal(res)
		fmt.Println(string(resJson))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(resJson)))
	}
}

func UpdateBoundary(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var boundary vehicles.Boundary
	json.Unmarshal(reqBody, &boundary)
	res, err := vehicles.UpdateBoundary(boundary)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad Request"}`))
	} else {
		resJson, _ := json.Marshal(res)
		fmt.Println(string(resJson))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(resJson)))
	}
}

func GetBoundary(w http.ResponseWriter, r *http.Request) {
	boundary := vehicles.GetBoundary()
	boundaryString, _ := json.Marshal(boundary)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(boundaryString)))
}

func GetVehiclesWithinBoundary(w http.ResponseWriter, r *http.Request) {
	vehicles, err := vehicles.GetVehiclesWithinBoundary()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad Request"}`))
	} else {
		resJson, _ := json.Marshal(vehicles)
		fmt.Println(string(resJson))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(resJson)))
	}
}
