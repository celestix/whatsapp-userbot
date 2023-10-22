package sql

type Note struct {
	Name  string `gorm:"primary_key"`
	Value string
}

func AddNote(name, value string) {
	w := &Note{Name: name}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.Value = value
	tx.Save(w)
	tx.Commit()
}

func DeleteNote(name string) bool {
	w := &Note{Name: name}
	return SESSION.Delete(w).RowsAffected != 0
}

func GetNote(name string) *Note {
	w := Note{Name: name}
	SESSION.First(&w)
	return &w
}

func GetNotes() []Note {
	var notes []Note
	SESSION.Find(&notes)
	return notes
}
