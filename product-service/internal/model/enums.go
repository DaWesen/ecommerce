package model

//ProductStatus
type ProductStatus int32

const (
	ProductStatusDRAFT   ProductStatus = 0 //预上架
	ProductStatusONLINE  ProductStatus = 1 //上架
	ProductStatusOFFLINE ProductStatus = 2 //下架
	ProductStatusDELETED ProductStatus = 3 //删除
)
