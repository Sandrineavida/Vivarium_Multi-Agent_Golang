package main

import (
	"fmt"
	"math/rand"
	"time"
	"vivarium/enums"
	"vivarium/organisme"
	"vivarium/terrain"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 创建一个 Organisme 接口类型的切片
	var organismes []organisme.Organisme
	t := terrain.NewTerrain(10, 10)

	// 创建一些 Insecte 和 Plante 实例
	snail := organisme.NewInsecte(001, 0, 3, 3, 2, 1, 7, 0, 3, enums.Hermaphrodite, enums.Escargot, 30*time.Second, false)
	t.AddOrganism(snail.ID(), "Escargot", snail.GetPosX(), snail.GetPosY())
	// plante := organisme.NewPlante(002,0,5,4,0,2,9,5,enums.Graine)
	littleSnake := organisme.NewInsecte(007, 0, 5, 5, 4, 1, 7, 0, 3, enums.Femelle, enums.PetitSerpent, time.Minute, false)
	t.AddOrganism(littleSnake.ID(), "PetitSerpent", littleSnake.GetPosX(), littleSnake.GetPosY())

	// 将它们添加到 organismes 切片中
	organismes = append(organismes, snail)
	// organismes = append(organismes, plante)
	organismes = append(organismes, littleSnake)

	// TEST MANGER
	fmt.Println(t.Grid)
	littleSnake.Manger(organismes, t)
	fmt.Println(t.Grid)

	// // 遍历 organismes 切片
	// for _, orga := range organismes {
	// 	// 可以调用 Organisme 接口中定义的任何方法
	// 	orga.Vieillir()
	// 	// organisme.Mourir()

	// 	// 如果需要特定类型的方法或属性，需要类型断言
	// 	if ins, ok := orga.(*organisme.Insecte); ok {
	// 		// 现在 ins 是 *Insecte 类型，可以访问 Insecte 的特定字段和方法
	// 		fmt.Println(ins.Espece)
	// 		fmt.Println(ins.GetAge())
	// 		fmt.Println(ins.GetPosX())
	// 		fmt.Println(ins.GetPosY())
	// 		fmt.Println(t.Grid)
	// 		ins.SeDeplacer(t)
	// 		fmt.Println(ins.GetPosX())
	// 		fmt.Println(ins.GetPosY())
	// 		fmt.Println(t.Grid)
	// 		//================================================
	// 	}

	// 	// if pl, ok := organisme.(*organisme.Plante); ok {
	// 	//     // 现在 pl 是 *Plante 类型，可以访问 Plante 的特定字段和方法
	// 	//     fmt.Println(pl.EtatSante)
	// 	// }
	// }

}
