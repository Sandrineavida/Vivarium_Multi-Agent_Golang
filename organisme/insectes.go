package organisme

import (
	"math"
	"math/rand"
	"time"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/terrain"
	"vivarium/utils"
)

const (
	timeSleep = 1000
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

	// should be used in ebiten sprites
	IsManger     bool
	IsReproduire bool
	IsSeDeplacer bool
	IsSeBattre   bool
	IsWinner     bool
	IsLooser     bool
	IsNormal     bool
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
	// PetitHerbe, GrandHerbe, Champignon all have 0 as hierarchy level
} // Hierarchie: PetitHerbe, GrandHerbe, Champignon=0 < Escargot = Grillons = Lombric = 1 < AraignéeSauteuse = PetitSerpent = 2

// NewInsecte creates a new Insecte with the given attributes.
func NewInsecte(organismeID, age, posX, posY int, sexe enums.Sexe, espece enums.MyEspece, envieReproduire bool) *Insecte {

	attributes := enums.SpeciesAttributes[espece]
	attributesInsecte := enums.InsectAttributesMap[espece]
	vitesse := enums.InsectSpeeds[espece] // 从映射中获取速度

	hierarchie, ok := hierarchyMap[espece]
	if !ok {
		hierarchie = 0
	}

	insecte := &Insecte{
		BaseOrganisme: NewBaseOrganisme(organismeID, age, posX, posY, attributesInsecte.Rayon, espece,
			attributes.AgeRate, attributes.MaxAge, attributes.GrownUpAge, attributes.TooOldToReproduceAge, attributes.NbProgeniture, true),
		Sexe:                 sexe,
		Vitesse:              vitesse,
		Energie:              attributes.NiveauEnergie,
		PeriodReproduire:     attributesInsecte.PeriodReproduire,
		EnvieReproduire:      envieReproduire,
		ListePourManger:      foodMap[espece], // Assign the diet based on the species
		Hierarchie:           hierarchie,
		AgeGaveBirthLastTime: 0,

		IsManger:     false, // The default value is false (can omit it)
		IsReproduire: false,
		IsSeDeplacer: false,
		IsSeBattre:   false,
		IsWinner:     false,
		IsLooser:     false,
		IsNormal:     true,
	}

	return insecte
}

// SeDeplacer updates the insect's position within the terrain boundaries.
func (in *Insecte) SeDeplacer(t *terrain.Terrain) {
	// check if the insect is busy
	if in.Busy { // if busy, means it's already doing something, cannot move
		return
	}

	in.Busy = true
	in.IsSeDeplacer = true
	in.IsNormal = false
	defer func() {
		in.Busy = false
		in.IsSeDeplacer = false
		in.IsNormal = true
	}() // reset status after the action is done

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
}

func (in Insecte) AFaim() bool {
	attributes := enums.SpeciesAttributes[in.Espece]
	return in.Energie < attributes.NiveauEnergie/3*2
}

// ============================================= getTarget =======================================================
// getTarget: find the closest target that satisfies the condition
func getTarget(in *Insecte, organismes []Organisme, jud_func func(*Insecte, Organisme) bool) Organisme {
	var closestTarget Organisme
	minDistance := math.MaxFloat64

	for _, o := range organismes {
		if jud_func(in, o) {
			x, y := o.GetPosX(), o.GetPosY()
			distance := utils.Calcul_Distance(in.PositionX, in.PositionY, x, y)

			if distance <= float64(in.Rayon) && distance < minDistance {
				if !o.GetEtat() { // if the target is not dead, then it's a valid possible target
					closestTarget = o
					minDistance = distance
				}
			}
		}
	}

	return closestTarget
}

// ============================================= END of getTarget =======================================================

// ============================================= Manger =======================================================
// isEdible: check if the target is edible
func isEdible(in *Insecte, target Organisme) bool {
	for _, food := range in.ListePourManger {
		if target.GetEspece().String() == food {
			return true
		}
	}
	return false
}

