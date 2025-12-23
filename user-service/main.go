package main

import (
	api "ecommerce/user-service/kitex_gen/api/userservice"
	"log"
)

func main() {
	svr := api.NewServer(new(UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
