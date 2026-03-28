package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
	"gotest.tools/v3/assert"
)

func TestIsTestFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		want bool
	}{
		"plain_go_file": {
			name: "main.go",
			want: false,
		},
		"test_file": {
			name: "main_test.go",
			want: true,
		},
		"nested_test_file": {
			name: "tools/okassertcheck/main_test.go",
			want: true,
		},
		"similar_suffix_only": {
			name: "main_test.gox",
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isTestFile(tt.name))
		})
	}
}

func TestIsTestPackage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"regular_package": {
			src:  "package okassertcheck\n",
			want: false,
		},
		"external_test_package": {
			src:  "package okassertcheck_test\n",
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			file := parseFile(t, tt.src)
			assert.Equal(t, tt.want, isTestPackage(file))
		})
	}
}

func TestIsAnalyzerSelfFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		want bool
	}{
		"exact_relative_path": {
			name: "tools/okassertcheck/main.go",
			want: false,
		},
		"absolute_like_path": {
			name: "/workspace/es-typed-go/tools/okassertcheck/main.go",
			want: true,
		},
		"windows_separators": {
			name: `C:\workspace\es-typed-go\tools\okassertcheck\main.go`,
			want: false,
		},
		"different_file": {
			name: "tools/okassertcheck/main_test.go",
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isAnalyzerSelfFile(tt.name))
		})
	}
}

func TestNormalizePath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		want string
	}{
		"already_normalized": {
			name: "tools/okassertcheck/main.go",
			want: "tools/okassertcheck/main.go",
		},
		"duplicate_slashes": {
			name: "tools//okassertcheck//main.go",
			want: "tools/okassertcheck/main.go",
		},
		"windows_path": {
			name: `tools\okassertcheck\main.go`,
			want: `tools\okassertcheck\main.go`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, normalizePath(tt.name))
		})
	}
}

func TestIsTestContext(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		info analyzedFile
		want bool
	}{
		"regular_file": {
			info: analyzedFile{},
			want: false,
		},
		"test_file": {
			info: analyzedFile{isTestFile: true},
			want: true,
		},
		"test_package": {
			info: analyzedFile{isTestPackage: true},
			want: true,
		},
		"both": {
			info: analyzedFile{isTestFile: true, isTestPackage: true},
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isTestContext(tt.info))
		})
	}
}

func TestIdentName(t *testing.T) {
	t.Parallel()

	ident := &ast.Ident{Name: "ok"}
	assert.Equal(t, "ok", identName(ident))

	call := &ast.CallExpr{}
	assert.Equal(t, "", identName(call))
}

func TestConditionReferencesIdent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expr string
		name string
		want bool
	}{
		"direct_reference": {
			expr: "ok",
			name: "ok",
			want: true,
		},
		"negated_reference": {
			expr: "!ok",
			name: "ok",
			want: true,
		},
		"compound_reference": {
			expr: "ok && ready",
			name: "ok",
			want: true,
		},
		"missing_reference": {
			expr: "ready && done",
			name: "ok",
			want: false,
		},
		"blank_identifier_rejected": {
			expr: "ok",
			name: "_",
			want: false,
		},
		"empty_identifier_rejected": {
			expr: "ok",
			name: "",
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			expr := parseExpr(t, tt.expr)
			assert.Equal(t, tt.want, conditionReferencesIdent(expr, tt.name))
		})
	}
}

func TestNegatedIdent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expr     string
		wantBool bool
		wantName string
	}{
		"negated_ident": {
			expr:     "!ok",
			wantBool: true,
			wantName: "ok",
		},
		"non_negated_ident": {
			expr:     "ok",
			wantBool: false,
			wantName: "",
		},
		"negated_selector": {
			expr:     "!state.ok",
			wantBool: false,
			wantName: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gotBool, gotName := negatedIdent(parseExpr(t, tt.expr))
			assert.Equal(t, tt.wantBool, gotBool)
			assert.Equal(t, tt.wantName, gotName)
		})
	}
}

func TestStatementInBlock(t *testing.T) {
	t.Parallel()

	file := parseFile(t, `
package sample

func f() {
	a := 1
	b := 2
	_ = a + b
}
`)

	fn := firstFuncDecl(t, file)
	block := fn.Body
	assert.Assert(t, len(block.List) >= 2)

	assert.Equal(t, true, statementInBlock(block, block.List[0]))
	assert.Equal(t, false, statementInBlock(block, &ast.ReturnStmt{}))
	assert.Equal(t, false, statementInBlock(nil, block.List[0]))
}

