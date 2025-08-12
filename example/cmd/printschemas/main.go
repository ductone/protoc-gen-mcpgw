package main

import (
	"encoding/json"
	"fmt"
	"sort"

	v1 "github.com/ductone/protoc-gen-mcpgw/example/bookstore/v1"
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
)

type registrar struct {
	methodDescs map[string]*mcpgw_v1.MethodDesc
}

func newRegistrar() *registrar {
	return &registrar{methodDescs: make(map[string]*mcpgw_v1.MethodDesc)}
}

func (r *registrar) RegisterService(sd *mcpgw_v1.ServiceDesc, ss any) {
	for _, md := range sd.Methods {
		r.methodDescs[md.Method] = md
	}
}

// nop server satisfies the interface; methods will not be invoked in this tool
type nopServer struct {
	v1.UnimplementedBookstoreServiceServer
}

func main() {
	reg := newRegistrar()

	// Register service (server value is not used by our registrar)
	v1.RegisterMCPBookstoreServiceServer(reg, &nopServer{})

	// Collect method names for stable output
	names := make([]string, 0, len(reg.methodDescs))
	for name := range reg.methodDescs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		md := reg.methodDescs[name]
		fmt.Printf("\n=== %s ===\n", name)
		if md.Title != "" {
			fmt.Printf("title: %s\n", md.Title)
		}
		if md.Description != "" {
			fmt.Printf("description: %s\n", md.Description)
		}

		if md.InputSchema == nil {
			fmt.Println("(no InputSchema function)")
			continue
		}

		schema := md.InputSchema()
		data, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			fmt.Printf("error marshaling schema: %v\n", err)
			continue
		}
		fmt.Println(string(data))
	}
}
