package main

import (
	"os"
	"io"
	"github.com/Philio/GoMySQL"
	"github.com/garyburd/twister/server"
	"github.com/garyburd/twister/web"
	"log"
	"strconv"
)

var (
	db *mysql.Client
	db_err os.Error
)

func indexHandler(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")

	orders := listOrders()
	params := make(map[string]interface{})
	params["Orders"] = orders
	if len(orders) > 0 {
		params["LastOrder"] = orders[0]
	}

	io.WriteString(w, RenderFile("templates/index.html", params))
}

// Order

func createOrderFormHandler(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")

	companies := listCompanies()
	params := make(map[string]interface{})
	params["Companies"] = companies

	io.WriteString(w, RenderFile("templates/order-create.html", params))
}

func createOrderHandler(req *web.Request) {
	name := req.Param.Get("name")
	companyId := req.Param.Get("company")
	id := NewUUID()
	storeOrder(id.String(), name, companyId)
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	io.WriteString(w, "ok")
}

func viewOrderHandler(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")

	id := req.Param.Get("id")
	order, err := getOrder(id)
	if err != nil {
		io.WriteString(w, "ERROR: " + err.String())
		return
	}
	log.Println(order)
	params := make(map[string]interface{})
	params["Order"] = order
	customers := getCustomers()
	log.Println(customers)
	params["Customers"] = customers
	io.WriteString(w, RenderFile("templates/order-view.html", params))
}

// Company

func createCompanyFormHandler(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	io.WriteString(w, RenderFile("templates/company-create.html", map[string]string{"c":"world"}))
}

func createCompanyHandler(req *web.Request) {
	name := req.Param.Get("name")
	address := req.Param.Get("address")
	phone := req.Param.Get("phone")
	url := req.Param.Get("url")
	id := NewUUID()
	storeCompany(id.String(), name, address, phone, url)
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	io.WriteString(w, "ok")
}

func viewCompanyHandler(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")

	id := req.Param.Get("id")
	company, err := getCompany(id)
	if err != nil {
		io.WriteString(w, "ERROR: " + err.String())
		return
	}
	params := make(map[string]interface{})
	params["Company"] = company
	io.WriteString(w, RenderFile("templates/company-view.html", params))
}

// Event

func createEventHandler(req *web.Request) {
	orderId := req.Param.Get("order")
	name := req.Param.Get("name")
	value := req.Param.Get("value")
	id := NewUUID()
	storeEvent(id.String(), orderId, name, value)
	req.Redirect("/order/?id=" + orderId, false)
}

// Customer

func createCustomerHandler(req *web.Request) {
	orderId := req.Param.Get("order")
	name := req.Param.Get("name")
	email := req.Param.Get("email")
	id := NewUUID()
	storeCustomer(id.String(), name, email)
	req.Redirect("/order/?id=" + orderId, false)
}

// Line

func createLineHandler(req *web.Request) {
	orderId := req.Param.Get("order")
	customer := req.Param.Get("customer")
	name := req.Param.Get("name")
	quantity := req.Param.Get("quantity")
	price := req.Param.Get("price")
	id := NewUUID()
	quantityInt, _ := strconv.Atoi(quantity)
	priceInt, _ := strconv.Atof64(price)
	storeLine(id.String(), orderId, customer, name, quantityInt, priceInt)
	req.Redirect("/order/?id=" + orderId, false)
}

// main loop

func main() {
	db, db_err = mysql.DialTCP("localhost", "root", "asd123", "orders")
	if db_err != nil {
		log.Println(db_err)
	    os.Exit(1)
	}

	port := ":8080"
	if envPort := os.Getenv("ORDERS_PORT"); envPort != "" {
		port = envPort
	}

	h := web.FormHandler(10000, false,
		web.NewRouter().
			Register("/", "GET", indexHandler).
			Register("/company/create", "GET", createCompanyFormHandler, "POST", createCompanyHandler).
			Register("/company/", "GET", viewCompanyHandler).
			Register("/order/create", "GET", createOrderFormHandler, "POST", createOrderHandler).
			Register("/order/", "GET", viewOrderHandler).
			Register("/event/create", "POST", createEventHandler).
			Register("/customer/create", "POST", createCustomerHandler).
			Register("/line/create", "POST", createLineHandler).
			Register("/static/<path:.*>", "GET", web.DirectoryHandler("./static/", new(web.ServeFileOptions))))
	server.Run(port, h)
}
