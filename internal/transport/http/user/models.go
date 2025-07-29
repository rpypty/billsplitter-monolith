package user

type UpdateUserProfileReq struct {
	Username  string `json:"username" example:"johndoe"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}
