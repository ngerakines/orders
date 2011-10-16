
package main

import (
	"log"
)

type Event struct {
	Id, OrderId, Name, Value string
}

func storeEvent(id, orderId, name, value string) {
	stmt, err := db.Prepare("insert into event values (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, orderId, name, value); error != nil {
		log.Println(error)
		return
	}
	if error := stmt.Execute(); error != nil {
		log.Println(error)
		return
	}
}

func getEvents(id string) []Event {
	err := db.Query("select * from event where order_id = \"" + id + "\"")
	if err != nil {
		log.Println(err)
	    return make([]Event, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Event, 0)
	}
	events := make([]Event, result.RowCount())
	for index, row := range result.FetchRows() {
		var event Event
		event.Id = row[0].(string)
		event.OrderId = row[1].(string)
		event.Name = row[2].(string)
		event.Value = row[3].(string)
		events[index] = event
	}
	db.FreeResult()

	return events
}
