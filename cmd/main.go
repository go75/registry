package main

import (
	"encoding/json"
	"log"
	"registry/route"
	"registry/server"
	"registry/service"
)

func main() {
	srv := server.Default()

	registHandlers(srv)

	err := srv.Run()
	if err != nil {
		panic(err)
	}
}

func registHandlers(srv *server.RegistryService) {
	serviceInfos := srv.ServiceInfos

	srv.Regist(route.REMOVE, func(b []byte) {
		serviceInfo := new(service.ServiceInfo)
		err := json.Unmarshal(b, serviceInfo)
		if err != nil {
			log.Println("service info unmarshal err:", err)
			return
		}

		serviceInfos.Remove(serviceInfo)
	})

}