func TestStmtUsesIdent(t *testing.T) {
	t.Parallel()

	stmt := parseStmt(t, `_ = ok`)
	assert.Equal(t, true, stmtUsesIdent(stmt, "ok"))
	assert.Equal(t, false, stmtUsesIdent(stmt, "missing"))
}

func TestBlockTerminatesControlFlow(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stmt string
		want bool
	}{
		"return": {
			stmt: `return`,
			want: true,
		},
		"continue": {
			stmt: `continue`,
			want: true,
		},
		"panic_call": {
			stmt: `panic("boom")`,
			want: true,
		},
		"plain_expression": {
			stmt: `_ = 1`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			block := &ast.BlockStmt{List: []ast.Stmt{parseStmt(t, tt.stmt)}}
			assert.Equal(t, tt.want, blockTerminatesControlFlow(block))
		})
	}

	t.Run("nested_if_with_both_branches_terminating", func(t *testing.T) {
		t.Parallel()
		ifStmt := parseIfStmt(t, `if ok { return } else { panic("boom") }`)
		block := &ast.BlockStmt{List: []ast.Stmt{ifStmt}}
		assert.Equal(t, true, blockTerminatesControlFlow(block))
	})

	t.Run("empty_block", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, false, blockTerminatesControlFlow(&ast.BlockStmt{}))
	})
}

func TestElseTerminatesControlFlow(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"else_block_return": {
			src:  `if ok { _ = value } else { return }`,
			want: true,
		},
		"else_if_terminates": {
			src:  `if ok { _ = value } else if ready { return } else { panic("boom") }`,
			want: true,
		},
		"missing_else": {
			src:  `if ok { _ = value }`,
			want: false,
		},
		"else_does_not_terminate": {
			src:  `if ok { _ = value } else { _ = value }`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, elseTerminatesControlFlow(parseIfStmt(t, tt.src)))
		})
	}
}

func TestIfStmtTerminates(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"both_branches_terminate": {
			src:  `if ok { return } else { panic("boom") }`,
			want: true,
		},
		"body_only_terminates": {
			src:  `if ok { return } else { _ = value }`,
			want: false,
		},
		"else_if_chain_terminates": {
			src:  `if ok { panic("boom") } else if ready { return } else { continue }`,
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, ifStmtTerminates(parseIfStmt(t, tt.src)))
		})
	}
}

func TestPositiveBodyContinuesAfterFiltering(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"continue_last": {
			src: `if ok {
				_ = value
				continue
			}`,
			want: true,
		},
		"return_last": {
			src: `if ok {
				return
			}`,
			want: false,
		},
		"empty_body": {
			src:  `if ok {}`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ifStmt := parseIfStmt(t, tt.src)
			assert.Equal(t, tt.want, positiveBodyContinuesAfterFiltering(ifStmt.Body))
		})
	}
}

func TestGuardUsesOKAndTerminates(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"negative_guard_return": {
			src:  `if !ok { return }`,
			want: true,
		},
		"positive_guard_else_return": {
			src:  `if ok { _ = value } else { return }`,
			want: true,
		},
		"positive_guard_continue_filter": {
			src: `if ok {
				_ = value
				continue
			}`,
			want: true,
		},
		"wrong_identifier": {
			src:  `if !ready { return }`,
			want: false,
		},
		"no_termination": {
			src:  `if !ok { _ = value }`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, guardUsesOKAndTerminates(parseIfStmt(t, tt.src), "ok"))
		})
	}
}

func TestNestedIfStmtUsesImmediateNegativeContinueGuard(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		stmt string
		want bool
	}{
		"negative_continue": {
			stmt: `if !ok { continue }`,
			want: true,
		},
		"negative_return": {
			stmt: `if !ok { return }`,
			want: true,
		},
		"with_init_not_allowed": {
			stmt: `if x := 1; !ok { continue }`,
			want: false,
		},
		"wrong_identifier": {
			stmt: `if !ready { continue }`,
			want: false,
		},
		"not_if_stmt": {
			stmt: `_ = ok`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, nestedIfStmtUsesImmediateNegativeContinueGuard(parseStmt(t, tt.stmt), "ok"))
		})
	}
}