// calculateScore: calculate the score of the predator and prey
func calculateScore(in *Insecte) float64 {
	// normalise the attributes
	normalizedVitesse := float64(in.Vitesse) / 5.0 // range of Vitesse: 1-5
	attributes := enums.SpeciesAttributes[in.Espece]
	MaxEnergie := attributes.NiveauEnergie
	normalizedEnergie := float64(in.Energie) / float64(MaxEnergie) // range of Energie: 1-MaxEnergie
	normalizedHierarchie := float64(in.Hierarchie) / 2.0           // range of Hierarchie: 1-2 (considering only insects)

	// set the weights
	w1, w2, w3 := 1.0, 2.0, 3.0

	// calculate the weighted average
	score := (w1*normalizedVitesse + w2*normalizedEnergie + w3*normalizedHierarchie) / (w1 + w2 + w3)

	// introduce some randomness as luck
	luck := rand.Float64()*0.6 - 0.3 // between -0.3 and 0.3
	finalScore := score + luck

	return finalScore
}

func (in *Insecte) Manger(organismes []Organisme, t *terrain.Terrain) {
	// check if the insect is busy
	if in.Busy { // if busy, means it's already doing something, cannot eat
		return
	}

	in.Busy = true
	defer func() {
		time.Sleep(timeSleep * time.Millisecond)
		in.Busy = false
		in.IsManger = false
		in.IsNormal = true
	}()

	// percepect the environment to find the closest target that is edible
	target := getTarget(in, organismes, isEdible)

	// if no target found, return
	if target == nil {
		//fmt.Println("Je n'ai rien trouvé à manger")
		return
	}

	if targetPlante, ok := target.(*Plante); ok {
		// Handling the case of plants as targets
		//if it's a GrandHerbe, then it can be eaten bit by bit; otherwise it's a whole meal
		if targetPlante.Espece == enums.GrandHerbe {
			if targetPlante.IsBeingEaten {
				// les insects sont sociables!!! Ils ne mangent pas les mêmes plantes en même temps ;)
				return
			} else {
				//when it's time to eat, teleport to the position of the plant
				in.PositionX = targetPlante.PositionX
				in.PositionY = targetPlante.PositionY

				targetPlante.IsBeingEaten = true
				// only set bools to true when it's really eating
				in.IsManger = true
				in.IsNormal = false
				defer func() {
					targetPlante.IsBeingEaten = false
				}()
				targetPlante.NbParts -= 1
				if targetPlante.NbParts == 0 {
					// if the GrandHerbe is eaten up, then it dies
					targetPlante.Mourir(t)
				}
			}
		} else {
			// when it's time to eat, teleport to the position of the PetitHerbe
			in.PositionX = targetPlante.PositionX
			in.PositionY = targetPlante.PositionY
			targetPlante.Mourir(t)
		}

		attributes := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes.NiveauEnergie, in.Energie+5))
		return
	}

	if targetInsecte, ok := target.(*Insecte); ok {
		// Handling the case of insects as targets
		targetInsecte.Busy = true

		// When starting to eat, move the predator to the position of the prey
		in.PositionX = targetInsecte.PositionX
		in.PositionY = targetInsecte.PositionY

		predatorScore := calculateScore(in)
		preyScore := calculateScore(targetInsecte)

		//fmt.Println("Essayer de Manger Insecte", targetInsecte.GetEspece().String())

		if predatorScore > preyScore {
			// Predator succeeds to catch the prey
			// only set bools to true when it's really eating
			in.IsManger = true
			in.IsNormal = false
			targetInsecte.Mourir(t)
			attributes := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes.NiveauEnergie, in.Energie+10))
			return
		} else {
			// Predator fails to catch the prey; both of them move away
			n := rand.Intn(2) + 1 // let them randomly SeDeplace 1-2 times
			for i := 0; i < n; i++ {
				in.SeDeplacer(t)
				targetInsecte.SeDeplacer(t)
			}
			targetInsecte.Busy = false
			return
		}
	}

	return
}

