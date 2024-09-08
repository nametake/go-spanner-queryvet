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
