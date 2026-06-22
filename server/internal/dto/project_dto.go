package dto

import "time"

type CreateProjectReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type ProjectResp struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEnvReq struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type EnvResp struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddMemberReq struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin editor viewer"`
}

type UpdateMemberReq struct {
	Role string `json:"role" validate:"required,oneof=admin editor viewer"`
}

type MemberResp struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
