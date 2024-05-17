package main

import (
	"fmt"
	"math"
	"raft/DB"
	"sort"
	"sync"
	"time"
)

var wg sync.WaitGroup
var wg1 sync.WaitGroup

type MatchCer struct {
	length float64
	bili   []float64
	angle  []int
}

type Result struct {
	dsum   int
	points map[int]struct{}
}

func AngleFromNorth(p1, p2 DB.Point) int {
	// 计算两点间的差值
	deltaX := p2.Xz - p1.Xz
	deltaY := p2.Yz - p1.Yz
	// 使用 Atan2 计算角度，注意参数顺序是 Y 坐标差，X 坐标差
	// Atan2 返回的是从正x轴（正东）逆时针旋转到向量的角度
	angle := math.Atan2(deltaY, deltaX)

	// 将角度从弧度转换为度
	degrees := angle * (180 / math.Pi)

	// 由于 Atan2 返回的是从正东（X轴正方向）逆时针到向量的角度，调整为从正北开始的角度
	degrees = 90 - degrees // 从正北开始顺时针
	if degrees < 0 {
		degrees += 360 // 调整负角度
	}
	return int(degrees / 5)
}

// 查询顺序
func DescidedRank(points []DB.Point) ([]DB.Point, *MatchCer) {
	//fmt.Println(points)
	sort.Slice(points, func(i, j int) bool {
		return DB.TypeMap[points[i].Type_1] < DB.TypeMap[points[j].Type_1]
	})
	index := 1
	maxdist := math.Sqrt((points[0].Xz-points[index].Xz)*(points[0].Xz-points[index].Xz) +
		(points[0].Yz-points[index].Yz)*(points[0].Yz-points[index].Yz))
	for i := index; i < len(points); i++ {
		temp := math.Sqrt((points[0].Xz-points[i].Xz)*(points[0].Xz-points[i].Xz) +
			(points[0].Yz-points[i].Yz)*(points[0].Yz-points[i].Yz))
		if temp > maxdist {
			index, maxdist = i, temp
		}
	}
	points[1], points[index] = points[index], points[1]
	//fmt.Println(points)
	bili := make([]float64, len(points)-2)
	angle := make([]int, len(points)-2)
	anglean := AngleFromNorth(points[0], points[1])
	for i := 2; i < len(points); i++ {
		temp := math.Sqrt((points[0].Xz-points[i].Xz)*(points[0].Xz-points[i].Xz) +
			(points[0].Yz-points[i].Yz)*(points[0].Yz-points[i].Yz))
		bili[i-2] = temp / maxdist
		tempangle := AngleFromNorth(points[0], points[i])
		angle[i-2] = tempangle - anglean
	}
	matchCer := &MatchCer{length: maxdist, bili: bili, angle: angle}
	return points, matchCer
}

func DescidedRank1(points []DB.Point) ([]DB.Point, *MatchCer) {
	//fmt.Println(points)
	sort.Slice(points, func(i, j int) bool {
		return DB.TypeMap[points[i].Type_1] < DB.TypeMap[points[j].Type_1]
	})
	index := 1
	maxdist := math.Sqrt((points[0].Xz-points[index].Xz)*(points[0].Xz-points[index].Xz) +
		(points[0].Yz-points[index].Yz)*(points[0].Yz-points[index].Yz))
	for i := index; i < len(points); i++ {
		temp := math.Sqrt((points[0].Xz-points[i].Xz)*(points[0].Xz-points[i].Xz) +
			(points[0].Yz-points[i].Yz)*(points[0].Yz-points[i].Yz))
		if temp > maxdist {
			index, maxdist = i, temp
		}
	}
	points[1], points[index] = points[index], points[1]
	//fmt.Println(points)
	bili := make([]float64, len(points)-2)
	angle := make([]int, len(points)-1)
	anglean := AngleFromNorth(points[0], points[1])
	for i := 2; i < len(points); i++ {
		temp := math.Sqrt((points[0].Xz-points[i].Xz)*(points[0].Xz-points[i].Xz) +
			(points[0].Yz-points[i].Yz)*(points[0].Yz-points[i].Yz))
		bili[i-2] = temp / maxdist
		angle[i-2] = AngleFromNorth(points[0], points[i])
	}
	angle[len(points)-2] = anglean
	matchCer := &MatchCer{length: maxdist, bili: bili, angle: angle}
	return points, matchCer
}

