package main

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

// Analyzer reports unchecked type assertions and requires explicit `ok` handling.
//
// In test files and external test packages, the analyzer accepts either positive
// or negative checks as long as the `ok` result is explicitly referenced. In
// non-test packages, it also accepts guard patterns that terminate control flow,
// such as `return`, `continue`, or `panic`, and ignores the analyzer's own
// implementation file to avoid self-reporting on intentionally analyzed cases.
var Analyzer = &analysis.Analyzer{
	Name: "okassertcheck",
	Doc:  "reports unchecked type assertions and requires explicit ok handling; tolerant of continue-style filtering in non-test packages and ignores self-analysis package",
	Run:  run,
}

// main runs the okassertcheck analyzer as a singlechecker-compatible command.
func main() {
	unitchecker.Main(Analyzer)
}

// run inspects each file in the current analysis pass and reports unchecked
// type assertions that do not satisfy the package-specific guard rules.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		fileInfo := analyzedFile{
			isTestFile:     isTestFile(filename),
			isTestPackage:  isTestPackage(file),
			isAnalyzerSelf: isAnalyzerSelfFile(filename),
			packageName:    file.Name.Name,
			normalizedPath: normalizePath(filename),
		}
		if fileInfo.isAnalyzerSelf {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			expr, ok := n.(*ast.TypeAssertExpr)
			if !ok {
				return true
			}

			// Type switches use x.(type) and are always allowed.
			if expr.Type == nil {
				return true
			}

			stmt, path := findParentStmtWithPath(file, expr)
			if stmt == nil {
				return true
			}

			switch s := stmt.(type) {
			case *ast.IfStmt:
				if checkIfStmtTypeAssertion(expr, s, fileInfo) {
					return true
				}
			case *ast.AssignStmt:
				if checkAssignStmtTypeAssertion(expr, s, path, fileInfo) {
					return true
				}
			}

			pass.Reportf(expr.Pos(), "unchecked type assertion; use the two-result form and verify ok")
			return true
		})
	}

	return nil, nil
}

// analyzedFile captures file-level context that influences analyzer rules.
type analyzedFile struct {
	isTestFile     bool
	isTestPackage  bool
	isAnalyzerSelf bool
	packageName    string
	normalizedPath string
}

// isTestFile reports whether the filename belongs to a Go test file.
func isTestFile(name string) bool {
	return strings.HasSuffix(filepath.Base(name), "_test.go")
}

// isTestPackage reports whether the AST file belongs to an external test package.
func isTestPackage(file *ast.File) bool {
	return strings.HasSuffix(file.Name.Name, "_test")
}

// isAnalyzerSelfFile reports whether the path is this analyzer's own source file.
func isAnalyzerSelfFile(name string) bool {
	path := normalizePath(name)
	return strings.HasSuffix(path, "/tools/okassertcheck/main.go")
}

// normalizePath converts a path to slash-separated form and removes duplicate separators.
func normalizePath(name string) string {
	return strings.ReplaceAll(filepath.ToSlash(name), "//", "/")
}

// isTestContext reports whether test-specific assertion-handling rules should apply.
func isTestContext(info analyzedFile) bool {
	return info.isTestFile || info.isTestPackage
}

// checkIfStmtTypeAssertion reports whether a type assertion used in an if
// statement is explicitly guarded in an allowed form.
func checkIfStmtTypeAssertion(expr *ast.TypeAssertExpr, ifStmt *ast.IfStmt, info analyzedFile) bool {
	assign, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok {
		return false
	}
	if len(assign.Rhs) != 1 || len(assign.Lhs) != 2 || assign.Rhs[0] != expr {
		return false
	}

	okName := identName(assign.Lhs[1])
	if okName == "" || okName == "_" {
		return false
	}

	if !conditionReferencesIdent(ifStmt.Cond, okName) {
		return false
	}

	// Tests may use either positive or negative checks.
	if isTestContext(info) {
		return true
	}

	// Non-test packages:
	// - allow explicit negative guard: if !ok { return ... } or if !ok { continue }
	// - also allow filter-style positive guard when the else path terminates:
	//   if ok { ... } else { return ... } / else { continue }
	negated, ident := negatedIdent(ifStmt.Cond)
	if negated && ident == okName {
		return blockTerminatesControlFlow(ifStmt.Body)
	}

	id, ok := ifStmt.Cond.(*ast.Ident)
	if ok && id.Name == okName {
		return elseTerminatesControlFlow(ifStmt)
	}

	return false
}

