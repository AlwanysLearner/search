package DB

func (point *Point) FindPointByType() []int {
	db := DataBaseSessoin()
	var pointId []int
	err := db.Model(&Point{}).Where("type_1=?", point.Type_1).Select("id").Find(&pointId).Error
	if err == nil {
		return pointId
	}
	return nil
}

func FindPointById(ids []int) []Point {
	db := DataBaseSessoin()
	var pointId []Point
	err := db.Model(&Point{}).Where("id in ?", ids).Find(&pointId).Error
	if err == nil {
		return pointId
	}
	return nil
}

func (line *Line) FindPointsByDist(k int) []Line {
	db := DataBaseSessoin()
	pointId := make([]Line, 0)
	/*for i := 1; i <= 12; i++ {
		temp := make([]int, 0)
		err := db.Table("ddrel1_part"+fmt.Sprintf("%d", i)).Where("first=? and type=? and dist>=1", line.First, line.Type_1).Select("second").
			Order("dist asc").Limit(k).Find(&temp).Error
		if err == nil {
			pointId = append(pointId, temp...)
		}
	}*/
	//err := db.Raw("SELECT second,dist,angle FROM ddrel1 where first=? and type=? ORDER BY dist asc LIMIT ?", line.First, line.Type_1, k).Scan(&pointId).Error
	//err := db.Model(line).Where("first=? and type=?", line.First, line.Type_1).Select("second,min(dist),angle").Limit(k).Find(&pointId).Error
	err := db.Table("ddrel1 FORCE INDEX (ddrel1_index2)").Where("first=? and type=?", line.First, line.Type_1).Select("second,dist,angle").
		Order("dist asc").Limit(k).Find(&pointId).Error
	if err == nil {
		return pointId
	}
	/*if len(pointId) == 0 {
		return nil
	}*/
	return nil
}

func (line *Line) FindPointsByDist1(k int) []Line {
	db := DataBaseSessoin()
	pointId := make([]Line, 0)
	/*for i := 1; i <= 12; i++ {
		temp := make([]int, 0)
		err := db.Table("ddrel1_part"+fmt.Sprintf("%d", i)).Where("first=? and type=? and dist>=1", line.First, line.Type_1).Select("second").
			Order("dist asc").Limit(k).Find(&temp).Error
		if err == nil {
			pointId = append(pointId, temp...)
		}
	}*/
	err := db.Model(line).Where("first=? and type=? and angle in (?,?,?)", line.First, line.Type_1, (line.Angle+71)%72, line.Angle, (line.Angle+1)%72).
		Select("second,dist").
		Order("dist asc").Limit(k).Find(&pointId).Error
	if err == nil {
		return pointId
	}
	/*if len(pointId) == 0 {
		return nil
	}*/
	return nil
}

func (line *Line) FindPointsByFS() *Line {
	db := DataBaseSessoin()
	pointId := &Line{}
	/*for i := 1; i <= 12; i++ {
		err := db.Table("ddrel1_part"+fmt.Sprintf("%d", i)).Where("first=? and second=?", line.First, line.Second).Select("dist", "angle").
			Find(&pointId).Error
		if err == nil && pointId.Dist != 0 {
			return pointId
		}
	}*/
	err := db.Model(line).Where("first=? and second=?", line.First, line.Second).Select("dist", "angle").
		Find(&pointId).Error
	if err == nil {
		return pointId
	}
	return nil
}

func (line *Line) FindPointsByline() []int {
	db := DataBaseSessoin()
	pointId := make([]int, 0)
	/*for i := 1; i <= 12; i++ {
		temp := make([]int, 0)
		err := db.Table("ddrel1_part"+fmt.Sprintf("%d", i)).Where("first=? and angle in (?,?,?) and dist between ? and ? and type=?",
			line.First, (line.Angle+71)%72, line.Angle, (line.Angle+1)%72, line.Dist-2, line.Dist+2, line.Type_1).
			Select("second").
			Find(&pointId).Error
		if err == nil {
			pointId = append(pointId, temp...)
		}
	}*/
	err := db.Model(line).Where("first=? and type=? and angle in (?,?,?) and dist between ? and ? ",
		line.First, line.Type_1, (line.Angle+71)%72, line.Angle, (line.Angle+1)%72, line.Dist-2, line.Dist+2).
		Select("second").
		Find(&pointId).Error
	if err == nil {
		return pointId
	}
	/*if len(pointId) == 0 {
		return nil
	}*/
	return nil
}
