package sql

type Settings struct {
	ChatId   int `gorm:"primary_key"`
	PmPermit bool
}

func TogglePmPermit(option bool) {
	w := &Settings{ChatId: DEFAULE_USER_ID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.PmPermit = option
	tx.Save(w)
	tx.Commit()
}

func GetSettings() *Settings {
	settings := &Settings{ChatId: DEFAULE_USER_ID}
	if SESSION.First(settings).RowsAffected == 0 {
		return &Settings{
			ChatId:   DEFAULE_USER_ID,
			PmPermit: false,
		}
	}
	return settings
}
