package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bufbuild/protocompile"
	"github.com/lawmatsuyama/protogui/helpers"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type ProtoCompiler interface {
	DecodeNoBase64(path string, tp string, message []byte) (string, error)
	Decode(path string, tp string, message string) (string, error)
	Encode(path string, tp string, jsonMessageString string) (string, error)
	RegisterProto(path string) error
	JSONToProto(path string, tp string, jsonMessageString string) (protoreflect.ProtoMessage, error)
	ProtoToJSON(msg protoreflect.ProtoMessage) (string, error)
	GetMessageType(path string, typename string) (protoreflect.MessageType, error)
	TemplateJSON(path string, typename string) (string, error)
	TemplateJSONFromMethod(path string, method string) (string, error)
	GetRegisteredTypes(path string) ([]string, error)
	GetRequestResponseFromMethod(path string, method string) (requestTypename string, responseTypename string, err error)
	GetRegisteredMethods(path string) ([]string, error)
}

func NewProtoCompile() ProtoCompiler {
	return &protoCompiler{
		mux:     &sync.Mutex{},
		methods: make(map[string]protoreflect.MethodDescriptor),
	}
}

type protoCompiler struct {
	mux     *sync.Mutex
	methods map[string]protoreflect.MethodDescriptor
}

func (p *protoCompiler) Decode(path string, typename string, message string) (string, error) {
	messageType, err := p.GetMessageType(path, typename)
	if err != nil {
		return "", fmt.Errorf("failed to register proto %w", err)
	}

	if messageType == nil {
		return "", fmt.Errorf("unknown message type %s", typename)
	}

	msg := messageType.New().Interface()

	message = strings.Trim(strings.TrimSpace(message), "\n")
	bb, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Printf("maybe message is not base64 %v\n", err)
		bb = []byte(message)
	}

	err = proto.Unmarshal(bb, msg)
	if err != nil {
		return "", fmt.Errorf("failed to proto unmarshal %w", err)
	}

	return p.ProtoToJSON(msg)
}

func (p *protoCompiler) DecodeNoBase64(path string, typename string, message []byte) (string, error) {
	messageType, err := p.GetMessageType(path, typename)
	if err != nil {
		return "", fmt.Errorf("failed to register proto %w", err)
	}

	if messageType == nil {
		return "", fmt.Errorf("unknown message type %s", typename)
	}
	msg := messageType.New().Interface()
	err = proto.Unmarshal(message, msg)
	if err != nil {
		return "", fmt.Errorf("failed to proto unmarshal %w", err)
	}

	return p.ProtoToJSON(msg)
}

func (p *protoCompiler) Encode(filePath string, tp string, jsonMessageString string) (string, error) {
	protoMessage, err := p.JSONToProto(filePath, tp, jsonMessageString)
	if err != nil {
		return "", err
	}

	b, err := proto.Marshal(protoMessage)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func (p *protoCompiler) GetMessageType(path string, typename string) (protoreflect.MessageType, error) {
	if path != "" {
		err := p.RegisterProto(path)
		if err != nil {
			return nil, err
		}
	}

	messageType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(typename))
	if err != nil {
		fmt.Printf("failed to find message type %v\n", err)
		return nil, err
	}

	return messageType, nil
}

func (p *protoCompiler) RegisterProto(path string) error {
	if path == "" {
		return nil
	}

	home := os.Getenv("HOME")

	importPaths := []string{
		path,
		home + "/.local/include",
		"/usr/include",
	}

	dirs, ok := helpers.WalkDeepDirectory(path, ".github")
	if ok {
		importPaths = append(importPaths, dirs...)
	}

	fsResolvers, err := p.fileDesriptorResolvers(path)
	if err != nil {
		return err
	}

	resolvers := protocompile.CompositeResolver{
		&protocompile.SourceResolver{
			ImportPaths: importPaths,
		},
	}

	resolvers = append(resolvers, fsResolvers...)

	compiler := protocompile.Compiler{
		Resolver: resolvers,
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir %w", err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), ".proto") {
			files = append(files, entry.Name())
		}
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	protoregistry.GlobalTypes = new(protoregistry.Types)
	p.methods = make(map[string]protoreflect.MethodDescriptor)
	for _, file := range files {
		f, err := compiler.Compile(context.Background(), file)
		if err != nil {
			fmt.Printf("failed to compile protobuf file %v", err)
			return err
		}

		ff := f.FindFileByPath(file)
		for i := range ff.Messages().Len() {
			msgType, err := f.AsResolver().FindMessageByName(ff.Messages().Get(i).FullName())
			if err != nil {
				fmt.Printf("failed to find message by name %v", err)
				return err
			}

			err = protoregistry.GlobalTypes.RegisterMessage(msgType)
			if err != nil {
				fmt.Printf("failed to register message %v", err)
				return err
			}
		}

		for i := range ff.Services().Len() {
			for j := range ff.Services().Get(i).Methods().Len() {
				methodDescript := ff.Services().Get(i).Methods().Get(j)
				serviceName := ff.Services().Get(i).FullName()
				methodName := methodDescript.FullName().Name()

				methodFullName := fmt.Sprintf("/%s/%s", serviceName, methodName)

				p.methods[methodFullName] = methodDescript
			}
		}

	}

	return err
}

