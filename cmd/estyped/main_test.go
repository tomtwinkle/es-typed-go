package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestToPascalCase(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path string
		want string
	}{
		"simple":              {path: "status", want: "Status"},
		"dotted":              {path: "items.color", want: "ItemsColor"},
		"multi_field":         {path: "title.keyword", want: "TitleKeyword"},
		"underscore":          {path: "field_name", want: "FieldName"},
		"deep_nested":         {path: "items.color.value", want: "ItemsColorValue"},
		"underscore_and_dot":  {path: "field_name.sub_field", want: "FieldNameSubField"},
		"double_underscore":   {path: "field__name", want: "FieldName"},
		"leading_underscore":  {path: "_field", want: "Field"},
		"trailing_underscore": {path: "field_", want: "Field"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := toPascalCase(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToStructFieldName(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		path string
		want string
	}{
		"simple":              {path: "status", want: "Status"},
		"dotted":              {path: "items.color", want: "Items_Color"},
		"multi_field":         {path: "title.keyword", want: "Title_Keyword"},
		"underscore":          {path: "field_name", want: "FieldName"},
		"deep_nested":         {path: "items.color.value", want: "Items_Color_Value"},
		"underscore_and_dot":  {path: "field_name.sub_field", want: "FieldName_SubField"},
		"double_underscore":   {path: "field__name", want: "FieldName"},
		"leading_underscore":  {path: "_field", want: "Field"},
		"trailing_underscore": {path: "field_", want: "Field"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := toStructFieldName(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFieldType(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		typ  string
		want string
	}{
		"keyword":                  {typ: "keyword", want: "keyword"},
		"text":                     {typ: "text", want: "text"},
		"nested":                   {typ: "nested", want: "nested"},
		"empty_defaults_to_object": {typ: "", want: "object"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fieldType(tt.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

// generateSource executes the template and returns formatted Go source.
func generateSource(t *testing.T, mappingData []byte, pkgName, structName string) []byte {
	t.Helper()

	mapping, err := estype.ParseMapping(mappingData)
	assert.NilError(t, err)

	entries := make([]fieldEntry, 0, len(mapping.Fields))
	for _, f := range mapping.Fields {
		entries = append(entries, fieldEntry{
			ConstName: toPascalCase(f.Path),
			FieldName: toStructFieldName(f.Path),
			Path:      f.Path,
			Type:      fieldType(f.TypeName()),
		})
	}

	td := templateData{
		Package: pkgName,
		Name:    structName,
		Fields:  entries,
	}

	var buf bytes.Buffer
	if structName != "" {
		err = structTemplate.Execute(&buf, td)
	} else {
		err = constTemplate.Execute(&buf, td)
	}
	assert.NilError(t, err)

	formatted, err := format.Source(buf.Bytes())
	assert.NilError(t, err)

	return formatted
}

// TestGenerate_ConstMode runs the estyped code generator in constant mode
// and uses go/ast to verify the correctness of the generated Go source.
func TestGenerate_ConstMode(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	outFile := filepath.Join(dir, "model.go")

	mapping := []byte(`{
		"properties": {
			"status": { "type": "keyword" },
			"title": {
				"type": "text",
				"fields": {
					"keyword": { "type": "keyword" }
				}
			},
			"items": {
				"type": "nested",
				"properties": {
					"color": { "type": "keyword" }
				}
			}
		}
	}`)

	generated := generateSource(t, mapping, "model", "")
	err := os.WriteFile(outFile, generated, 0o644)
	assert.NilError(t, err)

	// Parse the generated file with go/ast.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	// Verify package name.
	assert.Equal(t, "model", f.Name.Name)

	// Verify that "estype" is imported.
	assertImport(t, f, `"github.com/tomtwinkle/es-typed-go/estype"`)

	// Collect all top-level constant declarations.
	consts := collectConstDecls(f)

	// Expected constants: items, items.color, status, title, title.keyword
	expectedConsts := map[string]string{
		"FieldItems":        "items",
		"FieldItemsColor":   "items.color",
		"FieldStatus":       "status",
		"FieldTitle":        "title",
		"FieldTitleKeyword": "title.keyword",
	}
	assert.Assert(t, len(consts) == len(expectedConsts), "expected %d constants, got %d", len(expectedConsts), len(consts))
	for name, wantValue := range expectedConsts {
		gotValue, ok := consts[name]
		assert.Assert(t, ok, "missing constant %q", name)
		assert.Equal(t, wantValue, gotValue, "constant %s has wrong value", name)
	}
}

// TestGenerate_StructMode runs the estyped code generator in struct mode
// and uses go/ast to verify the correctness of the generated Go source.
func TestGenerate_StructMode(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	outFile := filepath.Join(dir, "model.go")

	mapping := []byte(`{
		"properties": {
			"status": { "type": "keyword" },
			"title": {
				"type": "text",
				"fields": {
					"keyword": { "type": "keyword" }
				}
			},
			"items": {
				"type": "nested",
				"properties": {
					"color": { "type": "keyword" }
				}
			}
		}
	}`)

	generated := generateSource(t, mapping, "model", "Sample")
	err := os.WriteFile(outFile, generated, 0o644)
	assert.NilError(t, err)

	// Parse the generated file with go/ast.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	// Verify package name.
	assert.Equal(t, "model", f.Name.Name)

	// Verify that "estype" is imported.
	assertImport(t, f, `"github.com/tomtwinkle/es-typed-go/estype"`)

	// Find the var declaration for "Sample".
	structFields := collectStructVarFields(f, "Sample")
	assert.Assert(t, len(structFields) > 0, "expected struct var 'Sample' with fields")

	expectedFields := map[string]string{
		"Items":         "items",
		"Items_Color":   "items.color",
		"Status":        "status",
		"Title":         "title",
		"Title_Keyword": "title.keyword",
	}
	assert.Assert(t, len(structFields) == len(expectedFields), "expected %d struct fields, got %d", len(expectedFields), len(structFields))
	for name, wantValue := range expectedFields {
		gotValue, ok := structFields[name]
		assert.Assert(t, ok, "missing struct field %q", name)
		assert.Equal(t, wantValue, gotValue, "struct field %s has wrong value", name)
	}
}

// TestGenerate_NestedFields verifies deeply nested field generation using go/ast.
func TestGenerate_NestedFields(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	outFile := filepath.Join(dir, "model.go")

	mapping := []byte(`{
		"mappings": {
			"properties": {
				"items": {
					"type": "nested",
					"properties": {
						"name": { "type": "text" },
						"tags": {
							"type": "nested",
							"properties": {
								"value": { "type": "keyword" }
							}
						}
					}
				}
			}
		}
	}`)

	generated := generateSource(t, mapping, "model", "")
	err := os.WriteFile(outFile, generated, 0o644)
	assert.NilError(t, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	consts := collectConstDecls(f)
	expectedConsts := map[string]string{
		"FieldItems":          "items",
		"FieldItemsName":      "items.name",
		"FieldItemsTags":      "items.tags",
		"FieldItemsTagsValue": "items.tags.value",
	}
	assert.Assert(t, len(consts) == len(expectedConsts), "expected %d constants, got %d", len(expectedConsts), len(consts))
	for name, wantValue := range expectedConsts {
		gotValue, ok := consts[name]
		assert.Assert(t, ok, "missing constant %q", name)
		assert.Equal(t, wantValue, gotValue, "constant %s has wrong value", name)
	}
}

// TestGenerate_ObjectField verifies that fields with empty type get "object" in comments.
func TestGenerate_ObjectField(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	outFile := filepath.Join(dir, "model.go")

	mapping := []byte(`{
		"properties": {
			"data": {
				"properties": {
					"value": { "type": "keyword" }
				}
			}
		}
	}`)

	generated := generateSource(t, mapping, "model", "")
	err := os.WriteFile(outFile, generated, 0o644)
	assert.NilError(t, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	consts := collectConstDecls(f)
	// "data" has no type → object, "data.value" → keyword
	_, ok := consts["FieldData"]
	assert.Assert(t, ok, "missing constant FieldData")
	_, ok = consts["FieldDataValue"]
	assert.Assert(t, ok, "missing constant FieldDataValue")
}

// writeGoFile writes src as a .go file named name in dir and returns its path.
func writeGoFile(t *testing.T, dir, name, src string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	assert.NilError(t, os.WriteFile(path, []byte(src), 0o644))
	return path
}

// TestParseGoStruct_FlatFields verifies that a flat struct with json tags produces
// the expected field entries.
func TestParseGoStruct_FlatFields(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Status   string `+"`"+`json:"status"`+"`"+`
	Title    string `+"`"+`json:"title"`+"`"+`
	Price    int    `+"`"+`json:"price"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.ConstName] = e.Path
	}
	assert.Equal(t, 3, len(got))
	assert.Equal(t, "price", got["Price"])
	assert.Equal(t, "status", got["Status"])
	assert.Equal(t, "title", got["Title"])
}

// TestParseGoStruct_NestedStruct verifies that a field whose type is another
// struct defined in the same package is expanded recursively.
func TestParseGoStruct_NestedStruct(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Items  Item   `+"`"+`json:"items"`+"`"+`
}

type Item struct {
	Name  string `+"`"+`json:"name"`+"`"+`
	Value int    `+"`"+`json:"value"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.ConstName] = e.Path
	}
	assert.Equal(t, 4, len(got))
	assert.Equal(t, "items", got["Items"])
	assert.Equal(t, "items.name", got["ItemsName"])
	assert.Equal(t, "items.value", got["ItemsValue"])
	assert.Equal(t, "status", got["Status"])
}

// TestParseGoStruct_SliceOfStruct verifies that a slice-of-struct field uses
// "nested" as the ES type.
func TestParseGoStruct_SliceOfStruct(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Tags []Tag `+"`"+`json:"tags"`+"`"+`
}

type Tag struct {
	Value string `+"`"+`json:"value"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	gotTypes := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.ConstName] = e.Path
		gotTypes[e.ConstName] = e.Type
	}
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "tags", got["Tags"])
	assert.Equal(t, "nested", gotTypes["Tags"])
	assert.Equal(t, "tags.value", got["TagsValue"])
}

// TestParseGoStruct_SkipMinusTag verifies that fields with json:"-" are omitted.
func TestParseGoStruct_SkipMinusTag(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Status  string `+"`"+`json:"status"`+"`"+`
	Ignored string `+"`"+`json:"-"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "status", entries[0].Path)
}

// TestParseGoStruct_OmitemptyTag verifies that omitempty options are handled and
// the field name before the comma is used.
func TestParseGoStruct_OmitemptyTag(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Status string `+"`"+`json:"status,omitempty"`+"`"+`
	Title  string `+"`"+`json:"title,omitempty"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.ConstName] = e.Path
	}
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "status", got["Status"])
	assert.Equal(t, "title", got["Title"])
}

// TestParseGoStruct_EmbeddedStruct verifies that anonymous embedded struct fields
// are inlined at the same nesting level.
func TestParseGoStruct_EmbeddedStruct(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Base struct {
	ID string `+"`"+`json:"id"`+"`"+`
}

type Document struct {
	Base
	Status string `+"`"+`json:"status"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.ConstName] = e.Path
	}
	assert.Equal(t, 2, len(got))
	assert.Equal(t, "id", got["Id"])
	assert.Equal(t, "status", got["Status"])
}

// TestParseGoStruct_TypeNotFound verifies that a helpful error is returned when
// the requested type does not exist.
func TestParseGoStruct_TypeNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
}
`)
	_, _, err := parseGoStruct(dir, "Missing")
	assert.ErrorContains(t, err, `"Missing" not found`)
}

// TestParseGoStruct_ConstModeOutput verifies that struct-based parsing produces
// valid Go source in constant output mode.
func TestParseGoStruct_ConstModeOutput(t *testing.T) {
	t.Parallel()
	srcDir := t.TempDir()
	writeGoFile(t, srcDir, "doc.go", `package model

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Items  []Item `+"`"+`json:"items"`+"`"+`
}

type Item struct {
	Name string `+"`"+`json:"name"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(srcDir, "Document")
	assert.NilError(t, err)

	td := templateData{Package: "model", Fields: entries}
	var buf bytes.Buffer
	assert.NilError(t, constTemplate.Execute(&buf, td))
	formatted, err := format.Source(buf.Bytes())
	assert.NilError(t, err)

	// Parse the generated source and verify constants.
	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "fields.go")
	assert.NilError(t, os.WriteFile(outFile, formatted, 0o644))

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	assertImport(t, f, `"github.com/tomtwinkle/es-typed-go/estype"`)

	consts := collectConstDecls(f)
	expectedConsts := map[string]string{
		"FieldItems":     "items",
		"FieldItemsName": "items.name",
		"FieldStatus":    "status",
	}
	assert.Assert(t, len(consts) == len(expectedConsts), "expected %d constants, got %d", len(expectedConsts), len(consts))
	for name, wantValue := range expectedConsts {
		gotValue, ok := consts[name]
		assert.Assert(t, ok, "missing constant %q", name)
		assert.Equal(t, wantValue, gotValue, "constant %s has wrong value", name)
	}
}

// TestParseGoStruct_StructModeOutput verifies that struct-based parsing produces
// valid Go source in struct variable output mode.
func TestParseGoStruct_StructModeOutput(t *testing.T) {
	t.Parallel()
	srcDir := t.TempDir()
	writeGoFile(t, srcDir, "doc.go", `package model

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Items  []Item `+"`"+`json:"items"`+"`"+`
}

type Item struct {
	Name string `+"`"+`json:"name"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(srcDir, "Document")
	assert.NilError(t, err)

	td := templateData{Package: "model", Name: "Sample", Fields: entries}
	var buf bytes.Buffer
	assert.NilError(t, structTemplate.Execute(&buf, td))
	formatted, err := format.Source(buf.Bytes())
	assert.NilError(t, err)

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "fields.go")
	assert.NilError(t, os.WriteFile(outFile, formatted, 0o644))

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, outFile, nil, parser.AllErrors)
	assert.NilError(t, err)

	assertImport(t, f, `"github.com/tomtwinkle/es-typed-go/estype"`)

	structFields := collectStructVarFields(f, "Sample")
	expectedFields := map[string]string{
		"Items":      "items",
		"Items_Name": "items.name",
		"Status":     "status",
	}
	assert.Assert(t, len(structFields) == len(expectedFields), "expected %d struct fields, got %d", len(expectedFields), len(structFields))
	for name, wantValue := range expectedFields {
		gotValue, ok := structFields[name]
		assert.Assert(t, ok, "missing struct field %q", name)
		assert.Equal(t, wantValue, gotValue, "struct field %s has wrong value", name)
	}
}

// TestGoFieldJSONName covers struct field JSON-name extraction for tagged,
// untagged, excluded, and unnamed fields.
func TestGoFieldJSONName(t *testing.T) {
	t.Parallel()

	parseField := func(t *testing.T, src string) *ast.Field {
		t.Helper()
		fset := token.NewFileSet()
		expr, err := parser.ParseExprFrom(fset, "", src, 0)
		assert.NilError(t, err)

		st, ok := expr.(*ast.StructType)
		assert.Assert(t, ok)
		assert.Assert(t, len(st.Fields.List) == 1)

		return st.Fields.List[0]
	}

	tests := map[string]struct {
		src  string
		want string
	}{
		"tagged": {
			src:  `struct{ Status string ` + "`json:\"status\"`" + ` }`,
			want: "status",
		},
		"tagged_skip": {
			src:  `struct{ Ignored string ` + "`json:\"-\"`" + ` }`,
			want: "-",
		},
		"tagged_empty_name_falls_back_to_go_name": {
			src:  `struct{ Status string ` + "`json:\",omitempty\"`" + ` }`,
			want: "Status",
		},
		"untagged_named_field": {
			src:  `struct{ Status string }`,
			want: "Status",
		},
		"anonymous_field_without_tag": {
			src:  `struct{ Embedded }`,
			want: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := goFieldJSONName(parseField(t, tt.src))
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestJSONTagKey covers the jsonTagKey helper for all relevant tag formats.
func TestJSONTagKey(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		raw  string
		want string
	}{
		"plain":            {raw: `json:"status"`, want: "status"},
		"omitempty":        {raw: `json:"title,omitempty"`, want: "title"},
		"skip":             {raw: `json:"-"`, want: "-"},
		"no_json_tag":      {raw: `db:"col"`, want: ""},
		"empty_name":       {raw: `json:",omitempty"`, want: ""},
		"unterminated_tag": {raw: `json:"status`, want: ""},
		"json_after_other": {raw: `xml:"status" json:"title,omitempty"`, want: "title"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := jsonTagKey(tt.raw)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestDerefTypeName covers identifier, pointer, slice, array, and unsupported
// expression shapes for derefTypeName.
func TestDerefTypeName(t *testing.T) {
	t.Parallel()

	parseExpr := func(t *testing.T, src string) ast.Expr {
		t.Helper()
		fset := token.NewFileSet()
		expr, err := parser.ParseExprFrom(fset, "", src, 0)
		assert.NilError(t, err)
		return expr
	}

	tests := map[string]struct {
		src  string
		want string
	}{
		"ident":            {src: `Tag`, want: "Tag"},
		"pointer":          {src: `*Tag`, want: "Tag"},
		"slice":            {src: `[]Tag`, want: "Tag"},
		"array":            {src: `[2]Tag`, want: "Tag"},
		"pointer_to_slice": {src: `*[]Tag`, want: "Tag"},
		"selector_expr":    {src: `time.Time`, want: ""},
		"map_type":         {src: `map[string]Tag`, want: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := derefTypeName(parseExpr(t, tt.src))
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestIsSliceExpr covers pointer-unwrapping and distinguishes slices from
// arrays and non-array expressions.
func TestIsSliceExpr(t *testing.T) {
	t.Parallel()

	parseExpr := func(t *testing.T, src string) ast.Expr {
		t.Helper()
		fset := token.NewFileSet()
		expr, err := parser.ParseExprFrom(fset, "", src, 0)
		assert.NilError(t, err)
		return expr
	}

	tests := map[string]struct {
		src  string
		want bool
	}{
		"slice":            {src: `[]Tag`, want: true},
		"pointer_to_slice": {src: `*[]Tag`, want: true},
		"array":            {src: `[2]Tag`, want: false},
		"pointer_to_array": {src: `*[2]Tag`, want: false},
		"pointer":          {src: `*Tag`, want: false},
		"ident":            {src: `Tag`, want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := isSliceExpr(parseExpr(t, tt.src))
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPascalToSnake covers the pascalToSnake helper used to derive ES type names
// from property constructor function name fragments.
func TestPascalToSnake(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input string
		want  string
	}{
		"single_word":  {input: "Text", want: "text"},
		"single_lower": {input: "Keyword", want: "keyword"},
		"two_words":    {input: "DenseVector", want: "dense_vector"},
		"three_words":  {input: "RankFeatures", want: "rank_features"},
		"all_lower":    {input: "nested", want: "nested"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := pascalToSnake(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPropertyCallTypeName covers the propertyCallTypeName helper that derives an
// ES type name from a NewXxxProperty constructor function AST expression.
func TestPropertyCallTypeName(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		src  string // a Go expression that is the Fun part of a CallExpr
		want string
	}{
		"ident_text":    {src: "NewTextProperty", want: "text"},
		"ident_keyword": {src: "NewKeywordProperty", want: "keyword"},
		"ident_nested":  {src: "NewNestedProperty", want: "nested"},
		"ident_dense":   {src: "NewDenseVectorProperty", want: "dense_vector"},
		"not_new":       {src: "MakeTextProperty", want: ""},
		"not_property":  {src: "NewTextField", want: ""},
		"empty_middle":  {src: "NewProperty", want: ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			expr, err := parser.ParseExprFrom(fset, "", tt.src, 0)
			assert.NilError(t, err)
			got := propertyCallTypeName(expr)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestFieldTypeFuncName covers function-name extraction from identifiers,
// selectors, and unsupported expression forms.
func TestFieldTypeFuncName(t *testing.T) {
	t.Parallel()

	parseExpr := func(t *testing.T, src string) ast.Expr {
		t.Helper()
		fset := token.NewFileSet()
		expr, err := parser.ParseExprFrom(fset, "", src, 0)
		assert.NilError(t, err)
		return expr
	}

	tests := map[string]struct {
		src  string
		want string
	}{
		"ident":      {src: `FieldType`, want: "FieldType"},
		"selector":   {src: `estype.FieldType`, want: "FieldType"},
		"call_expr":  {src: `FieldType("keyword")`, want: ""},
		"index_expr": {src: `fn[i]`, want: ""},
		"paren_expr": {src: `(FieldType)`, want: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fieldTypeFuncName(parseExpr(t, tt.src))
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPropertyValueTypeName covers the propertyValueTypeName helper that extracts
// an ES type name from a full Property value AST expression.
// It covers FieldType("...") conversions, NewXxxProperty(...) constructors,
// and the plain string literal fallback.
func TestPropertyValueTypeName(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		src  string // a Go expression used as the Property value
		want string
	}{
		"field_type_qualified":       {src: `estype.FieldType("integer")`, want: "integer"},
		"field_type_unqualified":     {src: `FieldType("keyword")`, want: "keyword"},
		"field_type_empty":           {src: `FieldType("")`, want: ""},
		"field_type_wrong_arg_count": {src: `FieldType("keyword", "text")`, want: ""},
		"field_type_non_string_arg":  {src: `FieldType(42)`, want: ""},
		"constructor_qualified":      {src: `estype.NewTextProperty()`, want: "text"},
		"constructor_unqualified":    {src: `NewKeywordProperty()`, want: "keyword"},
		"constructor_unknown":        {src: `BuildKeywordProperty()`, want: ""},
		"string_literal_fallback":    {src: `"date"`, want: "date"},
		"non_string_literal":         {src: `42`, want: ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			expr, err := parser.ParseExprFrom(fset, "", tt.src, 0)
			assert.NilError(t, err)
			got := propertyValueTypeName(expr)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestParseMappingBody covers statically analyzable and ignored branches in
// parseMappingBody.
func TestParseMappingBody(t *testing.T) {
	t.Parallel()

	parseBody := func(t *testing.T, src string) *ast.BlockStmt {
		t.Helper()
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", "package model\nfunc mapping() any "+src, 0)
		assert.NilError(t, err)

		assert.Assert(t, len(file.Decls) == 1)
		fd, ok := file.Decls[0].(*ast.FuncDecl)
		assert.Assert(t, ok)
		assert.Assert(t, fd.Body != nil)

		return fd.Body
	}

	tests := map[string]struct {
		src  string
		want map[string]string
	}{
		"single_return_with_fields": {
			src: `{
				return estype.Mapping{
					Fields: []estype.MappingField{
						{Path: "status", Property: estype.FieldType("keyword")},
						{Path: "title", Property: estype.NewTextProperty()},
					},
				}
			}`,
			want: map[string]string{
				"status": "keyword",
				"title":  "text",
			},
		},
		"ignores_non_return_statements": {
			src: `{
				x := 1
				_ = x
				return estype.Mapping{
					Fields: []estype.MappingField{
						{Path: "status", Property: estype.FieldType("keyword")},
					},
				}
			}`,
			want: map[string]string{
				"status": "keyword",
			},
		},
		"ignores_multi_result_return": {
			src: `{
				return estype.Mapping{}, nil
			}`,
			want: map[string]string{},
		},
		"ignores_non_composite_return": {
			src: `{
				mapping := estype.Mapping{}
				return mapping
			}`,
			want: map[string]string{},
		},
		"ignores_non_fields_keys": {
			src: `{
				return estype.Mapping{
					Name: "ignored",
				}
			}`,
			want: map[string]string{},
		},
		"ignores_non_composite_fields_value": {
			src: `{
				fields := []estype.MappingField{}
				_ = fields
				return estype.Mapping{
					Fields: fields,
				}
			}`,
			want: map[string]string{},
		},
		"ignores_invalid_field_elements": {
			src: `{
				return estype.Mapping{
					Fields: []estype.MappingField{
						123,
						{Path: "status", Property: estype.FieldType("keyword")},
						{Path: dynamicPath(), Property: estype.FieldType("text")},
						{Path: "title", Property: propertyExpr()},
						{Path: "price"},
						{Property: estype.FieldType("integer")},
					},
				}
			}`,
			want: map[string]string{
				"status": "keyword",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := make(map[string]string)
			parseMappingBody(parseBody(t, tt.src), got)
			assert.Equal(t, len(tt.want), len(got))
			for path, wantType := range tt.want {
				gotType, ok := got[path]
				assert.Assert(t, ok, "missing path %q", path)
				assert.Equal(t, wantType, gotType)
			}
		})
	}
}

// TestExtractMappingMethod verifies that the extractMappingMethod function
// correctly parses the return statement of a Mapping() method and returns
// the expected path→type map.
func TestExtractMappingMethod(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src      string
		typeName string
		want     map[string]string
	}{
		"value_receiver": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
	Title  string ` + "`" + `json:"title"` + "`" + `
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "title", Property: estype.FieldType("text")},
		},
	}
}
`,
			typeName: "Document",
			want:     map[string]string{"status": "keyword", "title": "text"},
		},
		"pointer_receiver": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
}

func (d *Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
		},
	}
}
`,
			typeName: "Document",
			want:     map[string]string{"status": "keyword"},
		},
		"no_method": {
			src: `package model

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
}
`,
			typeName: "Document",
			want:     map[string]string{},
		},
		"wrong_receiver_type": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}
type Other struct{}

func (Other) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
		},
	}
}
`,
			typeName: "Document",
			want:     map[string]string{},
		},
		"nested_field_paths": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
	Items  []Item ` + "`" + `json:"items"` + "`" + `
}

type Item struct {
	Name string ` + "`" + `json:"name"` + "`" + `
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "items", Property: estype.FieldType("nested")},
			{Path: "items.name", Property: estype.FieldType("text")},
		},
	}
}
`,
			typeName: "Document",
			want: map[string]string{
				"status":     "keyword",
				"items":      "nested",
				"items.name": "text",
			},
		},
		"typed_property_constructors": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
	Title  string ` + "`" + `json:"title"` + "`" + `
	Price  int    ` + "`" + `json:"price"` + "`" + `
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.NewKeywordProperty()},
			{Path: "title",  Property: estype.NewTextProperty(estype.WithSearchAnalyzer(estype.Analyzer("my_analyzer")))},
			{Path: "price",  Property: estype.FieldType("integer")},
		},
	}
}
`,
			typeName: "Document",
			want: map[string]string{
				"status": "keyword",
				"title":  "text",
				"price":  "integer",
			},
		},
		"unqualified_constructors": {
			src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string ` + "`" + `json:"status"` + "`" + `
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: NewKeywordProperty()},
		},
	}
}
`,
			typeName: "Document",
			want:     map[string]string{"status": "keyword"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writeGoFile(t, dir, "doc.go", tt.src)

			pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, 0)
			assert.NilError(t, err)

			got := extractMappingMethod(pkgs, tt.typeName)
			assert.Equal(t, len(tt.want), len(got))
			for path, wantType := range tt.want {
				gotType, ok := got[path]
				assert.Assert(t, ok, "missing path %q", path)
				assert.Equal(t, wantType, gotType)
			}
		})
	}
}

// TestParseGoStruct_WithMapping verifies that when a struct has a Mapping()
// method, the generated field entries use the types declared in that method
// instead of "unknown".
func TestParseGoStruct_WithMapping(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Title  string `+"`"+`json:"title"`+"`"+`
	Price  int    `+"`"+`json:"price"`+"`"+`
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "title", Property: estype.FieldType("text")},
			{Path: "price", Property: estype.FieldType("integer")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Path] = e.Type
	}
	assert.Equal(t, 3, len(got))
	assert.Equal(t, "keyword", got["status"])
	assert.Equal(t, "text", got["title"])
	assert.Equal(t, "integer", got["price"])
}

// TestParseGoStruct_WithMapping_PointerReceiver verifies that Mapping() with
// a pointer receiver is also detected and its types are applied.
func TestParseGoStruct_WithMapping_PointerReceiver(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
}

func (d *Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "status", entries[0].Path)
	assert.Equal(t, "keyword", entries[0].Type)
}

// TestParseGoStruct_TagFallbacksAndPointers verifies fallback-to-Go-name,
// pointer recursion, anonymous pointer embedding, arrays, and unsupported
// external struct types.
func TestParseGoStruct_TagFallbacksAndPointers(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "time"

type Base struct {
	ID string `+"`"+`json:"id"`+"`"+`
}

type Item struct {
	Name string `+"`"+`json:"name"`+"`"+`
}

type Document struct {
	*Base
	Status    string    `+"`"+`json:",omitempty"`+"`"+`
	Item      *Item     `+"`"+`json:"item"`+"`"+`
	ItemArray [2]Item   `+"`"+`json:"item_array"`+"`"+`
	CreatedAt time.Time `+"`"+`json:"created_at"`+"`"+`
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	gotTypes := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Path] = e.FieldName
		gotTypes[e.Path] = e.Type
	}

	assert.Equal(t, 7, len(entries))
	assert.Equal(t, "Id", got["id"])
	assert.Equal(t, "Status", got["Status"])
	assert.Equal(t, "Item", got["item"])
	assert.Equal(t, "Item_Name", got["item.name"])
	assert.Equal(t, "ItemArray", got["item_array"])
	assert.Equal(t, "ItemArray_Name", got["item_array.name"])
	assert.Equal(t, "unknown", gotTypes["created_at"])
	assert.Equal(t, "object", gotTypes["item"])
	assert.Equal(t, "object", gotTypes["item_array"])
	assert.Equal(t, "unknown", gotTypes["item.name"])
	assert.Equal(t, "unknown", gotTypes["item_array.name"])
}

// TestParseGoStruct_WithMapping_PartialOverride verifies that only fields
// present in Mapping() have their type overridden; fields absent from the
// method fall back to "unknown" or their struct-derived default.
func TestParseGoStruct_WithMapping_PartialOverride(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status  string `+"`"+`json:"status"`+"`"+`
	Title   string `+"`"+`json:"title"`+"`"+`
	Enabled bool   `+"`"+`json:"enabled"`+"`"+`
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Path] = e.Type
	}
	assert.Equal(t, 3, len(got))
	assert.Equal(t, "keyword", got["status"])
	assert.Equal(t, "unknown", got["title"])
	assert.Equal(t, "unknown", got["enabled"])
}

