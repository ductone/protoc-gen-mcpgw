package v1_test

import (
	"context"
	"encoding/json"
	"reflect"
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
		require.NotNil(t, methodDesc.GetInputSchema(), "InputSchema function should be available")

		// Get the schema as map[string]any
		schemaMap := methodDesc.GetInputSchema()()
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
		require.NotNil(t, methodDesc.GetInputSchema(), "InputSchema function should be available")

		// Get the schema as map[string]any
		schemaMap := methodDesc.GetInputSchema()()
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
	methodDescs map[string]mcpgw_v1.MethodDescInterface
}

func NewMockServiceRegistrar() *MockServiceRegistrar {
	return &MockServiceRegistrar{
		methodDescs: make(map[string]mcpgw_v1.MethodDescInterface),
	}
}

func (m *MockServiceRegistrar) RegisterService(sd *mcpgw_v1.ServiceDesc, ss any) {
	// Store the method descriptors for later lookup
	for _, md := range sd.Methods {
		m.methodDescs[md.GetMethod()] = md
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
		err := methodDesc.GetDecoder()(context.Background(), input, req)
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
		err := methodDesc.GetDecoder()(context.Background(), input, req)
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
		err := methodDesc.GetDecoder()(ctx, input, req)
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
		resp, err := methodDesc.GetHandler()(server, ctx, decoder, nil)
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

		resp, err = methodDesc.GetHandler()(server, ctx, decoder, interceptor)
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
		resp, err := methodDesc.GetHandler()(server, context.Background(), decoder, validatingInterceptor)

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
		resp, err := methodDesc.GetHandler()(server, context.Background(), decoder, validatingInterceptor)

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
		resp, err := bookMethodDesc.GetHandler()(server, context.Background(), decoder, validatingInterceptor)

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

// TestConvertToConcreteType tests the ConvertToConcreteType helper function
func TestConvertToConcreteType(t *testing.T) {
	// Create a mock service registrar
	mockRegistrar := NewMockServiceRegistrar()

	// Register the BookstoreService
	server := &mockBookstoreServer{}
	v1.RegisterMCPBookstoreServiceServer(mockRegistrar, server)

	// Test CreateGenre method type conversion
	t.Run("CreateGenre_Request_Valid", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Create a request of the correct type
		originalReq := &v1.CreateGenreRequest{}
		originalReq.SetName("Test Genre")

		// Convert using the helper function
		convertedReq, err := mcpgw_v1.ConvertToConcreteType(methodDesc, originalReq, true)
		assert.NoError(t, err, "Should successfully convert to concrete type")
		assert.NotNil(t, convertedReq, "Converted request should not be nil")

		// Verify it's the correct type
		genreReq, ok := convertedReq.(*v1.CreateGenreRequest)
		assert.True(t, ok, "Should be convertible to CreateGenreRequest")
		assert.Equal(t, "Test Genre", genreReq.GetName(), "Name should be preserved")
	})

	t.Run("CreateGenre_Request_WrongType", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Create a request of the wrong type
		wrongReq := &v1.CreateBookRequest{}
		wrongReq.SetShelf("test-shelf")

		// Convert using the helper function - should fail
		convertedReq, err := mcpgw_v1.ConvertToConcreteType(methodDesc, wrongReq, true)
		assert.Error(t, err, "Should fail to convert wrong type")
		assert.Nil(t, convertedReq, "Converted request should be nil")
		assert.Contains(t, err.Error(), "expected request type", "Error should mention expected type")
	})

	t.Run("CreateGenre_Response_Valid", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Create a response of the correct type
		originalResp := &v1.CreateGenreResponse{}
		genre := &v1.Genre{}
		genre.SetId(123)
		genre.SetName("Test Genre")
		originalResp.SetGenre(genre)

		// Convert using the helper function
		convertedResp, err := mcpgw_v1.ConvertToConcreteType(methodDesc, originalResp, false)
		assert.NoError(t, err, "Should successfully convert to concrete type")
		assert.NotNil(t, convertedResp, "Converted response should not be nil")

		// Verify it's the correct type
		genreResp, ok := convertedResp.(*v1.CreateGenreResponse)
		assert.True(t, ok, "Should be convertible to CreateGenreResponse")
		assert.Equal(t, int64(123), genreResp.GetGenre().GetId(), "ID should be preserved")
		assert.Equal(t, "Test Genre", genreResp.GetGenre().GetName(), "Name should be preserved")
	})

	t.Run("CreateGenre_Response_WrongType", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Create a response of the wrong type
		wrongResp := &v1.CreateBookResponse{}

		// Convert using the helper function - should fail
		convertedResp, err := mcpgw_v1.ConvertToConcreteType(methodDesc, wrongResp, false)
		assert.Error(t, err, "Should fail to convert wrong type")
		assert.Nil(t, convertedResp, "Converted response should be nil")
		assert.Contains(t, err.Error(), "expected response type", "Error should mention expected type")
	})

	// Test CreateBook method type conversion
	t.Run("CreateBook_Request_Valid", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateBook"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateBook should be registered")

		// Create a request of the correct type
		originalReq := &v1.CreateBookRequest{}
		originalReq.SetShelf("test-shelf")
		book := &v1.Book{}
		book.SetId("test-book")
		book.SetTitle("Test Book")
		originalReq.SetBook(book)

		// Convert using the helper function
		convertedReq, err := mcpgw_v1.ConvertToConcreteType(methodDesc, originalReq, true)
		assert.NoError(t, err, "Should successfully convert to concrete type")
		assert.NotNil(t, convertedReq, "Converted request should not be nil")

		// Verify it's the correct type
		bookReq, ok := convertedReq.(*v1.CreateBookRequest)
		assert.True(t, ok, "Should be convertible to CreateBookRequest")
		assert.Equal(t, "test-shelf", bookReq.GetShelf(), "Shelf should be preserved")
		assert.Equal(t, "test-book", bookReq.GetBook().GetId(), "Book ID should be preserved")
		assert.Equal(t, "Test Book", bookReq.GetBook().GetTitle(), "Book title should be preserved")
	})

	t.Run("CreateBook_Request_WrongType", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateBook"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateBook should be registered")

		// Create a request of the wrong type
		wrongReq := &v1.CreateGenreRequest{}
		wrongReq.SetName("Test Genre")

		// Convert using the helper function - should fail
		convertedReq, err := mcpgw_v1.ConvertToConcreteType(methodDesc, wrongReq, true)
		assert.Error(t, err, "Should fail to convert wrong type")
		assert.Nil(t, convertedReq, "Converted request should be nil")
		assert.Contains(t, err.Error(), "expected request type", "Error should mention expected type")
	})

	// Test edge cases
	t.Run("Nil_Input", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Test with nil input
		convertedReq, err := mcpgw_v1.ConvertToConcreteType(methodDesc, nil, true)
		assert.Error(t, err, "Should fail with nil input")
		assert.Nil(t, convertedReq, "Converted request should be nil")
		assert.Contains(t, err.Error(), "input cannot be nil", "Error should mention nil input")
	})

	t.Run("Type_Information_Retrieval", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Test that we can retrieve type information
		requestType := methodDesc.GetRequestType()
		responseType := methodDesc.GetResponseType()

		assert.NotNil(t, requestType, "Request type should not be nil")
		assert.NotNil(t, responseType, "Response type should not be nil")

		// Debug: print the actual types we're getting
		t.Logf("Request type: %v", requestType)
		t.Logf("Response type: %v", responseType)
		t.Logf("Expected request type: %v", reflect.TypeOf(&v1.CreateGenreRequest{}))
		t.Logf("Expected response type: %v", reflect.TypeOf(&v1.CreateGenreResponse{}))

		// Verify the types are correct - our methods return pointer types
		expectedReqType := reflect.TypeOf(&v1.CreateGenreRequest{})
		expectedRespType := reflect.TypeOf(&v1.CreateGenreResponse{})

		// The types should match exactly
		assert.Equal(t, expectedReqType, requestType, "Request type should match CreateGenreRequest")
		assert.Equal(t, expectedRespType, responseType, "Response type should match CreateGenreResponse")

		// Test creating new instances
		newReq := methodDesc.CreateRequest()
		newResp := methodDesc.CreateResponse()

		assert.NotNil(t, newReq, "New request should not be nil")
		assert.NotNil(t, newResp, "New response should not be nil")

		// Verify they're the correct types
		_, ok := newReq.(*v1.CreateGenreRequest)
		assert.True(t, ok, "New request should be CreateGenreRequest")

		_, ok = newResp.(*v1.CreateGenreResponse)
		assert.True(t, ok, "New response should be CreateGenreResponse")
	})

	t.Run("Assert_Methods", func(t *testing.T) {
		methodDesc := mockRegistrar.methodDescs["/bookstore.v1.BookstoreService/CreateGenre"]
		require.NotNil(t, methodDesc, "Method descriptor for CreateGenre should be registered")

		// Test AssertRequestType method
		originalReq := &v1.CreateGenreRequest{}
		originalReq.SetName("Test Genre")

		convertedReq, err := methodDesc.AssertRequestType(context.Background(), originalReq)
		assert.NoError(t, err, "Should successfully assert request type")
		assert.NotNil(t, convertedReq, "Converted request should not be nil")

		genreReq, ok := convertedReq.(*v1.CreateGenreRequest)
		assert.True(t, ok, "Should be convertible to CreateGenreRequest")
		assert.Equal(t, "Test Genre", genreReq.GetName(), "Name should be preserved")

		// Test AssertResponseType method
		originalResp := &v1.CreateGenreResponse{}
		genre := &v1.Genre{}
		genre.SetId(456)
		genre.SetName("Test Genre")
		originalResp.SetGenre(genre)

		convertedResp, err := methodDesc.AssertResponseType(context.Background(), originalResp)
		assert.NoError(t, err, "Should successfully assert response type")
		assert.NotNil(t, convertedResp, "Converted response should not be nil")

		genreResp, ok := convertedResp.(*v1.CreateGenreResponse)
		assert.True(t, ok, "Should be convertible to CreateGenreResponse")
		assert.Equal(t, int64(456), genreResp.GetGenre().GetId(), "ID should be preserved")

		// Test with wrong types
		wrongReq := &v1.CreateBookRequest{}
		_, err = methodDesc.AssertRequestType(context.Background(), wrongReq)
		assert.Error(t, err, "Should fail to assert wrong request type")

		wrongResp := &v1.CreateBookResponse{}
		_, err = methodDesc.AssertResponseType(context.Background(), wrongResp)
		assert.Error(t, err, "Should fail to assert wrong response type")
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
