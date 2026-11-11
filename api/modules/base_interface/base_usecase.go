package base_interface

type IBaseUseCase[T any] interface {
	AddData(data T) (T, error)
	FindDatas(page *int, limit *int) ([]T, error)
	FindDataById(id uint) (T, error)
	UpdateData(id uint, data T) (T, error)
	DeleteData(id uint) error
}
