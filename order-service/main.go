package main

import (
	api "ecommerce/order-service/kitex_gen/api/orderservice"
	"log"
)

func main() {
	svr := api.NewServer(new(OrderServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