func (p *protoCompiler) GetRegisteredTypes(path string) ([]string, error) {
	err := p.RegisterProto(path)
	if err != nil {
		return nil, fmt.Errorf("failed to register proto %v", err)
	}

	types := make([]string, 0)
	protoregistry.GlobalTypes.RangeMessages(func(mt protoreflect.MessageType) bool {
		types = append(types, string(mt.Descriptor().FullName()))
		return true
	})

	return types, nil
}

func (p *protoCompiler) JSONToProto(filePath string, tp string, jsonMessageString string) (protoreflect.ProtoMessage, error) {
	messageType, err := p.GetMessageType(filePath, tp)
	if err != nil {
		return nil, fmt.Errorf("failed to register proto %w", err)
	}

	protoMessage := messageType.New().Interface()

	jsonMessage := map[string]any{}
	err = json.Unmarshal([]byte(jsonMessageString), &jsonMessage)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to unmarshal json message to map %w", err)
	}

	marshal := NewUUIDMarshaller(protoMessage)
	marshal.UUIDs(jsonMessage)

	jsonMessageByte, err := json.Marshal(jsonMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json message %w", err)
	}

	err = protojson.Unmarshal(jsonMessageByte, protoMessage)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to unmarshal json message to proto %w", err)
	}

	return protoMessage, nil
}

func (p *protoCompiler) ProtoToJSON(msg protoreflect.ProtoMessage) (string, error) {
	jb, err := protojson.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to protojson marshal %w", err)
	}

	mapMsg := make(map[string]any)

	err = json.Unmarshal(jb, &mapMsg)
	if err != nil {
		return "", fmt.Errorf("failed to json unmarshal %w", err)
	}

	unmarshal := NewUUIDUnmarshaller()
	unmarshal.UUIDs(mapMsg)

	jb2, err := json.MarshalIndent(mapMsg, "", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to json marshal indent %w", err)
	}

	return string(jb2), nil
}

func (p *protoCompiler) TemplateJSON(path string, typename string) (string, error) {
	msg, err := p.GetMessageType(path, typename)
	if err != nil {
		return "", err
	}

	template := msg.New().Interface()

	setDefaultValuesToTemplate(template)

	marshaller := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}

	b, err := marshaller.Marshal(template)
	if err != nil {
		return "", err
	}

	return string(b), nil

}

func (p *protoCompiler) TemplateJSONFromMethod(path string, method string) (string, error) {
	err := p.RegisterProto(path)
	if err != nil {
		return "", err
	}

	met, ok := p.getMethod(method)
	if !ok {
		return "", errors.New("method not found in template json")
	}

	inputTypename := met.Input().FullName()

	return p.TemplateJSON("", string(inputTypename))
}

func setDefaultValuesToTemplate(m proto.Message) {
	descriptor := m.ProtoReflect().Descriptor()
	for i := range descriptor.Fields().Len() {
		field := descriptor.Fields().Get(i)

		if field.IsMap() {
			if field.MapValue().Kind() != protoreflect.MessageKind {
				continue
			}

			k := field.MapKey().Default().MapKey()
			x := m.ProtoReflect().Mutable(field).Map().Mutable(protoreflect.MapKey(k)).Message()

			setDefaultValuesToTemplate(x.Interface())
			continue
		}

		if !field.IsList() && field.Kind() == protoreflect.MessageKind {
			nestedMessage := m.ProtoReflect().Mutable(field).Message()
			setDefaultValuesToTemplate(nestedMessage.Interface())
			continue
		}

		if field.IsList() && field.Kind() == protoreflect.MessageKind {
			newMsg := m.ProtoReflect().Mutable(field).List().AppendMutable()
			setDefaultValuesToTemplate(newMsg.Message().Interface())
		}

	}
}

func (p *protoCompiler) getResolver(filePath string) (*CustomResolver, error) {
	descriptorBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to open descriptor set: %v\n", err)
		return nil, err
	}

	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(descriptorBytes, fds)
	if err != nil {
		fmt.Printf("Failed to unmarshal descriptor set: %v\n", err)
		return nil, err
	}

	resolver := NewCustomResolver(fds)
	return resolver, nil
}

func (p *protoCompiler) fileDesriptorResolvers(path string) ([]protocompile.Resolver, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	resolvers := make([]protocompile.Resolver, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".def") {
			continue
		}

		resolver, err := p.getResolver(path + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		resolvers = append(resolvers, resolver)

	}

	return resolvers, nil
}

func (p *protoCompiler) GetRequestResponseFromMethod(path string, method string) (requestTypename string, responseTypename string, err error) {
	err = p.RegisterProto(path)
	if err != nil {
		return "", "", err
	}

	met, ok := p.getMethod(method)
	if !ok {
		return "", "", errors.New("method not found")
	}

	return string(met.Input().FullName()), string(met.Output().FullName()), nil
}

func (p *protoCompiler) GetRegisteredMethods(path string) ([]string, error) {
	err := p.RegisterProto(path)
	if err != nil {
		return nil, err
	}

	methods := make([]string, len(p.methods))
	i := 0
	for k := range p.methods {
		methods[i] = string(k)
		i++
	}

	return methods, nil
}

func (p *protoCompiler) getMethod(method string) (protoreflect.MethodDescriptor, bool) {
	met, ok := p.methods[method]
	return met, ok
}

// func (p *protoCompiler) getMethodFullName(desc protoreflect.MethodDescriptor) string {

// }