func TestEnclosingBlockAndStmtIndex(t *testing.T) {
	t.Parallel()

	file := parseFile(t, `
package sample

func f() {
	ok := true
	if !ok {
		return
	}
	_ = ok
}
`)

	fn := firstFuncDecl(t, file)
	assert.Assert(t, len(fn.Body.List) == 3)

	block, idx := enclosingBlockAndStmtIndex([]ast.Node{file, fn, fn.Body}, fn.Body.List[1])
	assert.Assert(t, block != nil)
	assert.Equal(t, 1, idx)

	block, idx = enclosingBlockAndStmtIndex([]ast.Node{file, fn}, fn.Body.List[1])
	assert.Assert(t, block == nil)
	assert.Equal(t, -1, idx)
}

func TestFindParentStmtWithPath(t *testing.T) {
	t.Parallel()

	file := parseFile(t, `
package sample

func f(v any) {
	value, ok := v.(string)
	_, _ = value, ok
}
`)

	expr := firstTypeAssertExpr(t, file)
	stmt, path := findParentStmtWithPath(file, expr)

	_, isAssign := stmt.(*ast.AssignStmt)
	assert.Equal(t, true, isAssign)
	assert.Assert(t, len(path) > 0)
	assert.Assert(t, path[len(path)-1] == expr)
}

func TestCheckIfStmtTypeAssertion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		info analyzedFile
		want bool
	}{
		"test_context_positive_condition_allowed": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); ok {
		_ = value
	}
}`,
			info: analyzedFile{isTestFile: true},
			want: true,
		},
		"test_context_negative_condition_allowed": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); !ok {
		return
	} else {
		_ = value
	}
}`,
			info: analyzedFile{isTestPackage: true},
			want: true,
		},
		"non_test_negative_guard_return": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); !ok {
		return
	} else {
		_ = value
	}
}`,
			info: analyzedFile{},
			want: true,
		},
		"non_test_positive_guard_else_return": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); ok {
		_ = value
	} else {
		return
	}
}`,
			info: analyzedFile{},
			want: true,
		},
		"non_test_positive_without_else_not_allowed": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); ok {
		_ = value
	}
}`,
			info: analyzedFile{},
			want: false,
		},
		"blank_ok_not_allowed": {
			src: `package sample
func f(v any) {
	if value, _ := v.(string); value != "" {
		return
	}
}`,
			info: analyzedFile{},
			want: false,
		},
		"condition_must_reference_ok": {
			src: `package sample
func f(v any) {
	if value, ok := v.(string); value != "" {
		_ = ok
	}
}`,
			info: analyzedFile{isTestFile: true},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file := parseFile(t, tt.src)
			expr := firstTypeAssertExpr(t, file)
			ifStmt := firstIfStmt(t, file)
			assert.Equal(t, tt.want, checkIfStmtTypeAssertion(expr, ifStmt, tt.info))
		})
	}
}

func TestCheckAssignStmtTypeAssertion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		info analyzedFile
		want bool
	}{
		"test_context_later_ok_use": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	if !ok {
		t.Fatal("unexpected type")
	}
	_ = value
}`,
			info: analyzedFile{isTestFile: true},
			want: true,
		},
		"test_context_unused_ok": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	_ = value
}`,
			info: analyzedFile{isTestFile: true},
			want: false,
		},
		"non_test_following_negative_guard": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	if !ok {
		return
	}
	_ = value
}`,
			info: analyzedFile{},
			want: true,
		},
		"non_test_same_statement_guarded_block": {
			src: `package sample
func f(values []any) {
	for _, v := range values {
		ok := true
		if !ok {
			continue
		}
		value, ok := v.(string)
		_ = value
	}
}`,
			info: analyzedFile{},
			want: false,
		},
		"non_test_immediate_negative_guard_in_parent_block": {
			src: `package sample
func f(values []any) {
	for _, v := range values {
		value, ok := v.(string)
		if !ok {
			continue
		}
		_ = value
	}
}`,
			info: analyzedFile{},
			want: true,
		},
		"blank_ok_rejected": {
			src: `package sample
func f(v any) {
	value, _ := v.(string)
	_ = value
}`,
			info: analyzedFile{},
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file := parseFile(t, tt.src)
			expr := firstTypeAssertExpr(t, file)
			assign := firstAssignStmtContainingTypeAssert(t, file)
			_, path := findParentStmtWithPath(file, expr)
			assert.Equal(t, tt.want, checkAssignStmtTypeAssertion(expr, assign, path, tt.info))
		})
	}
}

