package organisme

import (
	"fmt"
	"math"
	"math/rand"
	"vivarium/enums"
	"vivarium/terrain"
	"vivarium/utils"
)

// Insecte represents an insect and embeds BaseOrganisme to inherit its properties.
type Insecte struct {
	*BaseOrganisme
	Sexe                 enums.Sexe
	Vitesse              int
	Energie              int
	PeriodReproduire     int
	EnvieReproduire      bool
	ListePourManger      []string
	Hierarchie           int
	AgeGaveBirthLastTime int
}

// foodMap defines what each insect species can eat.
var foodMap = map[enums.MyEspece][]string{
	enums.Escargot:         {"PetitHerbe", "GrandHerbe", "Champignon"},
	enums.Grillons:         {"Champignon"},
	enums.Lombric:          {"PetitHerbe", "GrandHerbe"},
	enums.PetitSerpent:     {"Lombric", "Escargot"},
	enums.AraignéeSauteuse: {"Grillons"},
}

/*
Escargot -> Petit herbe/Grande herbe/Champignon
Grillons (se reproduire hyper vite?) -> Champignon
Lombric (3) -> Petit herbe/Grande herbe
Petit Serpent (1) -> Lombric
Araignée sauteuse (seulement 2) -> Grillons
*/

// hierarchyMap defines the hierarchy level for each insect species.
var hierarchyMap = map[enums.MyEspece]int{
	enums.Escargot:         1,
	enums.Grillons:         1,
	enums.Lombric:          1,
	enums.PetitSerpent:     2,
	enums.AraignéeSauteuse: 2,
	// 默认情况，例如 PetitHerbe, GrandHerbe, Champignon
} // Hierarchie: PetitHerbe, GrandHerbe, Champignon=0 < Escargot = Grillons = Lombric = 1 < AraignéeSauteuse = PetitSerpent = 2

// NewInsecte creates a new Insecte with the given attributes.
func NewInsecte(organismeID int, age, posX, posY, vitesse, energie int,
	sexe enums.Sexe, espece enums.MyEspece, envieReproduire bool) *Insecte {

	attributes := enums.SpeciesAttributes[espece]
	attributesInsecte := enums.InsectAttributesMap[espece]

	// 如果映射中存在物种的层级，则使用它；否则默认为 0
	hierarchie, ok := hierarchyMap[espece]
	if !ok {
		hierarchie = 0
	}

	insecte := &Insecte{
		BaseOrganisme: NewBaseOrganisme(organismeID, age, posX, posY, attributesInsecte.Rayon, espece,
			attributes.AgeRate, attributes.MaxAge, attributes.GrownUpAge, attributes.TooOldToReproduceAge, attributes.NbProgeniture),
		Sexe:    sexe,
		Vitesse: vitesse,
		// Energie:              energie,
		Energie:              attributes.NiveauEnergie,
		PeriodReproduire:     attributesInsecte.PeriodReproduire,
		EnvieReproduire:      envieReproduire,
		ListePourManger:      foodMap[espece], // Assign the diet based on the species
		Hierarchie:           hierarchie,
		AgeGaveBirthLastTime: 0,
	}

	return insecte
}

// Other methods (Manger √, SeBattre √, SeReproduire, SeDeplacer √) need to be implemented here.

// SeDeplacer updates the insect's position within the terrain boundaries.
func (in *Insecte) SeDeplacer(t *terrain.Terrain) {
	// 检查是否忙碌
	if in.Busy {
		fmt.Println("Insecte", in.GetID(), "is busy, cannot move")
		return
	}

	in.Busy = true
	defer func() { in.Busy = false }() // 行为完成后重置状态

	// Generate random movement direction
	deltaX := rand.Intn(3) - 1 // Random int in {-1, 0, 1}
	deltaY := rand.Intn(3) - 1 // Random int in {-1, 0, 1}

	// Apply velocity and constrain the new position within the terrain boundaries
	newX := utils.Intmax(0, utils.Intmin(in.PositionX+deltaX*in.Vitesse, t.Width-1))
	newY := utils.Intmax(0, utils.Intmin(in.PositionY+deltaY*in.Vitesse, t.Length-1))

	// Update the insect's position in the Terrain and Insecte
	t.UpdateOrganismPosition(in.OrganismeID, in.Espece.String(), in.PositionX, in.PositionY, newX, newY)
	in.PositionX = newX
	in.PositionY = newY

	attributes := enums.SpeciesAttributes[in.Espece]
	in.Energie = utils.Intmax(0, utils.Intmin(attributes.NiveauEnergie, in.Energie-1))

	//fmt.Println(in.GetID(), " : ", in.Energie)
}

