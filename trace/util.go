package trace

import "log/uuid"

func NewTraceId() string {
	return string(uuid.NewV1().HexEncode())
}
