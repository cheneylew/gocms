package models



type Menu struct {
	MenuID int64
	ParentId int64
	Name string
	Class string
	Icon string
	Url string
	IsActive bool

	Chirldren []*Menu
}

