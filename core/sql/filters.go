package sql

type Filter struct {
	Name  string `gorm:"primary_key"`
	Value string
}

func AddFilter(name, value string) {
	w := &Filter{Name: name}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.Value = value
	tx.Save(w)
	tx.Commit()
}

func DeleteFilter(name string) bool {
	w := &Filter{Name: name}
	return SESSION.Delete(w).RowsAffected != 0
}

func GetFilter(name string) *Filter {
	w := Filter{Name: name}
	SESSION.First(&w)
	return &w
}

func GetFilters() []Filter {
	var w []Filter
	SESSION.Find(&w)
	return w
}
