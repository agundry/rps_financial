package api

type ExpenseRequest struct {
	AustinThrow	string 	`json:"austin_throw"`
	SamThrow 	string 	`json:"sam_throw"`
	Cost		int 	`json:"cost"`
}