// TestParseGoStruct_WithMapping_NestedFields verifies that Mapping() can
// supply types for both parent struct fields and their nested children.
func TestParseGoStruct_WithMapping_NestedFields(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Items  []Item `+"`"+`json:"items"`+"`"+`
}

type Item struct {
	Name  string `+"`"+`json:"name"`+"`"+`
	Value int    `+"`"+`json:"value"`+"`"+`
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "items", Property: estype.FieldType("nested")},
			{Path: "items.name", Property: estype.FieldType("text")},
			{Path: "items.value", Property: estype.FieldType("integer")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Path] = e.Type
	}
	assert.Equal(t, 4, len(got))
	assert.Equal(t, "keyword", got["status"])
	assert.Equal(t, "nested", got["items"])
	assert.Equal(t, "text", got["items.name"])
	assert.Equal(t, "integer", got["items.value"])
}

// TestParseGoStruct_WithMapping_ConstOutput verifies that struct-based parsing
// with a Mapping() method produces correct types in the generated constant output.
func TestParseGoStruct_WithMapping_ConstOutput(t *testing.T) {
	t.Parallel()
	srcDir := t.TempDir()
	writeGoFile(t, srcDir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Title  string `+"`"+`json:"title"`+"`"+`
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "title", Property: estype.FieldType("text")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(srcDir, "Document")
	assert.NilError(t, err)

	// Verify that types are correctly read from Mapping().
	gotTypes := make(map[string]string, len(entries))
	for _, e := range entries {
		gotTypes[e.Path] = e.Type
	}
	assert.Equal(t, "keyword", gotTypes["status"])
	assert.Equal(t, "text", gotTypes["title"])

	// Also verify the generated Go source is valid.
	td := templateData{Package: "model", Fields: entries}
	var buf bytes.Buffer
	assert.NilError(t, constTemplate.Execute(&buf, td))
	_, err = format.Source(buf.Bytes())
	assert.NilError(t, err)
}

// TestParseGoStruct_WithMapping_TypedProperty verifies that Mapping() using
// typed property constructors (NewTextProperty, NewKeywordProperty) instead of
// plain strings produces the correct ES type in the generated entries.
func TestParseGoStruct_WithMapping_TypedProperty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeGoFile(t, dir, "doc.go", `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct {
	Status string `+"`"+`json:"status"`+"`"+`
	Title  string `+"`"+`json:"title"`+"`"+`
	Price  int    `+"`"+`json:"price"`+"`"+`
}

func (Document) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.NewKeywordProperty(estype.WithIgnoreAbove(256))},
			{Path: "title",  Property: estype.NewTextProperty(
				estype.WithSearchAnalyzer(estype.Analyzer("my_search_analyzer")),
				estype.WithIndexAnalyzer(estype.Analyzer("my_index_analyzer")),
			)},
			{Path: "price",  Property: estype.FieldType("integer")},
		},
	}
}
`)
	entries, _, err := parseGoStruct(dir, "Document")
	assert.NilError(t, err)

	got := make(map[string]string, len(entries))
	for _, e := range entries {
		got[e.Path] = e.Type
	}
	assert.Equal(t, 3, len(got))
	assert.Equal(t, "keyword", got["status"])
	assert.Equal(t, "text", got["title"])
	assert.Equal(t, "integer", got["price"])
}