func TestOkIdentUsedAfter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"later_if_uses_ok": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	if !ok {
		return
	}
	_ = value
}`,
			want: true,
		},
		"later_assignment_uses_ok": {
			src: `package sample
func f(v any) {
	_, ok := v.(string)
	_ = ok
}`,
			want: true,
		},
		"no_later_use": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	_ = value
}`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			file := parseFile(t, tt.src)
			expr := firstTypeAssertExpr(t, file)
			assign := firstAssignStmtContainingTypeAssert(t, file)
			_, path := findParentStmtWithPath(file, expr)
			assert.Equal(t, tt.want, okIdentUsedAfter(path, assign, "ok"))
		})
	}

	t.Run("blank_identifier_rejected", func(t *testing.T) {
		t.Parallel()
		file := parseFile(t, `
package sample
func f(v any) {
	_, ok := v.(string)
	_ = ok
}`)
		expr := firstTypeAssertExpr(t, file)
		assign := firstAssignStmtContainingTypeAssert(t, file)
		_, path := findParentStmtWithPath(file, expr)
		assert.Equal(t, false, okIdentUsedAfter(path, assign, "_"))
	})
}

func TestSameStatementNegativeGuard(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"guarded_by_parent_if_continue": {
			src: `package sample
func f(values []any, ok bool) {
	for _, v := range values {
		if !ok {
			continue
		}
		value, matched := v.(string)
		_, _ = value, matched
	}
}`,
			want: false,
		},
		"guarded_by_parent_if_return": {
			src: `package sample
func f(v any, ok bool) {
	if !ok {
		return
	}
	value, matched := v.(string)
	_, _ = value, matched
}`,
			want: false,
		},
		"assignment_inside_negative_if_body": {
			src: `package sample
func f(v any) {
	if !ok {
		value, ok := v.(string)
		_, _ = value, ok
		return
	}
}`,
			want: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			file := parseFile(t, tt.src)
			expr := firstTypeAssertExpr(t, file)
			assign := firstAssignStmtContainingTypeAssert(t, file)
			_, path := findParentStmtWithPath(file, expr)
			assert.Equal(t, tt.want, sameStatementNegativeGuard(path, assign, "ok"))
		})
	}
}

func TestHasFollowingGuard(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src  string
		want bool
	}{
		"negative_return_guard": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	if !ok {
		return
	}
	_ = value
}`,
			want: true,
		},
		"positive_else_return_guard": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	if ok {
		_ = value
	} else {
		return
	}
}`,
			want: true,
		},
		"continue_filter_guard": {
			src: `package sample
func f(values []any) {
	for _, v := range values {
		value, ok := v.(string)
		if !ok {
			continue
		}
		_ = value
	}
}`,
			want: true,
		},
		"use_before_guard": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	_ = ok
	if !ok {
		return
	}
	_ = value
}`,
			want: false,
		},
		"no_guard": {
			src: `package sample
func f(v any) {
	value, ok := v.(string)
	_ = value
	_ = ok
}`,
			want: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			file := parseFile(t, tt.src)
			expr := firstTypeAssertExpr(t, file)
			assign := firstAssignStmtContainingTypeAssert(t, file)
			_, path := findParentStmtWithPath(file, expr)
			assert.Equal(t, tt.want, hasFollowingGuard(path, assign, "ok"))
		})
	}
}

func TestRunReportsUncheckedAssertions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		src       string
		filename  string
		wantCount int
	}{
		"single_result_assertion_reported": {
			src: `package sample
func f(v any) {
	_ = v.(string)
}`,
			filename:  "sample.go",
			wantCount: 1,
		},
		"test_file_if_init_allowed": {
			src: `package sample
func f(v any) {
	if s, ok := v.(string); ok {
		_ = s
	}
}`,
			filename:  "sample.go",
			wantCount: 1,
		},
		"analyzer_self_file_ignored": {
			src: `package main
