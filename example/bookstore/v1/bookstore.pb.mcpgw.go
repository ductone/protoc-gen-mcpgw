// Code generated by protoc-gen-mcpgw 0.1.0 from bookstore/v1/bookstore.proto. DO NOT EDIT
package v1

import (
	context "context"

	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
	mcpgw_schema "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1/schema"
	grpc "google.golang.org/grpc"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

func RegisterMCPBookstoreServiceServer(s mcpgw_v1.ServiceRegistrar, srv BookstoreServiceServer) {
	s.RegisterService(&mcpgw_desc_BookstoreServiceServer, srv)
}

var mcpgw_desc_BookstoreServiceServer = mcpgw_v1.ServiceDesc{
	Name:        "bookstore.v1.BookstoreService",
	HandlerType: (*BookstoreServiceServer)(nil),
	Methods: []*mcpgw_v1.MethodDesc{
		{
			Method:        BookstoreService_ListShelves_FullMethodName,
			Handler:       _BookstoreService_ListShelves_MCPGW_Handler,
			Decoder:       _BookstoreService_ListShelves_MCPGW_Decoder,
			InputSchema:   _BookstoreService_ListShelves_MCPGW_InputSchema,
			Title:         "List Shelves",
			Description:   "List all shelves in the bookstore",
			ReadOnlyHint:  true,
			Destructive:   false,
			Idempotent:    true,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_CreateShelf_FullMethodName,
			Handler:       _BookstoreService_CreateShelf_MCPGW_Handler,
			Decoder:       _BookstoreService_CreateShelf_MCPGW_Decoder,
			InputSchema:   _BookstoreService_CreateShelf_MCPGW_InputSchema,
			Title:         "Create Shelf",
			Description:   "Create a new shelf in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   false,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_DeleteShelf_FullMethodName,
			Handler:       _BookstoreService_DeleteShelf_MCPGW_Handler,
			Decoder:       _BookstoreService_DeleteShelf_MCPGW_Decoder,
			InputSchema:   _BookstoreService_DeleteShelf_MCPGW_InputSchema,
			Title:         "Delete Shelf",
			Description:   "Delete a shelf in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   true,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_ListGenres_FullMethodName,
			Handler:       _BookstoreService_ListGenres_MCPGW_Handler,
			Decoder:       _BookstoreService_ListGenres_MCPGW_Decoder,
			InputSchema:   _BookstoreService_ListGenres_MCPGW_InputSchema,
			Title:         "List Genres",
			Description:   "List all genres in the bookstore",
			ReadOnlyHint:  true,
			Destructive:   false,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_CreateGenre_FullMethodName,
			Handler:       _BookstoreService_CreateGenre_MCPGW_Handler,
			Decoder:       _BookstoreService_CreateGenre_MCPGW_Decoder,
			InputSchema:   _BookstoreService_CreateGenre_MCPGW_InputSchema,
			Title:         "Create Genre",
			Description:   "Create a new genre in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   false,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_GetGenre_FullMethodName,
			Handler:       _BookstoreService_GetGenre_MCPGW_Handler,
			Decoder:       _BookstoreService_GetGenre_MCPGW_Decoder,
			InputSchema:   _BookstoreService_GetGenre_MCPGW_InputSchema,
			Title:         "Get Genre",
			Description:   "Get a genre in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   false,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_DeleteGenre_FullMethodName,
			Handler:       _BookstoreService_DeleteGenre_MCPGW_Handler,
			Decoder:       _BookstoreService_DeleteGenre_MCPGW_Decoder,
			InputSchema:   _BookstoreService_DeleteGenre_MCPGW_InputSchema,
			Title:         "Delete Genre",
			Description:   "Delete a genre in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   true,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_CreateBook_FullMethodName,
			Handler:       _BookstoreService_CreateBook_MCPGW_Handler,
			Decoder:       _BookstoreService_CreateBook_MCPGW_Decoder,
			InputSchema:   _BookstoreService_CreateBook_MCPGW_InputSchema,
			Title:         "Create Book",
			Description:   "Create a new book in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   true,
			Idempotent:    true,
			OpenWorldHint: true,
		},
		{
			Method:        BookstoreService_GetBook_FullMethodName,
			Handler:       _BookstoreService_GetBook_MCPGW_Handler,
			Decoder:       _BookstoreService_GetBook_MCPGW_Decoder,
			InputSchema:   _BookstoreService_GetBook_MCPGW_InputSchema,
			Title:         "Get Book",
			Description:   "Get a book in the bookstore",
			ReadOnlyHint:  true,
			Destructive:   false,
			Idempotent:    true,
			OpenWorldHint: true,
		},
		{
			Method:        BookstoreService_ListBooks_FullMethodName,
			Handler:       _BookstoreService_ListBooks_MCPGW_Handler,
			Decoder:       _BookstoreService_ListBooks_MCPGW_Decoder,
			InputSchema:   _BookstoreService_ListBooks_MCPGW_InputSchema,
			Title:         "List Books",
			Description:   "List all books in the bookstore",
			ReadOnlyHint:  true,
			Destructive:   false,
			Idempotent:    true,
			OpenWorldHint: true,
		},
		{
			Method:        BookstoreService_DeleteBook_FullMethodName,
			Handler:       _BookstoreService_DeleteBook_MCPGW_Handler,
			Decoder:       _BookstoreService_DeleteBook_MCPGW_Decoder,
			InputSchema:   _BookstoreService_DeleteBook_MCPGW_InputSchema,
			Title:         "Delete Book",
			Description:   "Delete a book in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   true,
			Idempotent:    false,
			OpenWorldHint: false,
		},
		{
			Method:        BookstoreService_UpdateBook_FullMethodName,
			Handler:       _BookstoreService_UpdateBook_MCPGW_Handler,
			Decoder:       _BookstoreService_UpdateBook_MCPGW_Decoder,
			InputSchema:   _BookstoreService_UpdateBook_MCPGW_InputSchema,
			Title:         "Update Book",
			Description:   "Update a book in the bookstore",
			ReadOnlyHint:  false,
			Destructive:   true,
			Idempotent:    true,
			OpenWorldHint: true,
		},
	},
}

func _BookstoreService_ListShelves_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*ListShelvesRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&ListShelvesRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_ListShelves_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(ListShelvesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).ListShelves(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_ListShelves_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).ListShelves(ctx, req.(*ListShelvesRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_ListShelves_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_CreateShelf_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*CreateShelfRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&CreateShelfRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_CreateShelf_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(CreateShelfRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).CreateShelf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_CreateShelf_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).CreateShelf(ctx, req.(*CreateShelfRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_CreateShelf_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_DeleteShelf_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*DeleteShelfRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&DeleteShelfRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_DeleteShelf_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(DeleteShelfRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).DeleteShelf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_DeleteShelf_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).DeleteShelf(ctx, req.(*DeleteShelfRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_DeleteShelf_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_ListGenres_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*ListGenresRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&ListGenresRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_ListGenres_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(ListGenresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).ListGenres(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_ListGenres_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).ListGenres(ctx, req.(*ListGenresRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_ListGenres_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_CreateGenre_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*CreateGenreRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&CreateGenreRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_CreateGenre_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(CreateGenreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).CreateGenre(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_CreateGenre_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).CreateGenre(ctx, req.(*CreateGenreRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_CreateGenre_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_GetGenre_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*GetGenreRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&GetGenreRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_GetGenre_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(GetGenreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).GetGenre(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_GetGenre_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).GetGenre(ctx, req.(*GetGenreRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_GetGenre_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_DeleteGenre_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*DeleteGenreRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&DeleteGenreRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_DeleteGenre_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(DeleteGenreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).DeleteGenre(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_DeleteGenre_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).DeleteGenre(ctx, req.(*DeleteGenreRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_DeleteGenre_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_CreateBook_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*CreateBookRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&CreateBookRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_CreateBook_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(CreateBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).CreateBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_CreateBook_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).CreateBook(ctx, req.(*CreateBookRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_CreateBook_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_GetBook_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*GetBookRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&GetBookRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_GetBook_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(GetBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).GetBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_GetBook_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).GetBook(ctx, req.(*GetBookRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_GetBook_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_ListBooks_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*ListBooksRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&ListBooksRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_ListBooks_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(ListBooksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).ListBooks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_ListBooks_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).ListBooks(ctx, req.(*ListBooksRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_ListBooks_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_DeleteBook_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*DeleteBookRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&DeleteBookRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_DeleteBook_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(DeleteBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).DeleteBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_DeleteBook_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).DeleteBook(ctx, req.(*DeleteBookRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_DeleteBook_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}

func _BookstoreService_UpdateBook_MCPGW_InputSchema() map[string]any {
	return mcpgw_schema.MustGenerateSchema(((*UpdateBookRequest)(nil)).ProtoReflect().Descriptor())
	// return mcpgw_schema.MustGenerateSchema((&UpdateBookRequest{}).ProtoReflect().Descriptor())
}

func _BookstoreService_UpdateBook_MCPGW_Handler(srv any, ctx context.Context, dec func(proto.Message) error, interceptor grpc.UnaryServerInterceptor) (proto.Message, error) {
	in := new(UpdateBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookstoreServiceServer).UpdateBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookstoreService_UpdateBook_FullMethodName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(BookstoreServiceServer).UpdateBook(ctx, req.(*UpdateBookRequest))
	}

	rv, err := interceptor(ctx, in, info, handler)
	if err != nil {
		return nil, err
	}
	return rv.(proto.Message), nil
}

func _BookstoreService_UpdateBook_MCPGW_Decoder(ctx context.Context, input mcpgw_v1.DecoderInput, out proto.Message) error {
	var err error
	_ = err

	if len(input.RawArguments()) > 0 {
		return protojson.Unmarshal(input.RawArguments(), out)
	}

	return mcpgw_v1.UnmarshalFromMap(input.Arguments(), out)
}
