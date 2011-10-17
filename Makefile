include $(GOROOT)/src/Make.inc

TARG=orders
GOFILES=\
	src/main.go\
	src/mustache.go\
	src/uuid4.go\
	src/db.go\
	src/utils.go\
	src/model.go\
	src/company.go\
	src/order.go\
	src/event.go\
	src/line.go\
	src/customer.go

include $(GOROOT)/src/Make.cmd