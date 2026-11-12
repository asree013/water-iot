package base_interface

type IBaseRepository[T any] interface {
	Gets(page *int, limit *int) ([]T, error)
	GetById(id uint) (T, error)
	Create(data T) (T, error)
	Update(id uint, data T) (T, error)
	Delete(id uint) (T, error)
}
