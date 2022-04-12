package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
)

var db *gorm.DB
var err error

// Product is a representation of a product
type Product struct {
	ID    int             `form:"id" json:"id"`
	Code  string          `form:"code" json:"code"`
	Name  string          `form:"name" json:"name"`
	Price decimal.Decimal `form:"price" json:"price" sql:"type:decimal(16,2);"`
}

// Result is an array of product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// T_CNOTE_SINGLE is a representation of table cnote single
type T_CNOTE_SINGLE struct {
	CityName          string `json:"city_name"`
	CnoteAmount       string `json:"cnote_amount"`
	CnoteCustNo       string `json:"cnote_cust_no"`
	CnoteDate         string `json:"cnote_date"`
	CnoteDestination  string `json:"cnote_destination"`
	CnoteGoodsDescr   string `json:"cnote_goods_descr"`
	CnoteNo           string `json:"cnote_no"`
	CnoteOrigin       string `json:"cnote_origin"`
	CnotePodDate      string `json:"cnote_pod_date"`
	CnotePodReceiver  string `json:"cnote_pod_receiver"`
	CnoteReceiverName string `json:"cnote_receiver_name"`
	CnoteServicesCode string `json:"cnote_services_code"`
	CnoteWeight       string `json:"cnote_weight"`
	CustType          string `json:"cust_type"`
	EstimateDelivery  string `json:"estimate_delivery"`
	FreightCharge     string `json:"freight_charge"`
	Insuranceamount   string `json:"insuranceamount"`
	Keterangan        string `json:"keterangan"`
	LastStatus        string `json:"last_status"`
	Lat               string `json:"lat"`
	Long              string `json:"long"`
	Photo             string `json:"photo"`
	PodCode           string `json:"pod_code"`
	PodStatus         string `json:"pod_status"`
	Priceperkg        string `json:"priceperkg"`
	ReferenceNumber   string `json:"reference_number"`
	Servicetype       string `json:"servicetype"`
	Shippingcost      string `json:"shippingcost"`
	Signature         string `json:"signature"`
}

// Respond of cnote
type ResultCnote struct {
	Cnote interface{} `json:"cnote"`
}

// Main
func main() {
	db, err = gorm.Open("mysql", "root:root123@/go_rest_api_crud?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&T_CNOTE_SINGLE{})
	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/cnotes", createCnote).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/cnotes/{cnote_no}", getCnote).Methods("GET")
	//myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	//myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func createCnote(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var cnote T_CNOTE_SINGLE
	json.Unmarshal(payloads, &cnote)

	db.Create(&cnote)

	res := ResultCnote{Cnote: cnote}
	resultcnote, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultcnote)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)

	res := Result{Code: 200, Data: product, Message: "Success get product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getCnote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	CnoteNo := vars["cnote"]

	var cnote T_CNOTE_SINGLE

	db.First(&cnote, CnoteNo)

	res := ResultCnote{Cnote: cnote}
	resultcnote, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultcnote)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdates Product
	json.Unmarshal(payloads, &productUpdates)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	res := Result{Code: 200, Data: product, Message: "Success update product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

////func deleteProduct(w http.ResponseWriter, r *http.Request) {
////	vars := mux.Vars(r)
////	productID := vars["id"]
////
////	var product Product
////
////	db.First(&product, productID)
////	db.Delete(&product)
////
////	res := Result{Code: 200, Message: "Success delete product"}
////	result, err := json.Marshal(res)
////
////	if err != nil {
////		http.Error(w, err.Error(), http.StatusInternalServerError)
////	}
////
////	w.Header().Set("Content-Type", "application/json")
////	w.WriteHeader(http.StatusOK)
////	w.Write(result)
//}
