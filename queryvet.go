package queryvet

type DDL map[string]map[string]struct{}

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