// 查询
func search(points []DB.Point) []int {
	ch := make(chan []Result, 1500)
	/*ch := make(chan struct{}, 1500)
	for i := 0; i < 1500; i++ {
		ch <- struct{}{}
	}*/
	//1.确定查询顺序
	var matchCer *MatchCer
	//var lock sync.Mutex
	points, matchCer = DescidedRank(points)
	//2.查询第一个点
	point := &DB.Point{Type_1: points[0].Type_1}
	point1_ids := point.FindPointByType()
	results := make([]Result, 0)
	for _, v := range point1_ids {
		//3.查询第二个点
		v := v

		//<-ch
		wg.Add(1)
		go func() {
			/*defer func() {
				ch <- struct{}{}
			}()*/
			defer wg.Done()
			//db := DB.DataBaseSessoin()
			line := &DB.Line{First: v, Type_1: DB.TypeMap[points[1].Type_1]}
			point2_ids := line.FindPointsByDist(3)
			//db := DB.DataBaseSessoin()
			//line := &DB.Line{First: v, Type_1: DB.TypeMap[points[1].Type_1], Angle: matchCer.angle[len(points)-2]}
			//point2_ids := line.FindPointsByDist1(3)
			result1 := make([]Result, 0)
			//fmt.Println(point2_ids)
			for _, t := range point2_ids {
				/*line2 := &DB.Line{First: v, Second: t.Second}
				line3 := line2.FindPointsByFS()
				if line3 == nil {
					break
				}*/
				result := Result{dsum: t.Dist}
				result.points = make(map[int]struct{})
				result.points[v] = struct{}{}
				result.points[t.Second] = struct{}{}
				//4.查询剩余点
				i := 2
				for ; i < len(points); i++ {
					//fmt.Println(len(matchCer.angle), line3)
					line1 := &DB.Line{First: v, Angle: matchCer.angle[i-2] + t.Angle, Dist: int(float64(t.Dist) * matchCer.bili[i-2]), Type_1: DB.TypeMap[points[i].Type_1]}
					/*line1 := &DB.Line{First: v, Angle: matchCer.angle[i-2] + matchCer.angle[len(points)-2],
					Dist: int(float64(line3.Dist) * matchCer.bili[i-2]), Type_1: DB.TypeMap[points[i].Type_1]}*/
					//fmt.Println(line1)
					sresult := line1.FindPointsByline()
					//提前结束
					if len(sresult) == 0 {
						break
					}
					for j := 0; j < len(sresult); j++ {
						result.points[sresult[j]] = struct{}{}
					}
				}
				if i == len(points) {
					//fmt.Println(v, t, line3)
					result1 = append(result1, result)
				}
			}
			if len(result1) != 0 {
				ch <- result1
				/*lock.Lock()
				results = append(results, result1...)
				lock.Unlock()*/
			}
		}()
	}
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for {
			value := <-ch
			if len(value) == 0 {
				break
			}
			results = append(results, value...)
		}
	}()
	wg.Wait()
	ch <- []Result{}
	wg1.Wait()
	//5.排序结果集
	sort.Slice(results, func(i, j int) bool {
		return results[i].dsum < results[j].dsum
	})
	//fmt.Println(results)
	//6.计算正确率
	isOne := 0
	isFive := 0
	isTen := 0
	cg := 0
	for i := 0; i < len(results); i++ {
		istrue := true
		for j := 0; j < len(points); j++ {
			if _, ok := results[i].points[points[j].Id]; !ok {
				istrue = false
				break
			}
		}
		if istrue {
			cg++
			if i < 10 {
				isTen++
			}
			if i < 5 {
				isFive++
			}
			if i == 0 {
				isOne++
			}
			break
		}
	}
	return []int{isOne, isFive, isTen, cg}
}

