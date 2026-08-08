package domain

type SendInvitesRequest struct {
	SendToAll      bool   `json:"send_to_all"`
	DepartmentName string `json:"department_name"`
}