// ============================================= END of Manger =======================================================

// ============================================= CheckEtat =======================================================
func (in *Insecte) CheckEtat(t *terrain.Terrain) Organisme {
	// Check if the Energy is equal to 0
	if in.Energie <= 0 {
		in.Mourir(t) // if yes, then the insect dies
		return in
	}
	return nil
}

// ============================================= END of CheckEtat =======================================================

// ============================================= SeBattre =======================================================
func isFightable(in *Insecte, target Organisme) bool {
	// Right now, only insects of the same species can fight
	return in.Espece == target.GetEspece()
}

// Find a random target to fight with
func (in *Insecte) SeBattreRandom(organismes []Organisme, t *terrain.Terrain) {
	// check if the insect is busy
	if in.Busy {
		// maybe it's already fighting with another insect; maybe it's eating, currently set to not fight when eating; maybe it's mating, currently set to not fight when mating
		return
	}

	// Find the closest target that is fightable
	target := getTarget(in, organismes, isFightable)

	if target == nil {
		//fmt.Println("Damn can't find anyone to fight gonna explode dude.")
		return
	}

	in.Busy = true
	in.IsSeBattre = true
	in.IsNormal = false
	defer func() {
		time.Sleep(timeSleep * time.Millisecond)
		in.Busy = false
		in.IsSeBattre = false
		in.IsNormal = true
		in.IsWinner = false
		in.IsLooser = false
	}()

	if targetInsecte, ok := target.(*Insecte); ok {
		if targetInsecte.Busy {
			return
		}
		targetInsecte.Busy = true
		targetInsecte.IsSeBattre = true
		targetInsecte.IsNormal = false
		defer func() {
			time.Sleep(timeSleep * time.Millisecond)
			targetInsecte.Busy = false
			targetInsecte.IsSeBattre = false
			targetInsecte.IsNormal = true
			targetInsecte.IsWinner = false
			targetInsecte.IsLooser = false
		}()

		fighterScore := calculateScore(in)
		victimScore := calculateScore(targetInsecte)

		if fighterScore > victimScore {
			// Win
			attributes_target := enums.SpeciesAttributes[targetInsecte.Espece]
			targetInsecte.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, targetInsecte.Energie-3))
			attributes_in := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-1))
			/* 			fmt.Println("BEAT THE SHIT OUT OF ", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
			" !!! Score: fighter = ", fighterScore, "victim = ", victimScore) */
			in.IsWinner = true
			targetInsecte.IsLooser = true

			return
			// return targetInsecte
		} else {
			// Lose
			attributes_target := enums.SpeciesAttributes[targetInsecte.Espece]
			targetInsecte.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, targetInsecte.Energie-1))
			attributes_in := enums.SpeciesAttributes[in.Espece]
			in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-3))
			/* 			fmt.Println("Damn it I get fked up by", targetInsecte.GetEspece().String(), targetInsecte.GetID(),
			"... Score: fighter = ", fighterScore, "victim = ", victimScore) */
			targetInsecte.IsWinner = true
			in.IsLooser = true
		}
	}
	// return nil
}

