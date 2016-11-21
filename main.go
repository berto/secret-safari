package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/berto/secret-santa/db"
)

type match struct {
	Name   string `json:"name"`
	Animal string `json:"animal"`
	Error  string `json:"error"`
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "3000"
	}

	//preparing mux and server
	conn := fmt.Sprint(host, ":", port)
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir("./client")))
	router.HandleFunc("/pair", pairHandler)
	router.HandleFunc("/random", randomizeHandler)
	router.HandleFunc("/all", getAllHandler)

	//serving
	log.Printf("serving on %v", conn)
	log.Fatal(http.ListenAndServe(conn, router))
}

func pairHandler(response http.ResponseWriter, request *http.Request) {
	time.Sleep(1 * time.Second)
	pair := findPair(request.URL.Query()["name"][0], request.URL.Query()["animal"][0])
	json, err := json.Marshal(pair)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(json)
}

func findPair(name, animal string) match {
	ok := db.Validate(name, animal)
	if !ok {
		return match{Error: "Name and Animal do not match"}
	}
	matchName, err := db.GetMatch(name)
	if err != nil {
		return match{Error: "Name and Animal do not match"}
	}
	return match{Name: matchName, Animal: db.GetPair(matchName)}
}

func randomizeHandler(response http.ResponseWriter, request *http.Request) {
	err := db.RandomInsert()
	if err != nil {
		response.Write([]byte(fmt.Sprintf("Error: %+v", err)))
	}
	response.Write([]byte("Done!"))
}

func getAllHandler(response http.ResponseWriter, request *http.Request) {
	list, err := db.GetAll()
	if err != nil {
		response.Write([]byte(fmt.Sprintf("Error: %+v", err)))
	}
	response.Write([]byte(fmt.Sprintf("DB: %+v", list)))
}
