package queryvet

import (
	"fmt"
	"io"

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
