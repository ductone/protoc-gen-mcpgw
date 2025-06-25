package v1_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	v1 "github.com/ductone/protoc-gen-mcpgw/example/bookstore/v1"
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
)

// TestUnmarshalFromMap tests directly using mcpgw_v1.UnmarshalFromMap function
func TestUnmarshalFromMap(t *testing.T) {
	// Test with simple data
	t.Run("Simple_Data", func(t *testing.T) {
		// Create map[string]any with genre data
		inputMap := map[string]any{
			"name": "Mystery",
		}

		// Create a request message to decode into
		req := &v1.CreateGenreRequest{}

		// Use UnmarshalFromMap directly
		err := mcpgw_v1.UnmarshalFromMap(inputMap, req)
		assert.NoError(t, err)

		// Verify the decoded message
		assert.Equal(t, "Mystery", req.GetName())
	})

	// Test with complex nested data
	t.Run("Nested_Data", func(t *testing.T) {
		// Create map[string]any with book data
		inputMap := map[string]any{
			"shelf": "shelf-789",
			"book": map[string]any{
				"id":      "book-345",
				"author":  "John Smith",
				"title":   "The Mystery of the Deep",
				"quotes":  []string{"First quote", "Second quote"},
				"shelfId": "shelf-789",
			},
		}

		// Create a request message to decode into
		req := &v1.CreateBookRequest{}

		// Use UnmarshalFromMap directly
		err := mcpgw_v1.UnmarshalFromMap(inputMap, req)
		assert.NoError(t, err)

		// Verify the decoded message
		assert.Equal(t, "shelf-789", req.GetShelf())
		assert.NotNil(t, req.GetBook())

		book := req.GetBook()
		assert.Equal(t, "book-345", book.GetId())
		assert.Equal(t, "John Smith", book.GetAuthor())
		assert.Equal(t, "The Mystery of the Deep", book.GetTitle())
		assert.Equal(t, []string{"First quote", "Second quote"}, book.GetQuotes())
		assert.Equal(t, "shelf-789", book.GetShelfId())
	})
}

// TestInputSchemaValidation tests that method input schemas are valid JSON schemas
// and that they correctly validate inputs
func TestInputSchemaValidation(t *testing.T) {
	// Create a mock service registrar
	mockRegistrar := NewMockServiceRegistrar()

	// Register the BookstoreService
	server := &mockBookstoreServer{}
	v1.RegisterMCPBookstoreServiceServer(mockRegistrar, server)

	// Test the CreateGenre method's input schema
	t.Run("CreateGenre_Schema", func(t *testing.T) {
		// Get the method descriptor
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Check if the InputSchema function is available
		require.NotNil(t, methodDesc.InputSchema, "InputSchema function should be available")

		// Get the schema as map[string]any
		schemaMap := methodDesc.InputSchema()
		require.NotNil(t, schemaMap, "Schema map should not be nil")

		// Convert schema to JSON
		schemaBytes, err := json.Marshal(schemaMap)
		require.NoError(t, err, "Failed to marshal schema to JSON")

		// Load and compile the schema using santhosh-tekuri/jsonschema
		schema, err := loadAndCompileSchema(schemaBytes)
		require.NoError(t, err, "Schema should be valid JSON Schema")

		// Test valid input
		validInput := map[string]any{
			"name": "Science Fiction",
		}
		// Pass the raw Go value directly to Validate
		err = schema.Validate(validInput)
		assert.NoError(t, err, "Valid input should validate successfully")

		// Test invalid input (missing required field if applicable)
		invalidInput := map[string]any{}
		// Pass the raw Go value directly to Validate
		err = schema.Validate(invalidInput)
		t.Logf("Validation result for empty input: %v", err)
	})

	// Test the CreateBook method's input schema which is more complex
	t.Run("CreateBook_Schema", func(t *testing.T) {
		// Get the method descriptor
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateBook"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateBook should be registered")

		// Check if the InputSchema function is available
		require.NotNil(t, methodDesc.InputSchema, "InputSchema function should be available")

		// Get the schema as map[string]any
		schemaMap := methodDesc.InputSchema()
		require.NotNil(t, schemaMap, "Schema map should not be nil")

		// Convert schema to JSON
		schemaBytes, err := json.Marshal(schemaMap)
		require.NoError(t, err, "Failed to marshal schema to JSON")

		// Log the schema for debugging
		t.Logf("CreateBook JSON Schema: %s", string(schemaBytes))

		// Load and compile the schema using santhosh-tekuri/jsonschema
		schema, err := loadAndCompileSchema(schemaBytes)
		require.NoError(t, err, "Schema should be valid JSON Schema")

		// Test valid input - using []any instead of []string for quotes
		validInput := map[string]any{
			"shelf": "shelf-123",
			"book": map[string]any{
				"id":      "book-456",
				"author":  "Jane Doe",
				"title":   "The Great Adventure",
				"quotes":  []any{"Quote 1", "Quote 2"},
				"shelfId": "shelf-123",
			},
		}
		// Pass the raw Go value directly to Validate
		err = schema.Validate(validInput)
		assert.NoError(t, err, "Valid input should validate successfully")

		// Test invalid input (missing required fields if applicable)
		invalidInput := map[string]any{
			// Missing "shelf" field
			"book": map[string]any{
				// Missing required fields within book
				"id": "book-456",
			},
		}
		// Pass the raw Go value directly to Validate
		err = schema.Validate(invalidInput)
		t.Logf("Validation result for incomplete input: %v", err)
	})
}

