package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/ductone/protoc-gen-mcpgw/internal/mcpgw"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	fmt.Fprintf(os.Stderr, "XXXXX Starting protoc-gen-mcpgw\n")
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	slog.Info("YYYYYStarting protoc-gen-mcpgw")

	minEdition := int32(descriptorpb.Edition_EDITION_2023)
	maxEdition := int32(descriptorpb.Edition_EDITION_2023)
	features := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL | pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	options := []pgs.InitOption{
		pgs.DebugEnv("DEBUG_PROTOC_GEN_MCPGW"),
		pgs.SupportedFeatures(&features),
		pgs.MinimumEdition(&minEdition),
		pgs.MaximumEdition(&maxEdition),
	}

	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG_PROTOC_INPUT")); ok {
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, os.Stdin)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("input.data", buf.Bytes(), 0600)
		if err != nil {
			panic(err)
		}
		options = append(options,
			pgs.ProtocInput(bytes.NewReader(buf.Bytes())),
		)
	}

	if fname := os.Getenv("DEBUG_PROTOC_USE_FILE"); fname != "" {
		data, err := os.ReadFile(fname)
		if err != nil {
			panic(err)
		}
		options = append(options,
			pgs.ProtocInput(bytes.NewReader(data)),
		)
	}

	pgs.Init(options...).
		RegisterModule(mcpgw.New()).
		RegisterPostProcessor(pgsgo.GoImports()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
