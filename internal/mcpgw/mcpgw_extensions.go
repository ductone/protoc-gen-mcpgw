package mcpgw

import (
	apigw_v1 "github.com/ductone/protoc-gen-apigw/apigw/v1"
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func getServiceOptions(s pgs.Service) *mcpgw_v1.ServiceOptions {
	sopt := &mcpgw_v1.ServiceOptions{}
	ok, err := s.Extension(mcpgw_v1.E_Service, sopt)
	if err != nil || !ok {
		return nil
	}
	return sopt
}

func getMethodOptions(m pgs.Method) *mcpgw_v1.MethodOptions {
	mopt := &mcpgw_v1.MethodOptions{}
	ok, err := m.Extension(mcpgw_v1.E_Method, mopt)
	if err != nil || !ok {
		return nil
	}
	return mopt
}

// getApigwOperation returns the first apigw.v1.Operation declared on the
// method, or nil if the method has no apigw annotation. Services can have
// multiple operations (e.g. GET and POST to the same route); for MCP tool
// derivation we key on the first, which is the canonical route.
func getApigwOperation(m pgs.Method) *apigw_v1.Operation {
	mopt := &apigw_v1.MethodOptions{}
	ok, err := m.Extension(apigw_v1.E_Method, mopt)
	if err != nil || !ok {
		return nil
	}
	ops := mopt.GetOperations()
	if len(ops) == 0 {
		return nil
	}
	return ops[0]
}

// serviceEnabled reports whether a service should be rendered as an MCP
// service. A service is enabled when either:
//  1. It has (mcpgw.v1.service) = {enabled: true} — explicit opt-in.
//  2. At least one of its methods has an (apigw.v1.method) annotation —
//     implicit opt-in for REST-exposed services.
//
// An explicit {enabled: false} is respected as a hard opt-out regardless
// of apigw coverage.
func serviceEnabled(s pgs.Service) bool {
	sopt := getServiceOptions(s)
	if sopt != nil && sopt.HasEnabled() {
		return sopt.GetEnabled()
	}
	for _, m := range s.Methods() {
		if getApigwOperation(m) != nil {
			return true
		}
	}
	return false
}

func getMessageOptions(m pgs.Message) *mcpgw_v1.MessageOptions {
	mopt := &mcpgw_v1.MessageOptions{}
	ok, err := m.Extension(mcpgw_v1.E_Message, mopt)
	if err != nil || !ok {
		return nil
	}
	return mopt
}

func getFieldOptions(m pgs.Field) *mcpgw_v1.FieldOptions {
	fopt := &mcpgw_v1.FieldOptions{}
	ok, err := m.Extension(mcpgw_v1.E_Field, fopt)
	if err != nil || !ok {
		return nil
	}
	return fopt
}
