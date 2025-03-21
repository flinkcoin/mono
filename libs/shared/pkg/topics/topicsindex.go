// Code generated by "stringer -type=InboundTopic,OutboundTopic -output=topicsindex.go"; DO NOT EDIT.

package topics

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CashierInbound-0]
	_ = x[TraderInbound-1]
}

const _InboundTopic_name = "CashierInboundTraderInbound"

var _InboundTopic_index = [...]uint8{0, 14, 27}

func (i InboundTopic) String() string {
	if i < 0 || i >= InboundTopic(len(_InboundTopic_index)-1) {
		return "InboundTopic(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _InboundTopic_name[_InboundTopic_index[i]:_InboundTopic_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CashierOutbound-0]
	_ = x[TraderOutbound-1]
}

const _OutboundTopic_name = "CashierOutboundTraderOutbound"

var _OutboundTopic_index = [...]uint8{0, 15, 29}

func (i OutboundTopic) String() string {
	if i < 0 || i >= OutboundTopic(len(_OutboundTopic_index)-1) {
		return "OutboundTopic(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OutboundTopic_name[_OutboundTopic_index[i]:_OutboundTopic_index[i+1]]
}
