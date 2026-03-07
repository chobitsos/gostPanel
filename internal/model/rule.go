package model

import (
	"time"

	"gorm.io/gorm"
)

// RuleStatus is the runtime status of a rule.
type RuleStatus string

const (
	RuleStatusRunning RuleStatus = "running"
	RuleStatusStopped RuleStatus = "stopped"
	RuleStatusError   RuleStatus = "error"
)

// RuleProtocol is the protocol for a rule.
type RuleProtocol string

const (
	RuleProtocolTCP    RuleProtocol = "tcp"
	RuleProtocolUDP    RuleProtocol = "udp"
	RuleProtocolTCPUDP RuleProtocol = "tcp+udp"
)

// RuleType is the forwarding mode of a rule.
type RuleType string

const (
	RuleTypeForward      RuleType = "forward"
	RuleTypeTunnel       RuleType = "tunnel"
	RuleTypeLocalForward RuleType = "local_forward"
)

// GostRule is the forwarding rule model.
type GostRule struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	NodeID     *uint        `gorm:"index" json:"node_id"`
	Name       string       `gorm:"size:100;not null" json:"name"`
	Type       RuleType     `gorm:"size:20;not null;default:forward" json:"type"`
	TunnelID   *uint        `gorm:"index" json:"tunnel_id"`
	Protocol   RuleProtocol `gorm:"size:10;not null;default:tcp" json:"protocol"`
	ListenPort int          `gorm:"not null" json:"listen_port"`

	Targets   []string   `gorm:"type:json;serializer:json" json:"targets"`
	Strategy  string     `gorm:"size:20;default:round" json:"strategy"`
	EnableTLS bool       `gorm:"default:false" json:"enable_tls"`
	Status    RuleStatus `gorm:"size:20;default:stopped" json:"status"`
	ServiceID string     `gorm:"size:100" json:"service_id"`
	AutoResume bool      `gorm:"default:false" json:"-"`

	ObserverID string `gorm:"size:100" json:"observer_id"`

	InputBytes    int64 `gorm:"default:0" json:"input_bytes"`
	OutputBytes   int64 `gorm:"default:0" json:"output_bytes"`
	TotalBytes    int64 `gorm:"default:0" json:"total_bytes"`
	TotalRequests int64 `gorm:"default:0" json:"total_requests"`

	Remark    string         `gorm:"type:text" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Node   *GostNode   `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Tunnel *GostTunnel `gorm:"foreignKey:TunnelID" json:"tunnel,omitempty"`
}

// TableName specifies table name.
func (GostRule) TableName() string {
	return "rules"
}
