package queryvet

import (
	"fmt"
	"io"
	"slices"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

type DDL map[string]map[string]struct{}

func NewDDLFromReader(r io.Reader) (DDL, error) {
	sql, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	file := &token.File{
		Buffer: string(sql),
	}

	p := memefish.Parser{
		Lexer: &memefish.Lexer{File: file},
	}

	stmt, err := p.ParseDDLs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse DDL: %w", err)
	}

	ddl := DDL{}
	for _, s := range stmt {
		ct, ok := s.(*ast.CreateTable)
		if !ok {
			continue
		}
		for _, c := range ct.Columns {
			ddl.Add(ct.Name.Name, c.Name.Name)
		}
	}

	return ddl, nil
}

func (d DDL) Add(table, column string) {
	if _, ok := d[table]; !ok {
		d[table] = make(map[string]struct{})
	}
	d[table][column] = struct{}{}
}

func (d DDL) Has(table, column string) bool {
	if _, ok := d[table]; !ok {
		return false
	}
	_, ok := d[table][column]
	return ok
}

type Query struct {
	SelectTable    string
	WhereBoolExprs []*WhereBoolExpr
}

type WhereBoolExpr struct {
	Table  string
	Column string
}

func NewQuery(query string) (*Query, error) {
	file := &token.File{
		Buffer: query,
	}
	p := &memefish.Parser{
		Lexer: &memefish.Lexer{File: file},
	}

	stmt, err := p.ParseQuery()
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	selectStmt, ok := stmt.Query.(*ast.Select)
	if !ok {
		return nil, fmt.Errorf("expected SELECT statement, got %T", stmt)
	}

	tableNames := tableNamesFromSource(selectStmt.From.Source)
	if len(tableNames) == 0 {
		return nil, fmt.Errorf("expected table name, got %T", selectStmt.From.Source)
	}

	baseTableName := tableNames[0]

	if selectStmt.Where == nil {
		return &Query{
			SelectTable:    baseTableName,
			WhereBoolExprs: []*WhereBoolExpr{},
		}, nil
	}
	binaryExpr := selectStmt.Where.Expr.(*ast.BinaryExpr)

	fmt.Printf("%#v\n", binaryExpr)

	fmt.Printf("%#v\n", binaryExprToWhereBoolExpr(binaryExpr))

	return &Query{
		SelectTable:    baseTableName,
		WhereBoolExprs: []*WhereBoolExpr{},
	}, nil
}

func tableNamesFromSource(source ast.TableExpr) []string {
	switch s := source.(type) {
	case *ast.TableName:
		return []string{s.Table.Name}
	case *ast.Join:
		return slices.Concat(tableNamesFromSource(s.Left), tableNamesFromSource(s.Right))
	default:
		return []string{}
	}
}

func binaryExprToWhereBoolExpr(binaryExpr *ast.BinaryExpr) []*WhereBoolExpr {
	switch binaryExpr.Op {
	case "=":
		fmt.Printf("%#v\n", binaryExpr.Left)
		fmt.Printf("%#v\n", binaryExpr.Right)
		return []*WhereBoolExpr{}
	case "AND":
		return slices.Concat(binaryExprToWhereBoolExpr(binaryExpr.Left.(*ast.BinaryExpr)), binaryExprToWhereBoolExpr(binaryExpr.Right.(*ast.BinaryExpr)))
	default:
		return []*WhereBoolExpr{}
	}
}
