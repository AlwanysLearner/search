package DB

type Line struct {
	Tablet int
	First  int
	Angle  int
	Dist   int
	Type_1 int `gorm:"column:type"`
	Second int
}

func (Line) TableName() string {
	return "ddrel1"
}