// checkAssignStmtTypeAssertion reports whether a two-result type assertion
// assigned outside an if init is followed by an allowed guard or explicit use.
func checkAssignStmtTypeAssertion(expr *ast.TypeAssertExpr, assign *ast.AssignStmt, path []ast.Node, info analyzedFile) bool {
	if len(assign.Rhs) != 1 || len(assign.Lhs) != 2 || assign.Rhs[0] != expr {
		return false
	}

	okName := identName(assign.Lhs[1])
	if okName == "" || okName == "_" {
		return false
	}

	if isTestContext(info) {
		return okIdentUsedAfter(path, assign, okName)
	}

	if sameStatementNegativeGuard(path, assign, okName) {
		return true
	}

	return hasFollowingGuard(path, assign, okName)
}

// identName returns the identifier name for expr, or an empty string when expr
// is not an identifier.
func identName(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

// conditionReferencesIdent reports whether cond contains a reference to the
// named identifier.
func conditionReferencesIdent(cond ast.Expr, name string) bool {
	if name == "" || name == "_" {
		return false
	}

	found := false
	ast.Inspect(cond, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if ok && id.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

// negatedIdent returns whether expr is a logical negation of an identifier and,
// if so, the identifier name.
func negatedIdent(expr ast.Expr) (bool, string) {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok || unary.Op != token.NOT {
		return false, ""
	}
	id, ok := unary.X.(*ast.Ident)
	if !ok {
		return false, ""
	}
	return true, id.Name
}

// okIdentUsedAfter reports whether the ok identifier is referenced by a later
// statement in the same enclosing block.
func okIdentUsedAfter(path []ast.Node, stmt ast.Stmt, okName string) bool {
	if okName == "" || okName == "_" {
		return false
	}

	block, stmtIndex := enclosingBlockAndStmtIndex(path, stmt)
	if block == nil || stmtIndex < 0 {
		return false
	}

	for i := stmtIndex + 1; i < len(block.List); i++ {
		if stmtUsesIdent(block.List[i], okName) {
			return true
		}
	}
	return false
}

// sameStatementNegativeGuard reports whether stmt is enclosed by an if block
// guarded with `!ok` whose body terminates control flow.
func sameStatementNegativeGuard(path []ast.Node, stmt ast.Stmt, okName string) bool {
	if okName == "" || okName == "_" {
		return false
	}

	for i := len(path) - 1; i >= 0; i-- {
		ifStmt, ok := path[i].(*ast.IfStmt)
		if !ok {
			continue
		}
		if ifStmt.Init != nil {
			continue
		}
		if !statementInBlock(ifStmt.Body, stmt) {
			continue
		}
		negated, ident := negatedIdent(ifStmt.Cond)
		if !negated || ident != okName {
			continue
		}
		return blockTerminatesControlFlow(ifStmt.Body)
	}
	return false
}

// statementInBlock reports whether stmt appears directly in block.
func statementInBlock(block *ast.BlockStmt, stmt ast.Stmt) bool {
	if block == nil {
		return false
	}
	for _, s := range block.List {
		if s == stmt {
			return true
		}
	}
	return false
}

// hasFollowingGuard reports whether a valid guard for okName appears after stmt
// in the same enclosing block before any unsafe use.
func hasFollowingGuard(path []ast.Node, stmt ast.Stmt, okName string) bool {
	block, stmtIndex := enclosingBlockAndStmtIndex(path, stmt)
	if block == nil || stmtIndex < 0 {
		return false
	}

	for i := stmtIndex + 1; i < len(block.List); i++ {
		next := block.List[i]

		if ifStmt, ok := next.(*ast.IfStmt); ok {
			if guardUsesOKAndTerminates(ifStmt, okName) {
				return true
			}
		}

		if nestedIfStmtUsesImmediateNegativeContinueGuard(next, okName) {
			return true
		}

		// If ok is used before a proper guard appears, the assertion result
		// is being relied on without required failure handling.
		if stmtUsesIdent(next, okName) {
			return false
		}
	}

	return false
}

// guardUsesOKAndTerminates reports whether ifStmt uses okName as an accepted
// control-flow guard.
func guardUsesOKAndTerminates(ifStmt *ast.IfStmt, okName string) bool {
	if ifStmt == nil {
		return false
	}

	negated, ident := negatedIdent(ifStmt.Cond)
	if negated && ident == okName {
		return blockTerminatesControlFlow(ifStmt.Body)
	}

	id, ok := ifStmt.Cond.(*ast.Ident)
	if ok && id.Name == okName {
		if elseTerminatesControlFlow(ifStmt) {
			return true
		}
		return positiveBodyContinuesAfterFiltering(ifStmt.Body)
	}

	return false
}

// elseTerminatesControlFlow reports whether the else branch always terminates
// control flow.
func elseTerminatesControlFlow(ifStmt *ast.IfStmt) bool {
	if ifStmt == nil {
		return false
	}

	switch elseNode := ifStmt.Else.(type) {
	case *ast.BlockStmt:
		return blockTerminatesControlFlow(elseNode)
	case *ast.IfStmt:
		return ifStmtTerminates(elseNode)
	default:
		return false
	}
}

// positiveBodyContinuesAfterFiltering reports whether a positive ok branch ends
// with `continue`, allowing the loop to filter non-matching values.
func positiveBodyContinuesAfterFiltering(block *ast.BlockStmt) bool {
	if block == nil || len(block.List) == 0 {
		return false
	}

	last := block.List[len(block.List)-1]
	branch, ok := last.(*ast.BranchStmt)
	if !ok {
		return false
	}
	return branch.Tok == token.CONTINUE
}

// nestedIfStmtUsesImmediateNegativeContinueGuard reports whether stmt is an if
// statement that guards `!ok` and immediately terminates control flow.
func nestedIfStmtUsesImmediateNegativeContinueGuard(stmt ast.Stmt, okName string) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil {
		return false
	}

	negated, ident := negatedIdent(ifStmt.Cond)
	if !negated || ident != okName {
		return false
	}

	return blockTerminatesControlFlow(ifStmt.Body)
}

// blockTerminatesControlFlow reports whether the final statement in block
// unconditionally terminates the current control-flow path.
func blockTerminatesControlFlow(block *ast.BlockStmt) bool {
	if block == nil || len(block.List) == 0 {
		return false
	}

	last := block.List[len(block.List)-1]
	switch s := last.(type) {
	case *ast.ReturnStmt:
		return true
	case *ast.BranchStmt:
		return s.Tok == token.CONTINUE
	case *ast.ExprStmt:
		call, ok := s.X.(*ast.CallExpr)
		if !ok {
			return false
		}
		id, ok := call.Fun.(*ast.Ident)
		return ok && id.Name == "panic"
	case *ast.IfStmt:
		return ifStmtTerminates(s)
	default:
		return false
	}
}

// ifStmtTerminates reports whether both branches of ifStmt terminate control flow.
func ifStmtTerminates(ifStmt *ast.IfStmt) bool {
	if ifStmt == nil {
		return false
	}

	if !blockTerminatesControlFlow(ifStmt.Body) {
		return false
	}

	return elseTerminatesControlFlow(ifStmt)
}

// enclosingBlockAndStmtIndex returns the nearest enclosing block for stmt and
// its statement index within that block.
func enclosingBlockAndStmtIndex(path []ast.Node, stmt ast.Stmt) (*ast.BlockStmt, int) {
	for i := len(path) - 1; i >= 0; i-- {
		block, ok := path[i].(*ast.BlockStmt)
		if !ok {
			continue
		}
		for j, s := range block.List {
			if s == stmt {
				return block, j
			}
		}
	}
	return nil, -1
}

// stmtUsesIdent reports whether stmt references the named identifier anywhere
// within its subtree.
func stmtUsesIdent(stmt ast.Stmt, name string) bool {
	used := false
	ast.Inspect(stmt, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if ok && id.Name == name {
			used = true
			return false
		}
		return true
	})
	return used
}

// findParentStmtWithPath locates the nearest enclosing statement for target and
// returns both that statement and the AST path leading to target.
func findParentStmtWithPath(file *ast.File, target ast.Node) (ast.Stmt, []ast.Node) {
	var result ast.Stmt
	var resultPath []ast.Node
	var stack []ast.Node

	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			return true
		}

		stack = append(stack, n)

		if n == target {
			resultPath = append([]ast.Node(nil), stack...)
			for i := len(stack) - 2; i >= 0; i-- {
				if stmt, ok := stack[i].(ast.Stmt); ok {
					result = stmt
					return false
				}
			}
			return false
		}

		return true
	})

	return result, resultPath
}