func (in Insecte) AFaim() bool {
	attributes := enums.SpeciesAttributes[in.Espece]
	return in.Energie < attributes.NiveauEnergie/3*2
}

// ============================================= getTarget =======================================================
// getTarget 寻找周围最近的符合需求的目标
func getTarget(in *Insecte, organismes []Organisme, jud_func func(*Insecte, Organisme) bool) Organisme {
	var closestTarget Organisme
	minDistance := math.MaxFloat64

	for _, o := range organismes {
		if jud_func(in, o) {
			x, y := o.GetPosX(), o.GetPosY()
			distance := distance(in.PositionX, in.PositionY, x, y)

			if distance <= float64(in.Rayon) && distance < minDistance {
				closestTarget = o
				minDistance = distance
			}
		}
	}

	return closestTarget
}

// ============================================= END of getTarget =======================================================

// ============================================= Manger =======================================================
// isEdible 检查是否为可食用目标
func isEdible(in *Insecte, target Organisme) bool {
	for _, food := range in.ListePourManger {
		if target.GetEspece().String() == food {
			return true
		}
	}
	return false
}

// distance 计算两点之间的距离
func distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

// calculateScore 计算捕食者和猎物的分数
func calculateScore(in *Insecte) float64 {
	// 归一化属性值
	normalizedVitesse := float64(in.Vitesse) / 5.0 // Vitesse 范围是 1-5
	attributes := enums.SpeciesAttributes[in.Espece]
	MaxEnergie := attributes.NiveauEnergie
	normalizedEnergie := float64(in.Energie) / float64(MaxEnergie) // Energie 范围是 1-MaxEnergie
	normalizedHierarchie := float64(in.Hierarchie) / 2.0           // Hierarchie 范围是 1-2, 因为只考虑昆虫

	// 设置权重
	w1, w2, w3 := 1.0, 2.0, 3.0

	// 计算加权平均值
	score := (w1*normalizedVitesse + w2*normalizedEnergie + w3*normalizedHierarchie) / (w1 + w2 + w3)

	// 添加随机幸运值
	luck := rand.Float64()*0.6 - 0.3 // 在 -0.3 到 0.3 之间的随机数
	finalScore := score + luck

	return finalScore
}

func (in *Insecte) Manger(organismes []Organisme, t *terrain.Terrain) Organisme {
	// 检查是否忙碌
	if in.Busy {
		fmt.Println("Insecte", in.GetID(), "is busy, cannot eat")
		return nil
	}

	// 设置忙碌状态
	in.Busy = true
	defer func() { in.Busy = false }() // 行为完成后重置状态

	// 获取周围的生物
	target := getTarget(in, organismes, isEdible)

	// 如果没有找到目标，则直接退出函数
	if target == nil {
		fmt.Println("Je n'ai rien trouvé à manger")
		return nil
	}

	if targetPlante, ok := target.(*Plante); ok {
		// 处理植物的情况
		targetPlante.Mourir(t)
		attributes := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes.NiveauEnergie, in.Energie+1))
		fmt.Println(in.GetID(), "Manger Plante", targetPlante.GetEspece().String(), targetPlante.GetID())
		return targetPlante
	}

	if targetInsecte, ok := target.(*Insecte); ok {
		// 处理昆虫的情况
		targetInsecte.Busy = true

		predatorScore := calculateScore(in)
		preyScore := calculateScore(targetInsecte)

		fmt.Println("Essayer de Manger Insecte", targetInsecte.GetEspece().String())

		if predatorScore > preyScore {
			// 捕食成功
			targetInsecte.Mourir(t)
			attributes := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes.NiveauEnergie, in.Energie+1))

			fmt.Println("Success!!!! Manger Insecte", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				" Score: predator = ", predatorScore, "prey = ", preyScore)
			return targetInsecte
		} else {
			// 捕食失败；逃跑 and 离开
			fmt.Println("Fail!!!!!!!!!! Manger Insecte", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				" Score: predator = ", predatorScore, "prey = ", preyScore)
			n := rand.Intn(2) + 1 // 让二者分别SeDeplace1-2次
			for i := 0; i < n; i++ {
				in.SeDeplacer(t)
				targetInsecte.SeDeplacer(t)
			}
			targetInsecte.Busy = false
			return nil
		}
	}
	return nil
}