// Fight against a specific target (currently only used in SeReproduire)
func (in *Insecte) SeBattre(target *Insecte, t *terrain.Terrain) bool {
	// check if the insect is busy
	if in.Busy {
		return false
	}

	if target == nil {
		return false
	}

	if target.Busy {
		return false
	}

	// Move the insect to the position of the target
	in.PositionX = target.PositionX
	in.PositionY = target.PositionY

	in.IsSeBattre = true
	in.IsNormal = false
	defer func() {
		time.Sleep(2 * timeSleep * time.Millisecond)
		in.Busy = false
		in.IsSeBattre = false
		in.IsNormal = true
		in.IsWinner = false
		in.IsLooser = false
	}()
	target.Busy = true
	target.IsSeBattre = true
	target.IsNormal = false
	defer func() {
		time.Sleep(2 * timeSleep * time.Millisecond)
		target.Busy = false
		target.IsSeBattre = false
		target.IsNormal = true
		target.IsWinner = false
		target.IsLooser = false
	}()

	fighterScore := calculateScore(in)
	victimScore := calculateScore(target)

	if fighterScore > victimScore {
		// win
		attributes_target := enums.SpeciesAttributes[target.Espece]
		target.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, target.Energie-3))
		attributes_in := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-1))
		/* 		fmt.Println("EAT THE SHIT OUT OF", target.GetEspece().String(), target.GetID(),
		" !!! Score: fighter = ", fighterScore, "victim = ", victimScore) */
		time.Sleep(timeSleep * time.Millisecond)
		in.IsWinner = true
		target.IsLooser = true

		return true
	} else {
		// lose
		attributes_target := enums.SpeciesAttributes[target.Espece]
		target.Energie = utils.Intmax(0, utils.Intmin(attributes_target.NiveauEnergie, target.Energie-1))
		attributes_in := enums.SpeciesAttributes[in.Espece]
		in.Energie = utils.Intmax(0, utils.Intmin(attributes_in.NiveauEnergie, in.Energie-3))
		/* 		fmt.Println("Damn it I get fked up by", target.GetEspece().String(), target.GetID(),
		"... Score: fighter = ", fighterScore, "victim = ", victimScore) */
		time.Sleep(timeSleep * time.Millisecond)
		target.IsWinner = true
		in.IsLooser = true

		return false
	}

}

// ============================================= End of SeBattre =======================================================

// ============================================= SeReproduire =======================================================
// Set the insect's EnvieReproduire to true if it satisfies the conditions
func (in *Insecte) AvoirEnvieReproduire() {
	if (!in.AFaim()) && in.Age-in.AgeGaveBirthLastTime >= in.PeriodReproduire && in.Age >= in.GrownUpAge && in.Age <= in.TooOldToReproduceAge {
		in.EnvieReproduire = true
	} else {
		in.EnvieReproduire = false
	}
}

// isReproducible: check if the target is reproducible (used in "getTarget" func.)
func isReproducible(in *Insecte, target Organisme) bool {
	// only insects of the same species and both want to reproduce can reproduce
	if targetInsecte, ok := target.(*Insecte); ok { // cast the target to Insecte
		return in.Espece == targetInsecte.Espece && targetInsecte.EnvieReproduire
	}
	return false
}

