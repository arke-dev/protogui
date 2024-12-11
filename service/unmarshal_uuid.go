package service

import (
	"encoding/base64"

	"github.com/google/uuid"
)

type UUIDUnmarshaller struct{}

func NewUUIDUnmarshaller() *UUIDUnmarshaller {
	return &UUIDUnmarshaller{}
}

func (u *UUIDUnmarshaller) UUIDs(mapMsg map[string]any) {
	u.uuidsInMap(mapMsg)
}

func (u *UUIDUnmarshaller) uuidsInMap(mapMsg map[string]any) {
	for k, item := range mapMsg {
		s, ok := u.nextUUID(item)
		if ok {
			mapMsg[k] = s
		}
	}
}

func (u *UUIDUnmarshaller) uuidsInSlice(items []any) {
	for i, item := range items {
		s, ok := u.nextUUID(item)
		if ok {
			items[i] = s
		}
	}
}

func (u *UUIDUnmarshaller) uuid(v any) (string, bool) {
	s, ok := v.(string)
	if !ok {
		return "", false
	}

	bbb, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", false
	}

	uid, err := uuid.FromBytes(bbb)
	if err != nil {
		return "", false
	}

	return uid.String(), true
}

func (u *UUIDUnmarshaller) nextUUID(v any) (string, bool) {
	m, ok := v.(map[string]any)
	if ok {
		u.uuidsInMap(m)
		return "", false
	}

	sl, ok := v.([]any)
	if ok {
		u.uuidsInSlice(sl)
		return "", false
	}

	return u.uuid(v)
}