// ============================================= END of Manger =======================================================

// ============================================= CheckEtat =======================================================
func (in *Insecte) CheckEtat(t *terrain.Terrain) Organisme {
	// 检查能量和饥饿水平是否达到极限
	if in.Energie <= 0 {
		fmt.Println("Insecte", in.GetID(), "est mort de faim ou de fatigue.")
		in.Mourir(t)
		return in
	}
	return nil
}

// ============================================= END of CheckEtat =======================================================

// ============================================= SeBattre =======================================================
func isFightable(in *Insecte, target Organisme) bool {
	// 目前定义是只有同一种昆虫才能斗殴
	return in.Espece == target.GetEspece()
}

// 随便找人打，不知道能不能找到能打的对象
func (in *Insecte) SeBattreRandom(organismes []Organisme, t *terrain.Terrain) {
	// 检查是否忙碌
	if in.Busy {
		// 可能在另一个insect那边已经主动和当前insect打起来了；可能在吃，目前设定在吃就不打架；可能正在交配，目前设定在交配就不打架
		fmt.Println("Insecte", in.GetID(), "is busy, cannot fight")
		return
	}

	// 获取周围的生物
	target := getTarget(in, organismes, isFightable)

	// 如果没有找到目标，则直接退出函数
	if target == nil {
		fmt.Println("Damn can't find anyone to fight gonna explode dude.")
		return
		// return nil
	}

	in.Busy = true
	defer func() { in.Busy = false }() // 行为完成后重置状态

	if targetInsecte, ok := target.(*Insecte); ok {
		if targetInsecte.Busy {
			// 如果忙碌，回退或延迟操作 （不打群架之类的，估计有bug，后边再说）
			fmt.Println("Insecte", targetInsecte.GetID(), "is busy, cannot fight")
			return
		}
		targetInsecte.Busy = true
		defer func() { targetInsecte.Busy = false }() // 行为完成后重置状态

		fighterScore := calculateScore(in)
		victimScore := calculateScore(targetInsecte)

		fmt.Println("Essayer de SeBattreRandom Insecte", targetInsecte.GetEspece().String())

		if fighterScore > victimScore {
			// 干赢了
			attributes_target := enums.SpeciesAttributes[targetInsecte.Espece]
			targetInsecte.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, targetInsecte.Energie-3))
			attributes_in := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-1))

			fmt.Println("BEAT THE SHIT OUT OF ", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				" !!! Score: fighter = ", fighterScore, "victim = ", victimScore)
			return
			// return targetInsecte
		} else {
			// 被干爆
			attributes_target := enums.SpeciesAttributes[targetInsecte.Espece]
			targetInsecte.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, targetInsecte.Energie-1))
			attributes_in := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-3))
			fmt.Println("Damn it I get fked up by", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				"... Score: fighter = ", fighterScore, "victim = ", victimScore)
			// n := rand.Intn(3) + 1 // 让二者分别SeDeplace1-3次
			// for i := 0; i < n; i++ {
			// 	in.SeDeplacer(t)
			// 	targetInsecte.SeDeplacer(t)
			// }
			// return nil
		}
	}
	// return nil
}

