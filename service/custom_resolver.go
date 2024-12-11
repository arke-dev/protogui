package service

import (
	"fmt"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/types/descriptorpb"
)

type CustomResolver struct {
	descriptors map[string]*descriptorpb.FileDescriptorProto
}

func NewCustomResolver(fds *descriptorpb.FileDescriptorSet) *CustomResolver {
	resolver := &CustomResolver{
		descriptors: make(map[string]*descriptorpb.FileDescriptorProto),
	}
	for _, fd := range fds.File {
		resolver.descriptors[fd.GetName()] = fd
	}
	return resolver
}

func (r *CustomResolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	if fd, ok := r.descriptors[path]; ok {
		return protocompile.SearchResult{
			Proto: fd,
		}, nil
	}
	return protocompile.SearchResult{}, fmt.Errorf("file %s not found", path)
}
