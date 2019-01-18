package main

import (
	"github.com/coupons/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	CouponHandler := handlers.CouponHandler{}

	router.NewRoute().Path("/coupons").Methods(http.MethodPost).Handler(CouponHandler)

	log.Fatal(http.ListenAndServe(":6584", router))
}