package dto

// CreateRuleReq is the request payload for creating a rule.
type CreateRuleReq struct {
	NodeID     *uint  `json:"node_id"`
	TunnelID   *uint  `json:"tunnel_id"`
	Name       string `json:"name" binding:"required,min=1,max=100"`
	Type       string `json:"type" binding:"required,oneof=forward tunnel local_forward"`
	Protocol   string `json:"protocol" binding:"required,oneof=tcp udp tcp+udp"`
	ListenPort int    `json:"listen_port" binding:"required,min=1,max=65535"`

	Targets   []string `json:"targets"`
	Strategy  string   `json:"strategy" binding:"omitempty,oneof=round rand fifo hash"`
	EnableTLS bool     `json:"enable_tls"`
	Remark    string   `json:"remark"`
}

// UpdateRuleReq is the request payload for updating a rule.
type UpdateRuleReq struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	Protocol   string `json:"protocol" binding:"required,oneof=tcp udp tcp+udp"`
	ListenPort int    `json:"listen_port" binding:"required,min=1,max=65535"`

	Targets   []string `json:"targets"`
	Strategy  string   `json:"strategy" binding:"omitempty,oneof=round rand fifo hash"`
	EnableTLS bool     `json:"enable_tls"`
	Remark    string   `json:"remark"`
}

// RuleListReq is the query payload for listing rules.
type RuleListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=1000"`
	NodeID   uint   `form:"node_id"`
	TunnelID uint   `form:"tunnel_id"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
}

// SetDefaults applies default pagination values.
func (r *RuleListReq) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
}