// 传入"确定的"能打的对象 (目前只在繁殖中使用)
func (in *Insecte) SeBattre(target *Insecte, t *terrain.Terrain) {
	// 检查是否忙碌
	if in.Busy {
		// 可能在另一个insect那边已经主动和当前insect打起来了；可能在吃，目前设定在吃就不打架；可能正在交配，目前设定在交配就不打架
		fmt.Println("Insecte", in.GetID(), "is busy, cannot fight")
		return
	}

	// 如果没有找到目标，则直接退出函数
	if target == nil {
		fmt.Println("Error: target is nil (SeBattre)")
		return
		// return nil
	}

	if target.Busy {
		fmt.Println("Target insect", target.GetID(), "is busy, cannot fight")
		return
	}

	in.Busy = true
	defer func() { in.Busy = false }() // 行为完成后重置状态
	target.Busy = true
	defer func() { target.Busy = false }() // 行为完成后重置状态

	fighterScore := calculateScore(in)
	victimScore := calculateScore(target)

	fmt.Println("Essayer de SeBattre Insecte", target.GetEspece().String())

	if fighterScore > victimScore {
		// 干赢了
		attributes_target := enums.SpeciesAttributes[target.Espece]
		target.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, target.Energie-3))
		attributes_in := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-1))

		fmt.Println("EAT THE SHIT OUT OF", target.GetEspece().String(), target.GetID(),
			" !!! Score: fighter = ", fighterScore, "victim = ", victimScore)
		return
	} else {
		// 干赢了
		attributes_target := enums.SpeciesAttributes[target.Espece]
		target.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, target.Energie-1))
		attributes_in := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-3))
		fmt.Println("Damn it I get fked up by", target.GetEspece().String(), target.GetID(),
			"... Score: fighter = ", fighterScore, "victim = ", victimScore)
	}
}

// ============================================= End of SeBattre =======================================================

// ============================================= SeReproduire =======================================================

// 目前没有限制能不能乱伦

// 判断是否可以繁殖
//  1. 是否有足够的能量
//  2. 是否有足够的饥饿水平
//  3. PeriodReproduire是否已过 (Age - AgeGaveBirthLastTime >= PeriodReproduire)
//     （4. 是否有足够的空间；这个可以暂时先不考虑）
//     如果可以的话，就把EnvieReproduire设置为true
//
// 注：
// （这个函数应该在main里被不断地调用？）
// （func SeReproduire应该直接通过EnvieReproduire来进行判断，而不是这个函数）
func (in *Insecte) AvoirEnvieReproduire() {
	//fmt.Println("ID:", in.GetID(), "Energie:", in.Energie, "Age:", in.Age, "上次bang:", in.AgeGaveBirthLastTime, "bang周期:", in.PeriodReproduire, "成年:", in.GrownUpAge, "老了:", in.TooOldToReproduceAge)
	//if in.Energie >= 5 {
	//	fmt.Println("1:能量够bang")
	//}
	//if in.Age-in.AgeGaveBirthLastTime >= in.PeriodReproduire {
	//	fmt.Println("2:恢复了，又可以bang了")
	//}
	//if in.Age >= in.GrownUpAge {
	//	fmt.Println("3:成年了，可以bang了")
	//}
	//if in.Age <= in.TooOldToReproduceAge {
	//	fmt.Println("4:没老，没养胃")
	//}
	if (!in.AFaim()) && in.Age-in.AgeGaveBirthLastTime >= in.PeriodReproduire && in.Age >= in.GrownUpAge && in.Age <= in.TooOldToReproduceAge {
		in.EnvieReproduire = true
	} else {
		in.EnvieReproduire = false
	}
}

// isReproducible 判断是否可以繁殖
func isReproducible(in *Insecte, target Organisme) bool {
	// 目前定义是只有同一种昆虫才能繁殖
	// return in.Espece == target.GetEspece()
	if targetInsecte, ok := target.(*Insecte); ok {
		return in.Espece == targetInsecte.Espece && targetInsecte.EnvieReproduire
	}
	return false //一般不会走到这里，因为肯定是虫子不是植物
}