func f(v any) {
	_ = v.(string)
}`,
			filename:  "/workspace/es-typed-go/tools/okassertcheck/main.go",
			wantCount: 0,
		},
		"type_switch_ignored": {
			src: `package sample
func f(v any) {
	switch v.(type) {
	case string:
	}
}`,
			filename:  "sample.go",
			wantCount: 0,
		},
		"non_test_assign_with_guard_allowed": {
			src: `package sample
func f(v any) {
	s, ok := v.(string)
	if !ok {
		return
	}
	_ = s
}`,
			filename:  "sample.go",
			wantCount: 0,
		},
		"non_test_assign_without_guard_reported": {
			src: `package sample
func f(v any) {
	s, ok := v.(string)
	_, _ = s, ok
}`,
			filename:  "sample.go",
			wantCount: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file := parseFile(t, tt.src)
			fset := token.NewFileSet()
			parsed, err := parser.ParseFile(fset, tt.filename, tt.src, 0)
			assert.NilError(t, err)

			var reports []analysisReport
			pass := passForFile(fset, parsed, func(pos token.Pos, msg string) {
				reports = append(reports, analysisReport{pos: pos, msg: msg})
			})

			_, err = run(pass)
			assert.NilError(t, err)
			assert.Equal(t, tt.wantCount, len(reports))

			_ = file
		})
	}
}

type analysisReport struct {
	pos token.Pos
	msg string
}

type fakePass struct {
	Files []*ast.File
	Fset  *token.FileSet
}

func passForFile(fset *token.FileSet, file *ast.File, report func(token.Pos, string)) *analysis.Pass {
	return &analysis.Pass{
		Analyzer: Analyzer,
		Fset:     fset,
		Files:    []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			report(d.Pos, d.Message)
		},
	}
}

func parseFile(t *testing.T, src string) *ast.File {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sample.go", src, 0)
	assert.NilError(t, err)
	return file
}

func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()

	expr, err := parser.ParseExpr(src)
	assert.NilError(t, err)
	return expr
}

func parseStmt(t *testing.T, src string) ast.Stmt {
	t.Helper()

	file := parseFile(t, "package sample\nfunc f(){\n"+src+"\n}\n")
	fn := firstFuncDecl(t, file)
	assert.Assert(t, len(fn.Body.List) > 0)
	return fn.Body.List[0]
}

func parseIfStmt(t *testing.T, src string) *ast.IfStmt {
	t.Helper()

	stmt := parseStmt(t, src)
	ifStmt, ok := stmt.(*ast.IfStmt)
	assert.Assert(t, ok)
	return ifStmt
}

func firstFuncDecl(t *testing.T, file *ast.File) *ast.FuncDecl {
	t.Helper()

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fn
		}
	}
	t.Fatalf("expected function declaration")
	return nil
}

func firstIfStmt(t *testing.T, file *ast.File) *ast.IfStmt {
	t.Helper()

	var result *ast.IfStmt
	ast.Inspect(file, func(n ast.Node) bool {
		if result != nil {
			return false
		}
		if stmt, ok := n.(*ast.IfStmt); ok {
			result = stmt
			return false
		}
		return true
	})
	assert.Assert(t, result != nil)
	return result
}

func firstTypeAssertExpr(t *testing.T, file *ast.File) *ast.TypeAssertExpr {
	t.Helper()

	var result *ast.TypeAssertExpr
	ast.Inspect(file, func(n ast.Node) bool {
		if result != nil {
			return false
		}
		if expr, ok := n.(*ast.TypeAssertExpr); ok && expr.Type != nil {
			result = expr
			return false
		}
		return true
	})
	assert.Assert(t, result != nil)
	return result
}

func firstAssignStmtContainingTypeAssert(t *testing.T, file *ast.File) *ast.AssignStmt {
	t.Helper()

	var result *ast.AssignStmt
	ast.Inspect(file, func(n ast.Node) bool {
		if result != nil {
			return false
		}
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, rhs := range assign.Rhs {
			if _, ok := rhs.(*ast.TypeAssertExpr); ok {
				result = assign
				return false
			}
		}
		return true
	})
	assert.Assert(t, result != nil)
	return result
}
