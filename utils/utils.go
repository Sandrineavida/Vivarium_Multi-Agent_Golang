package utils

import (
	"math/rand"
)

// min returns the minimum of two integers.
func Intmin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func Intmax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// // RandomPositionInRadius(p.PositionX, p.PositionY, p.Rayon)
// func RandomPositionInRadius(posX, posY, rayon int) (int, int) {
// 	// 以posX、posY为圆心，rayon为半径，随机生成一个点

// 	// rand.Seed(time.Now().UnixNano())
// 	src := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(src)

// 	// 随机生成角度（0到2π）
// 	angle := r.Float64() * 2 * math.Pi

// 	// 随机生成半径，并进行平方根处理以保证均匀分布
// 	radius := math.Sqrt(r.Float64()) * float64(rayon)

// 	// 计算笛卡尔坐标
// 	x := posX + int(radius*math.Cos(angle))
// 	y := posY + int(radius*math.Sin(angle))

// 	// fmt.Println("随机生成的坐标：", x, y)

// 	return x, y
// }

func RandomPositionInRectangle(posX, posY, rayon, x_lower_bound, x_upper_bound, y_lower_bound, y_upper_bound int) (int, int) {
	// 设置随机数种子
	// rand.Seed(time.Now().UnixNano())

	// lower_bound = 0
	// upper_bound = 9

	// 计算范围
	X_minVal := Intmax(x_lower_bound, posX-rayon)
	X_maxVal := Intmin(x_upper_bound, posX+rayon)
	Y_minVal := Intmax(y_lower_bound, posY-rayon)
	Y_maxVal := Intmin(y_upper_bound, posY+rayon)

	// 生成随机数

	X := X_minVal + rand.Intn(X_maxVal-X_minVal+1)
	Y := Y_minVal + rand.Intn(Y_maxVal-Y_minVal+1)

	return X, Y
}
