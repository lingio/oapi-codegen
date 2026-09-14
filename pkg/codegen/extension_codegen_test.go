package codegen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const extensionsTestSpec = `
openapi: "3.0.0"
info:
  version: 1.0.0
  title: extensions
paths:
  /things:
    get:
      operationId: getThings
      parameters:
        - name: cursor
          in: query
          required: false
          schema:
            type: string
            x-go-type-skip-optional-pointer: true
      responses:
        '200':
          description: things
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Thing'
components:
  schemas:
    Thing:
      type: object
      properties:
        id:
          type: integer
          x-go-type: uint64
        name:
          type: string
          x-omitempty: false
        cursor:
          type: string
          x-go-type-skip-optional-pointer: true
`

// The values of OpenAPI extensions are handed to us already decoded, so make
// sure we still act on them. Getting this wrong silently changes the generated
// field types and JSON tags rather than failing the build.
func TestGenerateHonoursExtensions(t *testing.T) {
	swagger, err := openapi3.NewLoader().LoadFromData([]byte(extensionsTestSpec))
	require.NoError(t, err)

	code, err := Generate(swagger, "extensions", "", Options{
		GenerateTypes:      true,
		GenerateEchoServer: true,
	})
	require.NoError(t, err)

	// The generated code is gofmt'ed, so fields are padded into columns.
	fields := strings.Join(regexp.MustCompile(`[ \t]+`).Split(code, -1), " ")

	// x-go-type overrides the Go type we'd derive from the schema.
	assert.Contains(t, fields, "ID *uint64 `json:\"id,omitempty\"`")
	// x-omitempty: false keeps the field in the marshalled JSON.
	assert.Contains(t, fields, "Name *string `json:\"name\"`")
	// x-go-type-skip-optional-pointer drops the pointer on an optional field,
	// both on schemas and on query parameters.
	assert.Contains(t, fields, "Cursor string `json:\"cursor,omitempty\"`")
	oneLine := strings.Join(regexp.MustCompile(`\s+`).Split(code, -1), " ")
	assert.Contains(t, oneLine, "type GetThingsParams struct { Cursor string")
}
