package base_interface

type IBaseUseCase[T any, Req any] interface {
	AddData(data Req) (T, error)
	FindDatas(page *int, limit *int) ([]T, error)
	FindDataById(id uint) (T, error)
	UpdateData(id uint, data Req) (T, error)
	DeleteData(id uint) (T, error)
}
