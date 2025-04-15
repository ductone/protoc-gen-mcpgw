package mcpgw

import (
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
)

func getServiceOptions(s pgs.Service) *mcpgw_v1.ServiceOptions {
	sopt := &mcpgw_v1.ServiceOptions{}
	_, err := s.Extension(mcpgw_v1.E_Service, sopt)
	if err != nil {
		return nil
	}
	return sopt
}

func getMethodOptions(m pgs.Method) *mcpgw_v1.MethodOptions {
	mopt := &mcpgw_v1.MethodOptions{}
	_, err := m.Extension(mcpgw_v1.E_Method, mopt)
	if err != nil {
		return nil
	}
	return mopt
}

func getMessageOptions(m pgs.Message) *mcpgw_v1.MessageOptions {
	mopt := &mcpgw_v1.MessageOptions{}
	_, err := m.Extension(mcpgw_v1.E_Message, mopt)
	if err != nil {
		return nil
	}
	return mopt
}

func getFieldOptions(m pgs.Field) *mcpgw_v1.FieldOptions {
	fopt := &mcpgw_v1.FieldOptions{}
	_, err := m.Extension(mcpgw_v1.E_Field, fopt)
	if err != nil {
		return nil
	}
	return fopt
}
