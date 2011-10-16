package main

import (
	"os"
	"log"
)

type Order struct {
	Id, Name, CompanyId string
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

func storeOrder(id, name, companyId string) {
	stmt, err := db.Prepare("insert into orders values (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, name, companyId); error != nil {
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
