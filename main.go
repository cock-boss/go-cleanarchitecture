package main

import (
	"go-cleanarchitecture/infrastructure"
	"go-cleanarchitecture/interfaces"
	"go-cleanarchitecture/usecases"
	"net/http"
)

func main() {
	dbHandler := infrastructure.NewSqliteHandler("/var/tmp/production.sqlite")

	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbUserRepo"] = dbHandler
	handlers["DbCustomerRepo"] = dbHandler
	handlers["DbItemRepo"] = dbHandler
	handlers["DbOrderRepo"] = dbHandler


	userRepo := interfaces.NewDbUserRepo(handlers)
	itemRepo := interfaces.NewDbItemRepo(handlers)
	orderRepo := interfaces.NewDbOrderRepo(handlers)
	logger := new(infrastructure.Logger)
	orderInteractor := usecases.New(userRepo, orderRepo, itemRepo, logger)

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.OrderInteractor = orderInteractor

	http.HandleFunc("/orders", func(res http.ResponseWriter, req *http.Request) {
		webserviceHandler.ShowOrder(res, req)
	})
	http.ListenAndServe(":8080", nil)
}