// SeReproduire
func (in *Insecte) SeReproduire(organismes []Organisme, t *terrain.Terrain) (int, []Organisme, bool) {
	// check if the insect is busy
	if in.Busy {
		return 0, nil, false
	}

	// Check if the insect wants to reproduce
	if !in.EnvieReproduire {
		//fmt.Println("I don't wanna reproduce yet...")
		return 0, nil, false
	}

	// find the closest target that is reproducible
	target := getTarget(in, organismes, isReproducible)

	// if no target found, return
	if target == nil {
		// fmt.Println("Can't find anyone to bang with...!!")
		return 0, nil, false
	}

	if targetInsecte, ok := target.(*Insecte); ok {
		findTargetAgain := false
		if in.Sexe == enums.Hermaphrodite {
			// if the insect is hermaphrodite, then it can reproduce with any insect of its kind
			targetInsecte = findHermaphroditeTarget(targetInsecte, t)
		} else {
			//
			targetInsecte, findTargetAgain = findBisexualTarget(in, targetInsecte, t)
		}

		if targetInsecte == nil {
			// can't find a suitable target to reproduce
			return 0, nil, findTargetAgain
		}

		if targetInsecte.Busy {
			// fmt.Println("Target insect", targetInsecte.GetID(), "is busy, cannot reproduce")
			return 0, nil, true
		}

		in.Busy = true
		in.IsReproduire = true
		in.IsNormal = false
		defer func() {
			time.Sleep(3 * timeSleep * time.Millisecond)
			in.Busy = false
			in.IsReproduire = false
			in.IsNormal = true
		}()
		targetInsecte.Busy = true
		targetInsecte.IsReproduire = true
		targetInsecte.IsNormal = false
		defer func() {
			time.Sleep(3 * timeSleep * time.Millisecond)
			targetInsecte.Busy = false
			targetInsecte.IsReproduire = false
			targetInsecte.IsNormal = true
		}()

		// Start to reproduce!!!!!!!!!!!!!!!!!!!!
		var sliceNewBorn []Organisme

		// Move the insect to the position of the target
		in.PositionX = targetInsecte.PositionX
		in.PositionY = targetInsecte.PositionY

		in.AgeGaveBirthLastTime = in.Age
		targetInsecte.AgeGaveBirthLastTime = targetInsecte.Age
		in.EnvieReproduire = false
		targetInsecte.EnvieReproduire = false
		in.Energie = utils.Intmax(0, in.Energie-1)
		targetInsecte.Energie = utils.Intmax(0, targetInsecte.Energie-1)

		for i := 0; i < in.NbProgeniture; i++ {
			// Give birth to new insects
			newX := in.PositionX
			newY := in.PositionY

			newSexe := in.Sexe
			if rand.Intn(2) == 0 {
				// randomly decide the sexuality of the new insect (between its parents) (even if it's hermaphrodite, still randomly choose one)
				newSexe = targetInsecte.Sexe
			}
			newBorn := NewInsecte(-1, 0, newX, newY, newSexe, in.Espece, false) // Get the real ID in server
			sliceNewBorn = append(sliceNewBorn, newBorn)
		}

		return in.NbProgeniture, sliceNewBorn, false

	}

	return 0, nil, false

}

func findHermaphroditeTarget(targetInsecte *Insecte, t *terrain.Terrain) *Insecte {
	if targetInsecte.EnvieReproduire { // Check if the target wants to reproduce
		return targetInsecte
	}
	return nil
}

func findBisexualTarget(in *Insecte, targetInsecte *Insecte, t *terrain.Terrain) (mateInsect *Insecte, findTargetAgain bool) {
	mateInsect = nil
	findTargetAgain = false
	if in.Sexe == targetInsecte.Sexe {
		// Of the same sexuallity
		result := in.SeBattre(targetInsecte, t) // Fight
		if result {                             // If win
			findTargetAgain = true // Set true in order to find another target later
		}
	} else {
		// Of different sexuallity
		if targetInsecte.EnvieReproduire { // Check if the target wants to reproduce
			mateInsect = targetInsecte
		}
	}
	return
}

// ============================================= End of SeReproduire =======================================================

// ============================================= UpdateEnergie =======================================================
// PerceptClimat: return the severity of the environment (severity: 0-100)
func (in *Insecte) PerceptClimat(climat climat.Climat) int {
	severity := 0

	// Luminaire severity
	if climat.Luminaire < 5 || climat.Luminaire > 90 {
		severity += 10
	}

	// Temperature severity
	if climat.Temperature < 0 || climat.Temperature > 40 {
		severity += 30
	}

	// Humidity severity
	if climat.Humidite < 10 || climat.Humidite > 90 {
		severity += 15
	}

	// CO2 severity
	if climat.Co2 > 70 {
		severity += 25
	}

	// O2 severity
	if climat.O2 < 15 {
		severity += 20
	}

	return severity
}

// UpdateEnergie: update the energy of the insect based on the severity of the environment
func (in *Insecte) UpdateEnergie(severity int) {
	attributes := enums.SpeciesAttributes[in.Espece]
	maxEnergie := attributes.NiveauEnergie
	energyLossRatio := float64(severity) / 100.0
	energyLoss := float64(maxEnergie) * energyLossRatio

	// update the energy of the insect, make sure it's between 0 and maxEnergie
	in.Energie = utils.Intmax(0, utils.Intmin(maxEnergie, int(math.Floor(float64(in.Energie)-energyLoss))))
}

// ============================================= End of UpdateEnergie_Incendie =======================================================
