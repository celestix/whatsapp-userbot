package sql

type Afk struct {
	ChatId  int `gorm:"primary_key"`
	Working bool
	Reason  string
}

func ToggleAfk(option bool, reason string) {
	w := &Afk{ChatId: DEFAULE_USER_ID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.Working = option
	w.Reason = reason
	tx.Save(w)
	tx.Commit()
}

func GetAfkStatus() *Afk {
	w := &Afk{ChatId: DEFAULE_USER_ID}
	if SESSION.First(w).RowsAffected == 0 {
		return &Afk{
			ChatId: DEFAULE_USER_ID,
		}
	}
	return w
}
