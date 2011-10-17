package main

import (
	"os"
	"log"
	"time"
)

type Order struct {
	Id, Name, CompanyId string
	CreatedAt int64
}

func (order Order) Company() *Company {
	company, err := getCompany(order.CompanyId)
	if err == nil {
		return company
	}
	return nil
}

func (order Order) Events() []Event {
	return getEvents(order.Id)
}

func (order Order) Lines() []Line {
	return getLinesByOrder(order.Id)
}

func (order Order) CustomerSummary() []Customer {
	return getCustomerSummary(order.Id)
}

func (order Order) LineSummary() []Line {
	return getOrderLineSummary(order.Id)
}

func (order Order) Total() float64 {
	return getOrderTotal(order.Id)
}

func storeOrder(id, name, companyId string) {
	when := time.Seconds()
	stmt, err := db.Prepare("insert into orders values (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, name, companyId, when); error != nil {
		log.Println(error)
		return
	}
	if error := stmt.Execute(); error != nil {
		log.Println(error)
		return
	}
}

func getOrder(id string) (order *Order, err os.Error) {
	order = new(Order)
	dberr := db.Query("select * from orders where id = \"" + id + "\"")
	if dberr != nil {
	    return nil, dberr
	}
	result, dberr := db.UseResult()
	if dberr != nil {
	    return nil, err
	}
	for {
		row := result.FetchMap()
		if row == nil {
			break
		}
		order.Id = row["id"].(string)
		order.Name = row["name"].(string)
		order.CompanyId = row["company_id"].(string)
	}
	db.FreeResult()
	return order, nil
}

func listOrders() []Order {
	err := db.Query("select id, name, company_id from orders")
	if err != nil {
		log.Println(err)
	    return make([]Order, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Order, 0)
	}
	orders := make([]Order, result.RowCount())
	for index, row := range result.FetchRows() {
		var order Order
		order.Id = row[0].(string)
		order.Name = row[1].(string)
		order.CompanyId = row[2].(string)
		orders[index] = order
	}
	db.FreeResult()

	return orders
}

func getOrderTotal(id string) float64 {
	query := "select sum(quantity * price) from line where order_id = \"" + id + "\""
	log.Println(query)
	dberr := db.Query(query)
	if dberr != nil {
		log.Println(dberr)
	    return 0
	}
	result, dberr := db.StoreResult()
	if dberr != nil {
		log.Println(dberr)
	    return 0
	}
	if count := result.RowCount(); count == 1 {
		var value float64
		for _, row := range result.FetchRows() {
			value = row[0].(float64)
		}
		db.FreeResult()
		return value
	}
	db.FreeResult()
	return 0
}