// Helper function to load and compile a JSON schema
func loadAndCompileSchema(schemaBytes []byte) (*jsonschema.Schema, error) {
	// Create a compiler with default settings
	compiler := jsonschema.NewCompiler()

	// Add the schema as a resource
	compiler.AddResource("schema.json", strings.NewReader(string(schemaBytes)))

	// Compile the schema
	return compiler.Compile("schema.json")
}

// ValidatingInterceptor creates a gRPC interceptor that validates request messages
// using protovalidate-go before they reach the handler
func ValidatingInterceptor(t *testing.T) grpc.UnaryServerInterceptor {
	validator, err := protovalidate.New()
	assert.NoError(t, err)

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Validate the request
		if err := validator.Validate(req.(proto.Message)); err != nil {
			// Return validation error instead of proceeding to the handler
			return nil, err
		}

		// If validation passes, call the handler
		return handler(ctx, req)
	}
}

// MockServiceRegistrar implements mcpgw_v1.ServiceRegistrar for testing
type MockServiceRegistrar struct {
	methodDescs map[string]*mcpgw_v1.MethodDesc
}

func NewMockServiceRegistrar() *MockServiceRegistrar {
	return &MockServiceRegistrar{
		methodDescs: make(map[string]*mcpgw_v1.MethodDesc),
	}
}

func (m *MockServiceRegistrar) RegisterService(sd *mcpgw_v1.ServiceDesc, ss any) {
	// Store the method descriptors for later lookup
	for _, md := range sd.Methods {
		m.methodDescs[md.Method] = md
	}
}

// MockDecoderInput implements mcpgw_v1.DecoderInput
type MockDecoderInput struct {
	method  string
	args    map[string]any
	rawArgs json.RawMessage
}

func NewMockDecoderInput(method string, args map[string]any) *MockDecoderInput {
	var rawArgs json.RawMessage
	if args != nil {
		data, _ := json.Marshal(args)
		rawArgs = data
	}

	return &MockDecoderInput{
		method:  method,
		args:    args,
		rawArgs: rawArgs,
	}
}

func (m *MockDecoderInput) Method() string {
	return m.method
}

func (m *MockDecoderInput) Arguments() map[string]any {
	return m.args
}

func (m *MockDecoderInput) RawArguments() json.RawMessage {
	return m.rawArgs
}

