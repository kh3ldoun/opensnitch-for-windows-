package events

import "time"

type Decision string

const (
	DecisionAllowOnce   Decision = "allow_once"
	DecisionAllowAlways Decision = "allow_always"
	DecisionDenyOnce    Decision = "deny_once"
	DecisionDenyAlways  Decision = "deny_always"
)

type ConnectionEvent struct {
	ID          string    `json:"id"`
	NodeID      string    `json:"node_id"`
	ProcessPath string    `json:"process_path"`
	Domain      string    `json:"domain"`
	DstIP       string    `json:"dst_ip"`
	DstPort     uint16    `json:"dst_port"`
	Protocol    string    `json:"protocol"`
	IPv6        bool      `json:"ipv6"`
	Timestamp   time.Time `json:"timestamp"`
	State       string    `json:"state"`
}
