package queryvet

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewDDLFromReader(t *testing.T) {
	sql := `
CREATE TABLE Singers (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  SingerInfo BYTES(MAX),
) PRIMARY KEY (SingerId);

CREATE TABLE Albums (
  SingerId     INT64 NOT NULL,
  AlbumId      INT64 NOT NULL,
  AlbumTitle   STRING(MAX),
) PRIMARY KEY (SingerId, AlbumId);
`

	r := strings.NewReader(sql)
	ddl, err := NewDDLFromReader(r)
	if err != nil {
		t.Errorf("NewDDLFromReader(r) = _, %v; want _, nil", err)
	}

	want := DDL{
		"Singers": map[string]struct{}{
			"SingerId":   {},
			"FirstName":  {},
			"LastName":   {},
			"SingerInfo": {},
		},
		"Albums": map[string]struct{}{
			"SingerId":   {},
			"AlbumId":    {},
			"AlbumTitle": {},
		},
	}

	if diff := cmp.Diff(ddl, want); diff != "" {
		t.Errorf("NewDDLFromReader(r) mismatch (-got +want):\n%s", diff)
	}
}

func TestNewQuery(t *testing.T) {
	tests := []struct {
		query string
		want  *Query
	}{
		{
			query: "SELECT * FROM Singers",
			want: &Query{
				SelectTable:    "Singers",
				WhereBoolExprs: []*WhereBoolExpr{},
			},
		},
		{
			query: "SELECT * FROM Singers WHERE SingerId = @singerId",
			want: &Query{
				SelectTable: "Singers",
				WhereBoolExprs: []*WhereBoolExpr{
					{
						Table:  "Singers",
						Column: "SingerId",
					},
				},
			},
		},
		{
			query: "SELECT * FROM Singers JOIN Albums ON Singers.SingerId = Albums.SingerId WHERE Singers.SingerId = @singerId AND Albums.AlbumId = @albumId",
			want: &Query{
				SelectTable: "Singers",
				WhereBoolExprs: []*WhereBoolExpr{
					{
						Table:  "Singers",
						Column: "SingerId",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got, err := NewQuery(tt.query)
			if err != nil {
				t.Errorf("NewQuery(%q) = _, %v; want _, nil", tt.query, err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("NewQuery(%q) mismatch (-got +want):\n%s", tt.query, diff)
			}
		})
	}
}