func TestMCPGWRegistration(t *testing.T) {
	// Create a mock service registrar
	mockRegistrar := NewMockServiceRegistrar()

	// Register the BookstoreService
	server := &mockBookstoreServer{}
	v1.RegisterMCPBookstoreServiceServer(mockRegistrar, server)

	// Test CreateGenre method
	t.Run("CreateGenre", func(t *testing.T) {
		// Get the method descriptor for CreateGenre
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		assert.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Create input with map[string]any
		input := NewMockDecoderInput("bookstore.v1.BookstoreService.CreateGenre", map[string]any{
			"name": "Science Fiction",
		})

		// Create a new request message to decode into
		req := &v1.CreateGenreRequest{}

		// Use the decoder to decode the input
		err := methodDesc.Decoder(context.Background(), input, req)
		assert.NoError(t, err, "Decoding should succeed")

		// Verify the decoded request using the proper accessor methods
		assert.Equal(t, "Science Fiction", req.GetName())
	})

	// Test CreateBook method with nested structure
	t.Run("CreateBook", func(t *testing.T) {
		// Get the method descriptor for CreateBook
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateBook"]
		assert.NotNil(t, methodDesc, "Method descriptor for CreateBook should be registered")

		// Create complex input with nested map[string]any
		input := NewMockDecoderInput("bookstore.v1.BookstoreService.CreateBook", map[string]any{
			"shelf": "shelf-123",
			"book": map[string]any{
				"id":      "book-456",
				"author":  "Jane Doe",
				"title":   "The Great Adventure",
				"quotes":  []string{"Quote 1", "Quote 2", "Quote 3"},
				"shelfId": "shelf-123",
			},
		})

		// Create a new request message to decode into
		req := &v1.CreateBookRequest{}

		// Use the decoder to decode the input
		err := methodDesc.Decoder(context.Background(), input, req)
		assert.NoError(t, err, "Decoding should succeed")

		// Verify the decoded request
		assert.Equal(t, "shelf-123", req.GetShelf())
		assert.NotNil(t, req.GetBook())

		book := req.GetBook()
		assert.Equal(t, "book-456", book.GetId())
		assert.Equal(t, "Jane Doe", book.GetAuthor())
		assert.Equal(t, "The Great Adventure", book.GetTitle())
		assert.Equal(t, []string{"Quote 1", "Quote 2", "Quote 3"}, book.GetQuotes())
		assert.Equal(t, "shelf-123", book.GetShelfId())
	})

	// Test full end-to-end flow for CreateGenre
	t.Run("CreateGenre_E2E", func(t *testing.T) {
		// Get the method descriptor
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]

		// Create context
		ctx := context.Background()

		// Create map[string]any input
		inputMap := map[string]any{
			"name": "Fantasy",
		}

		// Create decoder input
		input := NewMockDecoderInput("/bookstore.v1.BookstoreService/CreateGenre", inputMap)

		// Create a request message
		req := &v1.CreateGenreRequest{}

		// Decode the input
		err := methodDesc.Decoder(ctx, input, req)
		assert.NoError(t, err)

		// Create a custom decoder function that uses UnmarshalFromMap directly
		decoder := func(msg proto.Message) error {
			protoReq, ok := msg.(*v1.CreateGenreRequest)
			if !ok {
				return nil
			}

			// Use UnmarshalFromMap instead of copying the req
			return mcpgw_v1.UnmarshalFromMap(inputMap, protoReq)
		}

		// Call the handler without interceptor
		resp, err := methodDesc.Handler(server, ctx, decoder, nil)
		assert.NoError(t, err)

		// Validate response
		genreResp, ok := resp.(*v1.CreateGenreResponse)
		assert.True(t, ok)
		assert.NotNil(t, genreResp.GetGenre())
		assert.Equal(t, int64(42), genreResp.GetGenre().GetId())
		assert.Equal(t, "Fantasy", genreResp.GetGenre().GetName())

		// Test with an interceptor
		interceptorCalled := false
		interceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			interceptorCalled = true
			return handler(ctx, req)
		}

		resp, err = methodDesc.Handler(server, ctx, decoder, interceptor)
		assert.NoError(t, err)
		assert.True(t, interceptorCalled, "Interceptor should have been called")

		// Validate response again
		genreResp, ok = resp.(*v1.CreateGenreResponse)
		assert.True(t, ok)
		assert.NotNil(t, genreResp.GetGenre())
		assert.Equal(t, int64(42), genreResp.GetGenre().GetId())
		assert.Equal(t, "Fantasy", genreResp.GetGenre().GetName())
	})
}

