
package main

import (
	"log"
)

type Line struct {
	Id, OrderId, CustomerId, Name string
	Quantity int64
	Price float64
}

func (line Line) Customer() *Customer {
	company, err := getCustomer(line.CustomerId)
	if err == nil {
		return company
	}
	return nil
}

func storeLine(id, orderId, customer, name string, quantity int, price float64) {
	stmt, err := db.Prepare("insert into line values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, orderId, customer, name, quantity, price); error != nil {
		log.Println(error)
		return
	}
	if error := stmt.Execute(); error != nil {
		log.Println(error)
		return
	}
}

func getLinesByOrder(id string) []Line {
	err := db.Query("select * from line where order_id = \"" + id + "\"")
	if err != nil {
		log.Println(err)
	    return make([]Line, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Line, 0)
	}
	lines := make([]Line, result.RowCount())
	for index, row := range result.FetchRows() {
		var line Line
		line.Id = row[0].(string)
		line.OrderId = row[1].(string)
		line.CustomerId = row[2].(string)
		line.Name = row[3].(string)
		line.Quantity = row[4].(int64)
		line.Price = row[5].(float64)
		lines[index] = line
	}
	db.FreeResult()
	return lines
}

func getOrderLineSummary(id string) []Line {
	err := db.Query("select name, count(*), sum(quantity * price) from line where order_id = \"" + id + "\" group by name order by name")
	if err != nil {
		log.Println(err)
	    return make([]Line, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Line, 0)
	}
	lines := make([]Line, result.RowCount())
	for index, row := range result.FetchRows() {
		var line Line
		line.Name = row[0].(string)
		line.Quantity = row[1].(int64)
		line.Price = row[2].(float64)
		lines[index] = line
	}
	db.FreeResult()
	return lines
}
