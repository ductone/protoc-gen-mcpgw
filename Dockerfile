FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/protoc-gen-mcpgw"]
COPY protoc-gen-mcpgw /