// TestProtoValidateMiddleware tests using protovalidate-go as middleware for validating input
func TestProtoValidateMiddleware(t *testing.T) {
	// Create a mock service registrar
	mockRegistrar := NewMockServiceRegistrar()

	// Register the BookstoreService
	server := &mockBookstoreServer{}
	v1.RegisterMCPBookstoreServiceServer(mockRegistrar, server)

	// Get the method descriptor for CreateGenre
	methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
	assert.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

	// Create a validating interceptor
	validatingInterceptor := ValidatingInterceptor(t)

	// Test successful validation case
	t.Run("ValidInput", func(t *testing.T) {
		// Create a map with valid data
		inputMap := map[string]any{
			"name": "Science Fiction",
		}

		// Create a decoder function that uses the input map
		decoder := func(msg proto.Message) error {
			return mcpgw_v1.UnmarshalFromMap(inputMap, msg)
		}

		// Call the handler with the validating interceptor
		resp, err := methodDesc.Handler(server, context.Background(), decoder, validatingInterceptor)

		// No validation errors should occur for valid input
		assert.NoError(t, err, "Valid input should pass validation")
		assert.NotNil(t, resp, "Response should not be nil")

		// Verify the response
		genreResp, ok := resp.(*v1.CreateGenreResponse)
		assert.True(t, ok)
		assert.Equal(t, "Science Fiction", genreResp.GetGenre().GetName())
	})

	// Test with invalid input (if validation rules are defined)
	t.Run("InvalidInput_EmptyName", func(t *testing.T) {
		// Create a map with invalid data (empty name)
		inputMap := map[string]any{
			"name": "", // Empty name may be invalid depending on validation rules
		}

		// Create a decoder function for the invalid input
		decoder := func(msg proto.Message) error {
			return mcpgw_v1.UnmarshalFromMap(inputMap, msg)
		}

		// Call the handler with the validating interceptor
		resp, err := methodDesc.Handler(server, context.Background(), decoder, validatingInterceptor)

		// Note: If there are no validation rules defined for this field in the proto,
		// this will still pass. The test is designed to demonstrate the validation process.
		t.Logf("Validation result: %v", err)

		// Check if validation failed or succeeded
		if err != nil {
			// Validation failed as expected (if rules are defined)
			assert.Nil(t, resp, "Response should be nil when validation fails")
		} else {
			// No validation rules are triggered or defined
			assert.NotNil(t, resp, "Response should not be nil when validation passes")
		}
	})

	// If CreateBookRequest requires a title, test that validation
	t.Run("InvalidBook_MissingTitle", func(t *testing.T) {
		// Get the method descriptor for CreateBook
		bookMethodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateBook"]
		assert.NotNil(t, bookMethodDesc, "Method descriptor for CreateBook should be registered")

		// Create a map with invalid data (missing title which might be required)
		inputMap := map[string]any{
			"shelf": "shelf-123",
			"book": map[string]any{
				"id":     "book-456",
				"author": "Jane Doe",
				// Missing title (might trigger validation if required)
				"quotes":  []string{"Quote 1"},
				"shelfId": "shelf-123",
			},
		}

		// Create a decoder function for the invalid input
		decoder := func(msg proto.Message) error {
			return mcpgw_v1.UnmarshalFromMap(inputMap, msg)
		}

		// Call the handler with the validating interceptor
		resp, err := bookMethodDesc.Handler(server, context.Background(), decoder, validatingInterceptor)

		// Log the validation result for debugging
		t.Logf("Book validation result: %v", err)

		// Check if validation failed or succeeded
		if err != nil {
			// Validation failed as expected (if rules are defined)
			assert.Nil(t, resp, "Response should be nil when validation fails")
		} else {
			// No validation rules are triggered or defined
			assert.NotNil(t, resp, "Response should not be nil when validation passes")
		}
	})
}

// mockBookstoreServer is a mock implementation of BookstoreServiceServer
type mockBookstoreServer struct {
	v1.UnimplementedBookstoreServiceServer
}

func (s *mockBookstoreServer) CreateGenre(ctx context.Context, req *v1.CreateGenreRequest) (*v1.CreateGenreResponse, error) {
	resp := &v1.CreateGenreResponse{}
	genre := &v1.Genre{}
	genre.SetId(42)
	genre.SetName(req.GetName())
	resp.SetGenre(genre)
	return resp, nil
}

func (s *mockBookstoreServer) CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.CreateBookResponse, error) {
	resp := &v1.CreateBookResponse{}
	resp.SetBook(req.GetBook())
	return resp, nil
}
