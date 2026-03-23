package main

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/unitchecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "okassertcheck",
	Doc:  "reports unchecked type assertions and requires explicit ok handling; tolerant of continue-style filtering in non-test packages and ignores self-analysis package",
	Run:  run,
}

func main() {
	unitchecker.Main(Analyzer)
}

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

type analyzedFile struct {
	isTestFile     bool
	isTestPackage  bool
	isAnalyzerSelf bool
	packageName    string
	normalizedPath string
}

func isTestFile(name string) bool {
	return strings.HasSuffix(filepath.Base(name), "_test.go")
}

func isTestPackage(file *ast.File) bool {
	return strings.HasSuffix(file.Name.Name, "_test")
}

func isAnalyzerSelfFile(name string) bool {
	path := normalizePath(name)
	return strings.HasSuffix(path, "/tools/okassertcheck/main.go")
}

func normalizePath(name string) string {
	return strings.ReplaceAll(filepath.ToSlash(name), "//", "/")
}

func isTestContext(info analyzedFile) bool {
	return info.isTestFile || info.isTestPackage
}

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

func identName(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

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

func ifStmtTerminates(ifStmt *ast.IfStmt) bool {
	if ifStmt == nil {
		return false
	}

	if !blockTerminatesControlFlow(ifStmt.Body) {
		return false
	}

	return elseTerminatesControlFlow(ifStmt)
}

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
