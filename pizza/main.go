package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Pizza struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Pizzas []Pizza

func (ps Pizzas) FindById(id int) (Pizza, error) {
	for _, pizza := range ps {
		if pizza.Id == id {
			return pizza, nil
		}
	}

	return Pizza{}, fmt.Errorf("Couldn't find pizza with ID: %d", id)
}

type Order struct {
	PizzaId  int `json:"pizza_id"`
	Quantity int `json:"quantity"`
	Total    int `json:"total"`
}

type Orders []Order

type pizzasHandler struct {
	pizzas *Pizzas
}

func (ph pizzasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if len(*ph.pizzas) == 0 {
			http.Error(w, "Error: No pizzas found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(ph.pizzas)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type ordersHandler struct {
	pizzas *Pizzas
	orders *Orders
}

func (oh ordersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var o Order

		if len(*oh.pizzas) == 0 {
			http.Error(w, "Error: No pizzas found", http.StatusNotFound)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		p, err := oh.pizzas.FindById(o.PizzaId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
			return
		}

		o.Total = p.Price * o.Quantity
		*oh.orders = append(*oh.orders, o)
		json.NewEncoder(w).Encode(o)
	case http.MethodGet:
		json.NewEncoder(w).Encode(oh.orders)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	var orders Orders

	pizzas := Pizzas{
		Pizza{
			Id:    1,
			Name:  "Pepperoni",
			Price: 12,
		},
		Pizza{
			Id:    2,
			Name:  "Capricciosa",
			Price: 11,
		},
		Pizza{
			Id:    3,
			Name:  "Margherita",
			Price: 10,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/pizzas", pizzasHandler{&pizzas})
	mux.Handle("/orders", ordersHandler{&pizzas, &orders})

	log.Fatal(http.ListenAndServe(":8080", mux))
}