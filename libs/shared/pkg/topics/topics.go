//go:generate stringer -type=InboundTopic,OutboundTopic -output=topicsindex.go
package topics

import ()

type InboundTopic int

const (
	CashierInbound InboundTopic = iota
	TraderInbound
)

type OutboundTopic int

const (
	CashierOutbound OutboundTopic = iota
	TraderOutbound
)
