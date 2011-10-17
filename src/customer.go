
package main

import (
	"os"
	"log"
)

type Customer struct {
	Id, Name, Email string
	Total float64
}

type CustomerOrder struct {
	CustomerId, OrderId, OrderName string
}

type CustomerOrderLines struct {
	CustomerId, OrderId, LineId string
}

func (customer Customer) Orders() []CustomerOrder {
	orders := getCustomerOrders(customer.Id)
	customerOrders := make([]CustomerOrder, len(orders))
	for index, order := range orders {
		var customerOrder CustomerOrder
		customerOrder.CustomerId = customer.Id
		customerOrder.OrderId = order.Id
		customerOrder.OrderName = order.Name
		customerOrders[index] = customerOrder
	}
	return customerOrders
}

func (customerOrder CustomerOrder) Total() float64 {
	return getOrderTotalForCustomer(customerOrder.OrderId, customerOrder.CustomerId)
}

func (customerOrder CustomerOrder) Lines() []Line {
	return  getLinesByOrderAndCustomer(customerOrder.OrderId, customerOrder.CustomerId)
}

func storeCustomer(id, name, email string) {
	stmt, err := db.Prepare("insert into customer values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, name, email); error != nil {
		log.Println(error)
		return
	}
	if error := stmt.Execute(); error != nil {
		log.Println(error)
		return
	}
}

func getCustomer(id string) (customer *Customer, err os.Error) {
	dberr := db.Query("select * from customer where id = \"" + id + "\"")
	if dberr != nil {
	    return nil, dberr
	}
	result, dberr := db.StoreResult()
	if dberr != nil {
	    return nil, err
	}
	if count := result.RowCount(); count == 0 {
		db.FreeResult()
		return nil, os.NewError("Invalid customer id")
	}
	customer = new(Customer)
	for _, row := range result.FetchRows() {
		customer.Id = row[0].(string)
		customer.Name = row[1].(string)
		customer.Email = row[2].(string)
	}
	db.FreeResult()
	return customer, nil
}

func getCustomers() []Customer {
	err := db.Query("select * from customer")
	if err != nil {
		log.Println(err)
	    return make([]Customer, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Customer, 0)
	}
	customers := make([]Customer, result.RowCount())
	for index, row := range result.FetchRows() {
		var customer Customer
		customer.Id = row[0].(string)
		customer.Name = row[1].(string)
		customer.Email = row[2].(string)
		customers[index] = customer
	}
	db.FreeResult()
	return customers
}

func getCustomerSummary(id string) []Customer {
	err := db.Query("select customer.id, customer.name, customer.email, sum(line.quantity * line.price) from line, customer where customer.id = line.customer and order_id = \"" + id + "\" group by customer")
	if err != nil {
		log.Println(err)
	    return make([]Customer, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Customer, 0)
	}
	customers := make([]Customer, result.RowCount())
	for index, row := range result.FetchRows() {
		var customer Customer
		customer.Id = row[0].(string)
		customer.Name = row[1].(string)
		customer.Email = row[2].(string)
		customer.Total = row[3].(float64)
		customers[index] = customer
	}
	db.FreeResult()
	return customers
}