// --- ast helpers ---

// assertImport verifies that an import path exists in the parsed file.
func assertImport(t *testing.T, f *ast.File, importPath string) {
	t.Helper()
	for _, imp := range f.Imports {
		if imp.Path.Value == importPath {
			return
		}
	}
	t.Errorf("import %s not found", importPath)
}

// collectConstDecls collects all top-level constant declarations as name→value pairs.
func collectConstDecls(f *ast.File) map[string]string {
	consts := make(map[string]string)
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) == 0 || len(vs.Values) == 0 {
				continue
			}
			bl, ok := vs.Values[0].(*ast.BasicLit)
			if !ok {
				continue
			}
			// Strip surrounding quotes from the string literal.
			val := bl.Value
			if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			}
			consts[vs.Names[0].Name] = val
		}
	}
	return consts
}

// collectStructVarFields extracts struct field name→value pairs from a var declaration
// of the form: var Name = struct{ ... }{ ... }
func collectStructVarFields(f *ast.File, varName string) map[string]string {
	fields := make(map[string]string)
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) == 0 || vs.Names[0].Name != varName {
				continue
			}
			if len(vs.Values) == 0 {
				continue
			}
			// The value is a composite literal: struct{...}{...}
			cl, ok := vs.Values[0].(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, elt := range cl.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				keyIdent, ok := kv.Key.(*ast.Ident)
				if !ok {
					continue
				}
				valLit, ok := kv.Value.(*ast.BasicLit)
				if !ok {
					continue
				}
				val := valLit.Value
				if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
					val = val[1 : len(val)-1]
				}
				fields[keyIdent.Name] = val
			}
		}
	}
	return fields
}

