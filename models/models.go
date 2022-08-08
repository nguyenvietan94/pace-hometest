package models

// contains info of a merchant
type Merchant struct {
	MerchantID int64  `json:"merchantID"`
	Name       string `json:"name"`
	Age        int64  `json:"age"`
	Location   string `json:"location"`
}

// contains info of a member
type Member struct {
	MemberID   int64  `json:"memberID"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	MerchantID int64  `json:"merchantID"`
}
