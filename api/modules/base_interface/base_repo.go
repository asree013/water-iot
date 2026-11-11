package base_interface

type IBaseRepository[T any] interface {
	Gets(page *int, limit *int) ([]T, error)
	GetById(id int) (T, error)
	Create(data T) (T, error)
	Update(data T) (T, error)
	Delete(id int) (T, error)
}
