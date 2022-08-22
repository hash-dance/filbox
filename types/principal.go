package types

// Principal api from principal
type Principal struct {
	BaseModel
	ExternalID string `json:"external_id"` // 外部系统ID,唯一识别号
	Username   string `json:"username"`    // 登录名称,也是单点登录系统的uuid
	Role       int    `json:"role"`
}

// UpdateRoleArg update user args
type UpdateRoleArg struct {
	Role int `json:"role" validate:"min=1,max=2,required"`
}

type CreatePrincipalOptions struct {
	UserName string `json:"user_name" validate:"required"`
	Role     *int   `json:"role" validate:"min=0,max=2,required"`
}
