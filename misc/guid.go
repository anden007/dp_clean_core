package misc

import "github.com/google/uuid"

func NewGuid() (result uuid.UUID) {
	result, _ = uuid.NewUUID()
	return
}

func NewGuidString() (result string) {
	result = ""
	uid, err := uuid.NewUUID()
	if err == nil {
		result = uid.String()
	}
	return
}