func TestExtractStringReturnMethod(t *testing.T) {
t.Parallel()

tests := map[string]struct {
src        string
typeName   string
methodName string
want       string
wantFound  bool
}{
"alias_plain_string": {
src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}

func (Document) Alias() estype.Alias { return "my-alias" }
`,
typeName:   "Document",
methodName: "Alias",
want:       "my-alias",
wantFound:  true,
},
"alias_conversion_call": {
src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}

func (Document) Alias() estype.Alias { return estype.Alias("product") }
`,
typeName:   "Document",
methodName: "Alias",
want:       "product",
wantFound:  true,
},
"index_plain_string": {
src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}

func (Document) Index() estype.Index { return "product-000001" }
`,
typeName:   "Document",
methodName: "Index",
want:       "product-000001",
wantFound:  true,
},
"index_conversion_call": {
src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}

func (Document) Index() estype.Index { return estype.Index("product-000001") }
`,
typeName:   "Document",
methodName: "Index",
want:       "product-000001",
wantFound:  true,
},
"method_not_found": {
src: `package model

type Document struct{}
`,
typeName:   "Document",
methodName: "Alias",
want:       "",
wantFound:  false,
},
"pointer_receiver": {
src: `package model

import "github.com/tomtwinkle/es-typed-go/estype"

type Document struct{}

func (d *Document) Alias() estype.Alias { return "ptr-alias" }
`,
typeName:   "Document",
methodName: "Alias",
want:       "ptr-alias",
wantFound:  true,
},
}

