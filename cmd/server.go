package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Products struct {
	ProductList []Product
	OrderList   []Order
}

type Product struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type Order struct {
	Products  []Product `json:"products"`
	SKU       string    `json:"sku"`
	Qty       int       `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Status   int       `json:"status"`
	Message  string    `json:"message"`
	Products []Product `json:"products,omitempty"`
	Product  Product   `json:"product,omitempty"`
	Orders   []Order   `json:"orders,omitempty"`
}

func serveHTTP(products *Products) {
	router := mux.NewRouter()
	router.HandleFunc("/create_product", products.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/get_product_list", products.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/create_order", products.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/get_order_list", products.GetOrderList).Methods(http.MethodGet)

	// Start the server and pass the router
	if err := http.ListenAndServe(":8000", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func (pd *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println("Got error when decoding: ", err)
		return
	}

	pd.ProductList = append(pd.ProductList, product)
	response := Response{
		Message: "Success",
		Status:  http.StatusOK,
		Product: product,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error when writing response: ", err)
	}
}

func (pd *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  200,
		Message: "Success",
	}
	response.Products = make([]Product, len(pd.ProductList))
	copy(response.Products, pd.ProductList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error when writing response: ", err)
	}
}

func (pd *Products) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Println("Got error when decoding: ", err)
		return
	}

	response := Response{
		Status:  200,
		Message: "Success",
	}

	for idx, product := range pd.ProductList {
		if product.SKU == order.SKU {
			if product.Stock < order.Qty {
				response.Message = "Order cannot be created. Quantity is less than available stock."
				break
			} else {
				orderResponse := Order{
					Products:  []Product{product},
					Qty:       order.Qty,
					CreatedAt: time.Now(),
				}
				pd.ProductList[idx].Stock -= order.Qty
				pd.OrderList = append(pd.OrderList, orderResponse)
				response.Message = "Order created."
				response.Orders = append(response.Orders, orderResponse)
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error when writing response: ", err)
	}
}

func (pd *Products) GetOrderList(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  200,
		Message: "Success",
	}
	response.Orders = make([]Order, len(pd.OrderList))

	copy(response.Orders, pd.OrderList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error when writing response: ", err)
	}
}
