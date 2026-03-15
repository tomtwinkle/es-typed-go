package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestToPascalCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		want string
	}{
		{"simple", "status", "Status"},
		{"dotted", "items.color", "ItemsColor"},
		{"multi_field", "title.keyword", "TitleKeyword"},
		{"underscore", "field_name", "FieldName"},
		{"deep_nested", "items.color.value", "ItemsColorValue"},
		{"underscore_and_dot", "field_name.sub_field", "FieldNameSubField"},
		{"double_underscore", "field__name", "FieldName"},
		{"leading_underscore", "_field", "Field"},
		{"trailing_underscore", "field_", "Field"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := toPascalCase(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToStructFieldName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		path string
		want string
	}{
		{"simple", "status", "Status"},
		{"dotted", "items.color", "Items_Color"},
		{"multi_field", "title.keyword", "Title_Keyword"},
		{"underscore", "field_name", "FieldName"},
		{"deep_nested", "items.color.value", "Items_Color_Value"},
		{"underscore_and_dot", "field_name.sub_field", "FieldName_SubField"},
		{"double_underscore", "field__name", "FieldName"},
		{"leading_underscore", "_field", "Field"},
		{"trailing_underscore", "field_", "Field"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := toStructFieldName(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFieldType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		typ  string
		want string
	}{
		{"keyword", "keyword", "keyword"},
		{"text", "text", "text"},
		{"nested", "nested", "nested"},
		{"empty_defaults_to_object", "", "object"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			Type:      fieldType(f.Type),
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
