package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

type Driver struct {
	gorm.Model
	Name    string
	License string
	Cars    []Car
}

type Car struct {
	gorm.Model
	Year      int
	Make      string
	ModelName string
	DriverID  int
}

type User struct {
	gorm.Model
	Name string
	Age  int
}

type Order struct {
	gorm.Model
	DriverID    int
	UserID      int
	Destination string
	Price       int
	Finished    bool
	Rating      int
}

var db *gorm.DB
var err error

var (
	drivers = []Driver{
		{Name: "Jimmy Johnson", License: "ABC"},
		{Name: "Howard Hills", License: "XYZ789"},
		{Name: "Craig Colbin", License: "DEF333"},
	}

	cars = []Car{
		{Year: 2000, Make: "Toyota", ModelName: "Tundra", DriverID: 1},
		{Year: 2001, Make: "Honda", ModelName: "Accord", DriverID: 1},
		{Year: 2002, Make: "Nissan", ModelName: "Sentra", DriverID: 2},
		{Year: 2003, Make: "Ford", ModelName: "F-150", DriverID: 3},
	}

	users = []User{
		{Name: "Ardhi", Age: 25},
		{Name: "Kim Go Eun", Age: 30},
	}
)

func main() {
	router := mux.NewRouter()

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=ardhihdra dbname=postgres sslmode=disable password=postgres")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	defer db.Close()
	InitDB()

	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")
	router.HandleFunc("/order", NewOrder).Methods("POST")

	handler := cors.Default().Handler(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Up and Running...")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))

}

func InitDB() {
	db.AutoMigrate(&Driver{})
	db.AutoMigrate(&Car{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})
	ClearDB()

	for index := range cars {
		db.Create(&cars[index])
	}
	for index := range drivers {
		db.Create(&drivers[index])
	}
	for index := range users {
		db.Create(&users[index])
	}
}

func ClearDB() {
	var cars []Car
	db.Find(&cars)
	db.Delete(&cars)

	var drivers []Driver
	db.Find(&drivers)
	db.Delete(&drivers)

	var users []User
	db.Find(&users)
	db.Delete(&users)

	var orders []Order
	db.Find(&orders)
	db.Delete(&orders)
}

func GetCars(res http.ResponseWriter, req *http.Request) {
	var cars []Car
	db.Find(&cars)
	json.NewEncoder(res).Encode(&cars)
}

func GetCar(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var car Car
	db.First(&cars, params["id"])
	json.NewEncoder(res).Encode(&car)
}

func GetDriver(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var driver Driver
	var cars []Car
	db.First(&driver, params["id"])
	db.Model(&driver).Related(&cars)
	driver.Cars = cars
	json.NewEncoder(res).Encode(&driver)
}

func DeleteCar(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var car Car
	db.First(&car, params["id"])
	db.Delete(&car)

	GetCars(res, req)
}

func NewOrder(res http.ResponseWriter, req *http.Request) {
	var order Order
	err := json.NewDecoder(req.Body).Decode(&order)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	db.Create(&order)
	fmt.Fprintf(res, "New Order : %+v", order)
}
