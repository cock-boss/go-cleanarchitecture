package usecases

import (
	"fmt"
	"go-cleanarchitecture.bak/src/usecases"
	"go-cleanarchitecture/domain"
)

type UserRepository interface {
	Store(user User)
	FindById(id int) User
}

type User struct {
	Id       int
	IsAdmin  bool
	Customer domain.Customer
}

//type Item struct {
//	Id    int
//	Name  string
//	Value float64
//}

type Logger interface {
	Log(args ...interface{})
}

type OrderInteractor interface {
	Items(userId, orderId int) ([]usecases.Item, error)
	Add(userId, orderId, itemId int) error
}


type orderInteractor struct {
	UserRepository  UserRepository
	OrderRepository domain.OrderRepository
	ItemRepository  domain.ItemRepository
	Logger          Logger
}

func New(userRepo UserRepository, orderRepo domain.OrderRepository, itemRepo domain.ItemRepository, log Logger) OrderInteractor {
	it := new(orderInteractor)
	it.UserRepository = userRepo
	it.OrderRepository = orderRepo
	it.ItemRepository = itemRepo
	it.Logger = log
	return it
}

func (interactor *orderInteractor) Items(userId, orderId int) ([]usecases.Item, error) {
	var items []usecases.Item
	user := interactor.UserRepository.FindById(userId)
	order := interactor.OrderRepository.FindById(orderId)
	if user.Customer.Id != order.Customer.Id {
		message := "User #%d (customer #%d) "
		message += "is not allowed to see items "
		message += "in order #%d (of customer #%d)"
		err := fmt.Errorf(message,
			user.Id,
			user.Customer.Id,
			order.Id,
			order.Customer.Id)
		interactor.Logger.Log(err.Error())
		items = make([]usecases.Item, 0)
		return items, err
	}
	items = make([]usecases.Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = usecases.Item{item.Id, item.Name, item.Value}
	}
	return items, nil
}

func (interactor *orderInteractor) Add(userId, orderId, itemId int) error {
	var message string
	user := interactor.UserRepository.FindById(userId)
	order := interactor.OrderRepository.FindById(orderId)
	if user.Customer.Id != order.Customer.Id {
		message = "User #%d (customer #%d) "
		message += "is not allowed to add items "
		message += "to order #%d (of customer #%d)"
		err := fmt.Errorf(message,
			user.Id,
			user.Customer.Id,
			order.Id,
			order.Customer.Id)
		interactor.Logger.Log(err.Error())
		return err
	}
	item := interactor.ItemRepository.FindById(itemId)
	if domainErr := order.Add(item); domainErr != nil {
		message = "Could not add item #%d "
		message += "to order #%d (of customer #%d) "
		message += "as user #%d because a business "
		message += "rule was violated: '%s'"
		err := fmt.Errorf(message,
			item.Id,
			order.Id,
			order.Customer.Id,
			user.Id,
			domainErr.Error())
		interactor.Logger.Log(err.Error())
		return err
	}
	interactor.OrderRepository.Store(order)
	interactor.Logger.Log(fmt.Sprintf(
		"User added item '%s' (#%d) to order #%d",
		item.Name, item.Id, order.Id))
	return nil
}

type AdminOrderInteractor struct {
	orderInteractor
}

func (interactor *AdminOrderInteractor) Add(userId, orderId, itemId int) error {
	var message string
	user := interactor.UserRepository.FindById(userId)
	order := interactor.OrderRepository.FindById(orderId)
	if !user.IsAdmin {
		message = "User #%d (customer #%d) "
		message += "is not allowed to add items "
		message += "to order #%d (of customer #%d), "
		message += "because he is not an administrator"
		err := fmt.Errorf(message,
			user.Id,
			user.Customer.Id,
			order.Id,
			order.Customer.Id)
		interactor.Logger.Log(err.Error())
		return err
	}
	item := interactor.ItemRepository.FindById(itemId)
	if domainErr := order.Add(item); domainErr != nil {
		message = "Could not add item #%d "
		message += "to order #%d (of customer #%d) "
		message += "as user #%d because a business "
		message += "rule was violated: '%s'"
		err := fmt.Errorf(message,
			item.Id,
			order.Id,
			order.Customer.Id,
			user.Id,
			domainErr.Error())
		interactor.Logger.Log(err.Error())
		return err
	}
	interactor.OrderRepository.Store(order)
	interactor.Logger.Log(fmt.Sprintf(
		"Admin added item '%s' (#%d) to order #%d",
		item.Name, item.Id, order.Id))
	return nil
}
