package main // import "github.com/madeleinesmith/coupons"

import (
	"github.com/madeleinesmith/coupons/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/madeleinesmith/coupons/dbservices"
	"github.com/madeleinesmith/coupons/handlers"
	"github.com/madeleinesmith/coupons/model/coupon"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	db := initializeDb()

	couponService := dbservices.CouponService{
		DB: db,
	}
	couponSerializer := coupon.Serializer{}
	couponHandler := handlers.CouponHandler{
		CouponService: couponService,
		Serializer: couponSerializer,
	}

	couponDetailsHandler := handlers.CouponDetailsHandler{
		CouponService: couponService,
		Serializer: couponSerializer,
	}

	router.NewRoute().Path("/coupons").Handler(couponHandler)
	router.NewRoute().Path("/coupon/{couponId}").Handler(couponDetailsHandler)

	log.Fatal(http.ListenAndServe(":6584", router))
}

func initializeDb() *sql.DB {
	applicationConfiguration := loadConfiguration()

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		applicationConfiguration.Database.User,
		applicationConfiguration.Database.Password,
		applicationConfiguration.Database.DBName)

	db, _ := sql.Open("postgres", connectionString)

	return db
}


func loadConfiguration() model.Config {
	var config model.Config
	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}