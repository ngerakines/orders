
package main

import (
	"os"
	"log"
	"strings"
)

type Company struct {
	Id, Name, Address, Phone, Url string
}

func (company Company) AddressLines() []string {
	return strings.Split(company.Address, "\n")
}

func storeCompany(id, name, address, phone, url string) {
	stmt, err := db.Prepare("insert into company values (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return
	}
	if error := stmt.BindParams(id, name, address, phone, url); error != nil {
		log.Println(error)
		return
	}
	if error := stmt.Execute(); error != nil {
		log.Println(error)
		return
	}
}

func getCompany(id string) (company *Company, err os.Error) {
	company = new(Company)
	dberr := db.Query("select * from company where id = \"" + id + "\"")
	if dberr != nil {
	    return nil, dberr
	}
	result, dberr := db.UseResult()
	if dberr != nil {
	    return company, err
	}
	for {
		row := result.FetchMap()
		if row == nil {
			break
		}
		log.Println(row)
		company.Id = row["id"].(string)
		company.Name = row["name"].(string)
		company.Address = string([]uint8( row["address"].([]uint8)  ))
		company.Phone = row["phone"].(string)
		company.Url = row["url"].(string)
	}
	db.FreeResult()
	return company, nil
}

func listCompanies() []Company {
	err := db.Query("select id, name, address, phone, url from company")
	if err != nil {
		log.Println(err)
	    return make([]Company, 0)
	}
	result, err := db.StoreResult()
	if err != nil {
		log.Println(err)
	    return make([]Company, 0)
	}
	companies := make([]Company, result.RowCount())
	for index, row := range result.FetchRows() {
		var company Company
		company.Id = row[0].(string)
		company.Name = row[1].(string)
		company.Address = string([]uint8( row[2].([]uint8)  ))
		company.Phone = row[3].(string)
		company.Url = row[4].(string)
		companies[index] = company
	}
	db.FreeResult()

	return companies
}
