
package main

import (
	"log"
)

type Line struct {
	Id, OrderId, Customer, Name string
	Quantity int64
	Price float64
}

func storeLine(id, orderId, customer, name string, quantity int, price float64) {
	log.Println(quantity)
	log.Println(price)
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
		line.Customer = row[2].(string)
		line.Name = row[3].(string)
		line.Quantity = row[4].(int64)
		line.Price = row[5].(float64)
		lines[index] = line
	}
	db.FreeResult()

	return lines
}