func search1(points []DB.Point) []int {
	ch := make(chan []Result, 1500)
	/*ch := make(chan struct{}, 1500)
	for i := 0; i < 1500; i++ {
		ch <- struct{}{}
	}*/
	//1.确定查询顺序
	var matchCer *MatchCer
	//var lock sync.Mutex
	points, matchCer = DescidedRank1(points)
	//fmt.Println(points, matchCer)
	//2.查询第一个点
	point := &DB.Point{Type_1: points[0].Type_1}
	point1_ids := point.FindPointByType()
	//fmt.Println(point1_ids)
	results := make([]Result, 0)
	for _, v := range point1_ids {
		//3.查询第二个点
		v := v

		//<-ch
		wg.Add(1)
		go func() {
			/*defer func() {
				ch <- struct{}{}
			}()*/
			defer wg.Done()

			//db := DB.DataBaseSessoin()
			line := &DB.Line{First: v, Type_1: DB.TypeMap[points[1].Type_1], Angle: matchCer.angle[len(points)-2]}
			point2_ids := line.FindPointsByDist1(3)
			/*if v == 3866 {
				fmt.Println(point2_ids)
			}*/
			//db := DB.DataBaseSessoin()
			//line := &DB.Line{First: v, Type_1: DB.TypeMap[points[1].Type_1], Angle: matchCer.angle[len(points)-2]}
			//point2_ids := line.FindPointsByDist1(3)
			result1 := make([]Result, 0)
			//fmt.Println(point2_ids)
			for _, t := range point2_ids {
				/*line2 := &DB.Line{First: v, Second: t.Second}
				line3 := line2.FindPointsByFS()
				if line3 == nil {
					break
				}*/
				result := Result{dsum: t.Dist}
				result.points = make(map[int]struct{})
				result.points[v] = struct{}{}
				result.points[t.Second] = struct{}{}
				//4.查询剩余点
				i := 2
				for ; i < len(points); i++ {
					//fmt.Println(len(matchCer.angle), line3)
					line1 := &DB.Line{First: v, Angle: matchCer.angle[i-2], Dist: int(float64(t.Dist) * matchCer.bili[i-2]), Type_1: DB.TypeMap[points[i].Type_1]}
					/*line1 := &DB.Line{First: v, Angle: matchCer.angle[i-2] + matchCer.angle[len(points)-2],
					Dist: int(float64(line3.Dist) * matchCer.bili[i-2]), Type_1: DB.TypeMap[points[i].Type_1]}*/
					//fmt.Println(line1)
					sresult := line1.FindPointsByline()
					//提前结束
					if len(sresult) == 0 {
						break
					}
					for j := 0; j < len(sresult); j++ {
						result.points[sresult[j]] = struct{}{}
					}
				}
				if i == len(points) {
					//fmt.Println(v, t, line3)
					result1 = append(result1, result)
				}
			}
			if len(result1) != 0 {
				ch <- result1
				/*lock.Lock()
				results = append(results, result1...)
				lock.Unlock()*/
			}
		}()
	}
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for {
			value := <-ch
			if len(value) == 0 {
				break
			}
			results = append(results, value...)
		}
	}()
	wg.Wait()
	ch <- []Result{}
	wg1.Wait()
	//5.排序结果集
	sort.Slice(results, func(i, j int) bool {
		return results[i].dsum < results[j].dsum
	})
	//fmt.Println(results)
	//6.计算正确率
	isOne := 0
	isFive := 0
	isTen := 0
	cg := 0
	for i := 0; i < len(results); i++ {
		istrue := true
		for j := 0; j < len(points); j++ {
			if _, ok := results[i].points[points[j].Id]; !ok {
				istrue = false
				break
			}
		}
		if istrue {
			cg++
			if i < 10 {
				isTen++
			}
			if i < 5 {
				isFive++
			}
			if i == 0 {
				isOne++
			}
			break
		}
	}
	return []int{isOne, isFive, isTen, cg}
}

func search3(data [][]int) {
	acc := make([]int, 4)
	startTime := time.Now()
	for _, v := range data {
		points := DB.FindPointById(v)
		result := search1(points)
		for i := 0; i < 4; i++ {
			acc[i] += result[i]
		}
	}
	elapsedTime := time.Since(startTime)
	fmt.Println(acc[0], acc[1], acc[2], acc[3])
	fmt.Println(elapsedTime)
}

func main() {
	DB.InitDatabase()
	/*db := DB.DataBaseSessoin()
	l := DB.Line{}
	db.Raw("SELECT second,dist,angle FROM `ddrel1` WHERE first=4539 and type=17 ORDER BY dist asc LIMIT 1").Scan(&l)
	fmt.Println(l)*/
	/*s := []int{889, 1358, 3866}
	point := DB.FindPointById(s)
	startTime := time.Now()
	result := search1(point)
	elapsedTime := time.Since(startTime)
	fmt.Println(elapsedTime)
	fmt.Println(result[0], result[1], result[2], result[3])*/
	/*line2 := &DB.Line{First: 31099, Second: 10231}
	line3 := line2.FindPointsByFS()
	fmt.Println(line3)*/
	search3(DB.Data3)
	search3(DB.Data5)
	search3(DB.Data7)
	//search3(DB.Data10)
}
