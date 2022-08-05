package models

type Merchant struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
}

type Member struct {
	Id         int64  `json:"id"`
	MerchantId int64  `json:"merchantId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}
