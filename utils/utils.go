package utils

import (
	"math"
	"math/rand"
	"time"
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

// RandomPositionInRadius(p.PositionX, p.PositionY, p.Rayon)
func RandomPositionInRadius(posX, posY, rayon int) (int, int) {
	// 以posX、posY为圆心，rayon为半径，随机生成一个点

	rand.Seed(time.Now().UnixNano()) //不一定要，但还是先写着

	// 随机生成角度（0到2π）
	angle := rand.Float64() * 2 * math.Pi

	// 随机生成半径，并进行平方根处理以保证均匀分布
	r := math.Sqrt(rand.Float64()) * float64(rayon)

	// 计算笛卡尔坐标
	x := posX + int(r*math.Cos(angle))
	y := posY + int(r*math.Sin(angle))

	return x, y
}
