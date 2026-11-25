package base_interface

type IBaseRepository[T any, Req any] interface {
	Gets(page *int, limit *int) ([]T, error)
	GetById(id uint) (T, error)
	Create(data Req) (T, error)
	UpdateById(id uint, data Req) (T, error)
	DeleteById(id uint) (T, error)
}