// SeReproduire 通过EnvieReproduire来判断是否繁殖
func (in *Insecte) SeReproduire(organismes []Organisme, t *terrain.Terrain) (int, []Organisme) {
	// 检查是否忙碌
	if in.Busy {
		fmt.Println("Insecte", in.GetID(), "is busy, cannot reproduce")
		return 0, nil
	}

	// 先判断本insect是否有繁殖的欲望
	if !in.EnvieReproduire {
		fmt.Println("I don't wanna reproduce yet...")
		return 0, nil
	} else {
		fmt.Println(in.GetID(), "好想bang啊")
	}

	// 在周围的生物里找能干的
	target := getTarget(in, organismes, isReproducible)
	//fmt.Println("操操操操操 能干的生物：", target)

	// 如果没有找到目标，则直接退出函数
	if target == nil {
		fmt.Println("Can't find anyone to bang with...!!")
		return 0, nil
	}

	if targetInsecte, ok := target.(*Insecte); ok {
		if in.Sexe == enums.Hermaphrodite {
			targetInsecte = findHermaphroditeTarget(targetInsecte, t)
		} else {
			targetInsecte = findBisexualTarget(in, targetInsecte, t)
		}

		if targetInsecte == nil {
			fmt.Println("找不到合适的繁殖对象")
			return 0, nil
		}

		if targetInsecte.Busy {
			fmt.Println("Target insect", targetInsecte.GetID(), "is busy, cannot reproduce")
			return 0, nil
		}

		in.Busy = true
		defer func() { in.Busy = false }() // 行为完成后重置状态
		targetInsecte.Busy = true
		defer func() { targetInsecte.Busy = false }() // 行为完成后重置状态

		// 开生！！！！！！！！！！！！
		var sliceNewBorn []Organisme

		in.AgeGaveBirthLastTime = in.Age
		targetInsecte.AgeGaveBirthLastTime = targetInsecte.Age
		in.EnvieReproduire = false
		targetInsecte.EnvieReproduire = false
		in.Energie = utils.Intmax(0, in.Energie-1)
		targetInsecte.Energie = utils.Intmax(0, targetInsecte.Energie-1)

		for i := 0; i < in.NbProgeniture; i++ {
			// 生成新生昆虫
			// 随机选择 in 和 target 的 positionX 和 positionY
			newX := in.PositionX
			newY := in.PositionY
			if rand.Intn(2) == 0 { // 随机数为0或1，决定使用 in 的位置还是 target 的位置
				newX = target.GetPosX()
				newY = target.GetPosY()
			}
			newSexe := in.Sexe
			if rand.Intn(2) == 0 { // 随机数为0或1，决定使用 in 的性别还是 target 的性别 (就算雌雄同体的也随机选它一个，虽然没必要但懒得判断了)
				newSexe = targetInsecte.Sexe
			}
			newBorn := NewInsecte(-1, 0, newX, newY, in.Vitesse, 10, newSexe, in.Espece, false) // ID为-1;；要去main里面更新terrain和organismes的list
			fmt.Println("生出来了！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！")
			sliceNewBorn = append(sliceNewBorn, newBorn)
		}

		return in.NbProgeniture, sliceNewBorn

	}

	return 0, nil

}

func findHermaphroditeTarget(targetInsecte *Insecte, t *terrain.Terrain) *Insecte {
	// 	== 直接判断对方想不想bang（target的EnvieReproduire）
	// 	 ++++ 想就返回对方对象
	//   ++++ 不想就返回nil
	if targetInsecte.EnvieReproduire {
		return targetInsecte
	}
	return nil
}

func findBisexualTarget(in *Insecte, targetInsecte *Insecte, t *terrain.Terrain) *Insecte {
	// 	== 如果遇到同性
	//    +++++ SeBattre
	//    +++++ 返回nil
	//  == 如果遇到异性
	//    +++++ 判断对方想不想bang （target的EnvieReproduire）
	//      ++++ 想就返回对方对象
	//      ++++ 不想就返回nil
	if in.Sexe == targetInsecte.Sexe {
		// 同性
		in.SeBattre(targetInsecte, t)
		return nil
	} else {
		// 异性
		if targetInsecte.EnvieReproduire {
			return targetInsecte
		}
		return nil
	}
}

// 先不考虑卵生胎生这种东西了

// ============================================= End of SeReproduire =======================================================
