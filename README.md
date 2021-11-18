
## Golang x Postgres exploratory

thanks to initial codebase from https://github.com/richyen/toolbox/blob/master/pg/orm/gorm_demo/gorm_demo.go 

# easy run :
go run main.go

# usage with httpie example :

- http :8080/cars
- http :8080/cars/1
- http --raw '{"destination":"garut","driverid":1,"userid":2,"price":100}' POST :8080/order