for name, tt := range tests {
t.Run(name, func(t *testing.T) {
t.Parallel()

dir := t.TempDir()
if err := os.WriteFile(filepath.Join(dir, "doc.go"), []byte(tt.src), 0o644); err != nil {
t.Fatal(err)
}

fset := token.NewFileSet()
pkgs, err := parser.ParseDir(fset, dir, nil, 0)
assert.NilError(t, err)

got, found := extractStringReturnMethod(pkgs, tt.typeName, tt.methodName)
assert.Equal(t, tt.wantFound, found)
assert.Equal(t, tt.want, got)
})
}
}

// collectModelVarFields extracts the Fields sub-struct field values and top-level
// Alias/Index string values from a generated model variable.
// Returns fields as map[fieldName]path, alias, index (empty strings if absent).
func collectModelVarFields(f *ast.File, varName string) (fields map[string]string, alias, index string) {
fields = make(map[string]string)
for _, decl := range f.Decls {
gd, ok := decl.(*ast.GenDecl)
if !ok || gd.Tok != token.VAR {
continue
}
for _, spec := range gd.Specs {
vs, ok := spec.(*ast.ValueSpec)
if !ok || len(vs.Names) == 0 || vs.Names[0].Name != varName {
continue
}
if len(vs.Values) == 0 {
continue
}
// Outer composite literal: struct{...}{...}
outerCL, ok := vs.Values[0].(*ast.CompositeLit)
if !ok {
continue
}
for _, elt := range outerCL.Elts {
kv, ok := elt.(*ast.KeyValueExpr)
if !ok {
continue
}
keyIdent, ok := kv.Key.(*ast.Ident)
if !ok {
continue
}
switch keyIdent.Name {
case "Fields":
// Inner composite literal: struct{...}{field: "path", ...}
innerCL, ok := kv.Value.(*ast.CompositeLit)
if !ok {
continue
}
for _, fElt := range innerCL.Elts {
fkv, ok := fElt.(*ast.KeyValueExpr)
if !ok {
continue
}
fKey, ok := fkv.Key.(*ast.Ident)
if !ok {
continue
}
fVal, ok := fkv.Value.(*ast.BasicLit)
if !ok {
continue
}
fields[fKey.Name] = strings.Trim(fVal.Value, `"`)
}
case "Alias":
valLit, ok := kv.Value.(*ast.BasicLit)
if ok {
alias = strings.Trim(valLit.Value, `"`)
}
case "Index":
valLit, ok := kv.Value.(*ast.BasicLit)
if ok {
index = strings.Trim(valLit.Value, `"`)
}
}
}
}
}
return fields, alias, index
}

