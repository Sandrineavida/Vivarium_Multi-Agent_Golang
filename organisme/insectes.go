package organisme

import (
	"fmt"
	"math"
	"math/rand"
	"time"
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
	CapaciteReproduction int
	NiveauFaim           int
	PeriodReproduire     time.Duration
	EnvieReproduire      bool
	ListePourManger      []string
	Hierarchie           int
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
func NewInsecte(organismeID int, age, posX, posY, rayon, vitesse, energie, capaciteReproduction, niveauFaim int,
	sexe enums.Sexe, espece enums.MyEspece, periodReproduire time.Duration, envieReproduire bool) *Insecte {

	// 如果映射中存在物种的层级，则使用它；否则默认为 0
	hierarchie, ok := hierarchyMap[espece]
	if !ok {
		hierarchie = 0
	}

	insecte := &Insecte{
		BaseOrganisme:        NewBaseOrganisme(organismeID, age, posX, posY, rayon, espece),
		Sexe:                 sexe,
		Vitesse:              vitesse,
		Energie:              energie,
		CapaciteReproduction: capaciteReproduction,
		NiveauFaim:           niveauFaim,
		PeriodReproduire:     periodReproduire,
		EnvieReproduire:      envieReproduire,
		ListePourManger:      foodMap[espece], // Assign the diet based on the species
		Hierarchie:           hierarchie,
	}

	return insecte
}

// Other methods (Manger, SeBattre, SeReproduire, SeDeplacer) need to be implemented here.

// SeDeplacer updates the insect's position within the terrain boundaries.
func (in *Insecte) SeDeplacer(t *terrain.Terrain) {
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

	in.Energie = max(0, min(10, in.Energie-1))
	in.NiveauFaim = max(0, min(10, in.NiveauFaim+1))

	fmt.Println(in.ID(), " : ", in.Energie, "and", in.NiveauFaim)
}

func (in Insecte) EstFaim() bool {
	return in.NiveauFaim < 5
}

// ============================================= Manger =======================================================

// getTarget 寻找周围可吃的最近目标
func getTarget(in *Insecte, organismes []Organisme) Organisme {
	var closestTarget Organisme
	minDistance := math.MaxFloat64

	for _, o := range organismes {
		if isEdible(in, o) {
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
	normalizedVitesse := float64(in.Vitesse) / 5.0       // Vitesse 范围是 1-5
	normalizedEnergie := float64(in.Energie) / 10.0      // Energie 范围是 1-10
	normalizedHierarchie := float64(in.Hierarchie) / 2.0 // Hierarchie 范围是 1-2, 因为只考虑昆虫

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
	// 获取周围的生物
	target := getTarget(in, organismes) // getTarget 需要根据您的逻辑实现

	// 如果没有找到目标，则直接退出函数
	if target == nil {
		fmt.Println("Je n'ai rien trouvé à manger")
		return nil
	}

	if targetInsecte, ok := target.(*Plante); ok {
		// 处理植物的情况
		targetInsecte.Mourir(t)
		in.Energie = max(0, min(10, in.Energie+1))
		in.NiveauFaim = max(0, min(10, in.NiveauFaim-1))
		fmt.Println(in.ID(), "Manger Plante", targetInsecte.GetEspece().String(), targetInsecte.GetID())
		return targetInsecte
	}

	if targetInsecte, ok := target.(*Insecte); ok {
		// 处理昆虫的情况
		predatorScore := calculateScore(in)
		preyScore := calculateScore(targetInsecte)

		fmt.Println("Essayer de Manger Insecte", targetInsecte.GetEspece().String())

		if predatorScore > preyScore {
			// 捕食成功
			targetInsecte.Mourir(t)
			in.Energie = max(0, min(10, in.Energie+1))
			in.NiveauFaim = max(0, min(10, in.NiveauFaim-1))

			fmt.Println("Success!!!! Manger Insecte", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				" Score: predator = ", predatorScore, "prey = ", preyScore)
			return targetInsecte
		} else {
			// 捕食失败；逃跑 and 离开
			fmt.Println("Fail!!!!!!!!!! Manger Insecte", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
				" Score: predator = ", predatorScore, "prey = ", preyScore)
			n := rand.Intn(3) + 1 // 让二者分别SeDeplace1-3次
			for i := 0; i < n; i++ {
				in.SeDeplacer(t)
				targetInsecte.SeDeplacer(t)
			}
			return nil
		}
		return nil
	}
	return nil
}

// ============================================= END of Manger =======================================================

func (in *Insecte) CheckEtat(t *terrain.Terrain) Organisme {
	// 检查能量和饥饿水平是否达到极限
	if in.Energie <= 0 || in.NiveauFaim >= 10 {
		fmt.Println("Insecte", in.ID(), "est mort de faim ou de fatigue.")
		in.Mourir(t)
		return in
	}
	return nil
}

/* Helper function */
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
