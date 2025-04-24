// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
)

// Embed swagger spec.
func GenerateEmbeddedSpec(t *template.Template, importMapping importMap, swagger *openapi3.T) (string, error) {
	var buf bytes.Buffer

	buf.Reset()
	w := bufio.NewWriter(&buf)
	err := t.ExecuteTemplate(w, "embedded.tmpl", struct {
		ImportMapping importMap
	}{ImportMapping: importMapping})
	if err != nil {
		return "", fmt.Errorf("error generating embedded spec: %s", err)
	}
	err = w.Flush()
	if err != nil {
		return "", fmt.Errorf("error flushing output buffer for embedded spec: %s", err)
	}
	return buf.String(), nil
}