func TestGenerate_ModelFormat(t *testing.T) {
t.Parallel()

src := `package model

import "github.com/tomtwinkle/es-typed-go/estype"

//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product

type Product struct {
Status   string ` + "`" + `json:"status"` + "`" + `
Title    string ` + "`" + `json:"title"` + "`" + `
Category string ` + "`" + `json:"category"` + "`" + `
}

func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }

func (Product) Mapping() estype.Mapping {
return estype.Mapping{
Fields: []estype.MappingField{
{Path: "status",   Property: estype.NewKeywordProperty()},
{Path: "title",    Property: estype.NewTextProperty()},
{Path: "category", Property: estype.NewKeywordProperty()},
},
}
}
`
dir := t.TempDir()
outFile := filepath.Join(dir, "product_gen.go")
if err := os.WriteFile(filepath.Join(dir, "product.go"), []byte(src), 0o644); err != nil {
t.Fatal(err)
}

entries, pkgs, err := parseGoStruct(dir, "Product")
assert.NilError(t, err)
assert.Equal(t, 3, len(entries))

alias, hasAlias := extractStringReturnMethod(pkgs, "Product", "Alias")
index, hasIndex := extractStringReturnMethod(pkgs, "Product", "Index")
assert.Assert(t, hasAlias, "expected Alias() to be found")
assert.Assert(t, hasIndex, "expected Index() to be found")
assert.Equal(t, "product", alias)
assert.Equal(t, "product-000001", index)

var buf bytes.Buffer
mdata := modelTemplateData{
Package:  "esmodel",
Name:     "Product",
Fields:   entries,
Alias:    alias,
Index:    index,
HasAlias: hasAlias,
HasIndex: hasIndex,
}
assert.NilError(t, modelTemplate.Execute(&buf, mdata))

formatted, err := format.Source(buf.Bytes())
assert.NilError(t, err)

if err := os.WriteFile(outFile, formatted, 0o644); err != nil {
t.Fatal(err)
}

fset := token.NewFileSet()
f, err := parser.ParseFile(fset, outFile, nil, 0)
assert.NilError(t, err)

assertImport(t, f, `"github.com/tomtwinkle/es-typed-go/estype"`)

fields, gotAlias, gotIndex := collectModelVarFields(f, "Product")

assert.Equal(t, "status", fields["Status"])
assert.Equal(t, "title", fields["Title"])
assert.Equal(t, "category", fields["Category"])
assert.Equal(t, "product", gotAlias)
assert.Equal(t, "product-000001", gotIndex)
}
