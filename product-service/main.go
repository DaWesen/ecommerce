package main

import (
	api "ecommerce/product-service/kitex_gen/api/productservice"
	"log"
)

func main() {
	svr := api.NewServer(new(ProductServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
