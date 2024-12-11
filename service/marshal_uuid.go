package service

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type UUIDMarshaller struct {
	fieldsType   map[string]protoreflect.Kind
	protoMessage protoreflect.ProtoMessage
}

func NewUUIDMarshaller(protoMessage protoreflect.ProtoMessage) *UUIDMarshaller {
	return &UUIDMarshaller{
		fieldsType:   make(map[string]protoreflect.Kind),
		protoMessage: protoMessage,
	}
}

func (unmarshal *UUIDMarshaller) UUIDs(jsonMessage map[string]any) {
	unmarshal.buildFieldsTypeFromDescriptors("", unmarshal.protoMessage.ProtoReflect().Descriptor().Fields())
	unmarshal.uuidsInMap("", jsonMessage)
}

func (unmarshal *UUIDMarshaller) buildFieldsTypeFromDescriptors(root string, mm protoreflect.FieldDescriptors) {
	for i := range mm.Len() {
		field := mm.Get(i)
		if root == "" {
			unmarshal.fieldsType[string(field.JSONName())] = field.Kind()
		} else {
			fieldKey := fmt.Sprintf("%s.%s", root, field.JSONName())
			unmarshal.fieldsType[fieldKey] = field.Kind()
		}

		if field.Kind() == protoreflect.MessageKind {
			var newRoot string
			if root != "" {
				newRoot = fmt.Sprintf("%s.%s", root, field.JSONName())
			} else {
				newRoot = string(field.JSONName())
			}

			unmarshal.buildFieldsTypeFromDescriptors(newRoot, field.Message().Fields())
		}
	}
}

func (unmarshal *UUIDMarshaller) uuidsInMap(key string, mapMsg map[string]any) {
	for k, item := range mapMsg {
		newKey := fmt.Sprintf("%s.%s", key, k)
		if key == "" {
			newKey = k
		}
		s, ok := unmarshal.nextUUID(newKey, item)
		if ok {
			mapMsg[k] = s
		}
	}
}

func (unmarshal *UUIDMarshaller) uuidsInSlice(key string, items []any) {
	for i, item := range items {
		s, ok := unmarshal.nextUUID(key, item)
		if ok {
			items[i] = s
		}
	}
}

func (unmarshal *UUIDMarshaller) uuid(key string, v any) (string, bool) {
	s, ok := v.(string)
	if !ok {
		return "", false
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return "", false
	}

	k, ok := unmarshal.fieldsType[key]
	if !ok {
		return "", false
	}

	if k != protoreflect.BytesKind {
		return "", false
	}

	// MarshalBinary never returns error so it's safe to omit error
	b, _ := id.MarshalBinary()

	bencoded := base64.StdEncoding.EncodeToString(b)

	return bencoded, true
}

func (unmarshal *UUIDMarshaller) nextUUID(key string, v any) (any, bool) {
	m, ok := v.(map[string]any)
	if ok {
		unmarshal.uuidsInMap(key, m)
		return "", false
	}

	sl, ok := v.([]any)
	if ok {
		unmarshal.uuidsInSlice(key, sl)
		return "", false
	}

	return unmarshal.uuid(key, v)

}
