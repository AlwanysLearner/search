package DB

type Point struct {
	Id     int
	Name   string
	Type_1 string
	Xz     float64
	Yz     float64
}

func (Point) TableName() string {
	return "obj_tb"
}
