package sprites

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"image"
	"log"
	"vivarium/ebiten/assets/images"
	"vivarium/enums"
	"vivarium/organisme"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32

	framePerSwitch = 10 // It decides the speed of animation: the greater the slower
)

type SpriteState int

const (
	Idle SpriteState = iota
	Moving
	Attacking
	Eating
	Sexing
	Winning
	Losing
)

func sign(x float64) float64 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// 用于存储每个生物agent的状态
type Sprite struct {
	X float64
	Y float64

	TargetX float64
	TargetY float64
	Speed   float64

	Id int

	image *ebiten.Image

	State        SpriteState
	IdleFrames   []*ebiten.Image
	MoveFrames   []*ebiten.Image
	AttackFrames []*ebiten.Image
	DieFrames    []*ebiten.Image

	frameIndex int

	Species string

	IsDead            bool
	DyingCount        int
	EatingCount       int
	AttackingCount    int
	IsDying           bool
	StatusCountWinner int
	StatusCountLoser  int
	IsInsect          bool

	IsManger     bool
	IsReproduire bool
	IsSeDeplacer bool
	IsSeBattre   bool
	IsWinner     bool
	IsLooser     bool
	IsNormal     bool

	NbParts int
}

// After each update request, the sprite status will be updated according to the agent.
// If the id is not in the map, the sprite will be automatically generated.
func UpdateOrganisme(spriteMap map[int]*Sprite, org organisme.Organisme) {
	switch o := org.(type) {
	case *organisme.Insecte:
		UpdateInsecte(spriteMap, o)
	case *organisme.Plante:
		UpdatePlante(spriteMap, o)
	}
}

func UpdateInsecte(spriteMap map[int]*Sprite, org *organisme.Insecte) {
	spriteInfo := spriteMap[org.GetID()]
	// Generate a random number from -0。25 to 0.25 to prevent sprites from overlapping
	randBigInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return
	}
	randNum := (float64(randBigInt.Int64())/1000000.0 - 0.5) / 2
	spriteInfo.TargetX = 15 * (float64(org.GetPosX()) + randNum + 1)
	spriteInfo.TargetY = 15 * (float64(org.GetPosY()) + randNum + 1)

	// Recenter the sprite, correct the offset caused by the 32*32 pixel sprite
	if spriteInfo.Species == "AraignéeSauteuse" {
		spriteInfo.TargetX = 16*(float64(org.GetPosX())+randNum+1) - 8
		spriteInfo.TargetY = 16*(float64(org.GetPosY())+1) - 16

	} else if spriteInfo.Species == "PetitSerpent" {
		spriteInfo.TargetX = 16*(float64(org.GetPosX())+randNum+1) - 4
		spriteInfo.TargetY = 16*(float64(org.GetPosY())+1) - 16

	} else if spriteInfo.Species == "Escargot" {
		spriteInfo.TargetX = 16*(float64(org.GetPosX())+randNum+1) - 9
		spriteInfo.TargetY = 16*(float64(org.GetPosY())+1) - 16

	} else if spriteInfo.Species == "Grillons" {
		spriteInfo.TargetX = 16*(float64(org.GetPosX())+randNum+1) - 7
		spriteInfo.TargetY = 16*(float64(org.GetPosY())+1) - 16
	} else if spriteInfo.Species == "Lombric" {
		spriteInfo.TargetX = 16*(float64(org.GetPosX())+randNum+1) - 6
		spriteInfo.TargetY = 16*(float64(org.GetPosY())+1) - 16
	}

	spriteInfo.Species = org.GetEspece().String()
	spriteInfo.DyingCount = 0
	spriteInfo.IsDying = org.GetEtat()
	spriteInfo.IsInsect = true
	spriteInfo.StatusCountWinner = 0
	spriteInfo.StatusCountLoser = 0

	// States peculiar to insects
	spriteInfo.IsManger = org.IsManger
	spriteInfo.IsSeDeplacer = org.IsSeDeplacer
	spriteInfo.IsSeBattre = org.IsSeBattre
	spriteInfo.IsWinner = org.IsWinner
	spriteInfo.IsLooser = org.IsLooser
	spriteInfo.IsNormal = org.IsNormal

	// States peculiar to plants
	spriteInfo.NbParts = 1

	// States for both insects and plants
	spriteInfo.IsReproduire = org.IsReproduire

	spriteMap[org.GetID()] = spriteInfo
}

func UpdatePlante(spriteMap map[int]*Sprite, org *organisme.Plante) {
	spriteInfo := spriteMap[org.GetID()]
	spriteInfo.X = 16 * float64(org.GetPosX()+1)
	spriteInfo.Y = 16 * float64(org.GetPosY()+1)

	spriteInfo.TargetX = 16 * float64(org.GetPosX()+1)
	spriteInfo.TargetY = 16 * float64(org.GetPosY()+1)

	spriteInfo.Species = org.GetEspece().String()
	spriteInfo.DyingCount = 0
	spriteInfo.IsDying = org.GetEtat()
	spriteInfo.IsInsect = false
	spriteInfo.StatusCountWinner = 0
	spriteInfo.StatusCountLoser = 0

	// States peculiar to plants
	spriteInfo.NbParts = org.NbParts

	// States for both insects and plants
	spriteInfo.IsReproduire = org.IsReproduire

	spriteMap[org.GetID()] = spriteInfo

}

func (s *Sprite) Update(deltaTime float64) {
	// Update sprite frame index
	s.frameIndex++

	if s.IsDead {
		// If the sprite is dead, no rendering occurs
		return
	}

	if s.IsNormal == false {
		if s.IsInsect {
			if s.IsManger {
				// Execute logic related to eating
				s.State = Eating
			}
			if s.IsReproduire {
				// Execute logic related to reproduction
				s.State = Sexing
			}
			if s.IsSeBattre {
				if s.IsWinner {
					if s.StatusCountWinner <= 20 {
						s.StatusCountWinner++
						// Execute logic related to winner
						s.State = Winning
					}
					s.StatusCountWinner = 0
				} else if s.IsLooser {
					if s.StatusCountLoser <= 20 {
						s.StatusCountLoser++
						// Execute logic related to loser
						s.State = Losing
					}
					s.StatusCountLoser = 0
				} else {
					// Execute logic related to fight
					s.State = Attacking
					s.AttackingCount = 0
				}
			}
		} else {
			// Execute logic related to reproduction
			if s.IsReproduire {
				if s.EatingCount <= 20 {
					s.EatingCount++
					s.State = Sexing
				} else {
					s.State = Idle
					s.EatingCount = 0
				}
			} else {
				s.State = Idle
			}
		}
	} else {
		// Execute logic for normal state
		s.State = Idle
	}

	// Calculate the distance to move this frame
	distX := s.TargetX - s.X
	distY := s.TargetY - s.Y

	if (distX != 0) && (distY != 0) {
		// If the sprite is moving, update the sprite status to Moving
		s.State = Moving
	} else if s.State == Moving {
		s.State = Idle
	}
	// Move the sprite X and Y towards the target position
	if math.Abs(distX) > s.Speed*deltaTime {
		s.X += s.Speed * deltaTime * sign(distX)
	} else {
		s.X = s.TargetX
	}

	if math.Abs(distY) > s.Speed*deltaTime {
		s.Y += s.Speed * deltaTime * sign(distY)
	} else {
		s.Y = s.TargetY
	}

}

func (s *Sprite) MoveTo(x, y float64) {
	s.TargetX = x
	s.TargetY = y
}

func (s *Sprite) Draw(screen *ebiten.Image, FrameIndex int) {
	var currentFrame *ebiten.Image

	if s.IsDead {
		return
	}

	if s.Species == "PetitHerbe" {
		if s.State == Sexing {
			img, _, err := image.Decode(bytes.NewReader(images.Seed1_png))
			if err != nil {
				log.Fatal(err)
			}
			Img := ebiten.NewImageFromImage(img)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(s.X, s.Y-6)
			screen.DrawImage(Img, op2)
		}
		if s.IsDying {
			s.IsDead = true
			return
		} else {
			currentFrame = s.image
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(s.X, s.Y)
			screen.DrawImage(currentFrame, op)

			return
		}
	}

	if s.Species == "GrandHerbe" {
		if s.State == Sexing {
			img, _, err := image.Decode(bytes.NewReader(images.Seed1_png))
			if err != nil {
				log.Fatal(err)
			}
			Img := ebiten.NewImageFromImage(img)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(s.X, s.Y-8)
			screen.DrawImage(Img, op2)
		}

		if s.IsDying {
			s.IsDead = true
			return
		} else {
			switch s.NbParts {
			case 0:
				s.IsDead = true
				return
			case 1:
				currentFrame = s.IdleFrames[0]
			case 2:
				currentFrame = s.MoveFrames[0]
			case 3:
				currentFrame = s.AttackFrames[0]
			case 4:
				currentFrame = s.DieFrames[0]
			}

			//currentFrame = s.image
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(s.X, s.Y)
			screen.DrawImage(currentFrame, op)
			return
		}
	}

	if s.Species == "Champignon" {
		if s.IsDying {
			s.IsDead = true
			return
		}

		if s.State == Sexing {
			img, _, err := image.Decode(bytes.NewReader(images.Baozi2_png))
			if err != nil {
				log.Fatal(err)
			}
			Img := ebiten.NewImageFromImage(img)

			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(s.X, s.Y-8)
			screen.DrawImage(Img, op2)
		}
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)
		return
	}

	// Animation of insect
	if s.IsDead {
		// If the sprite is dead, no rendering occurs
		return
	}

	if s.IsDying {
		currentFrame = s.DieFrames[(FrameIndex/(framePerSwitch))%len(s.DieFrames)]
		s.DyingCount++
		if s.DyingCount >= len(s.DieFrames) {
			if s.DyingCount >= 20 {
				s.IsDead = true
				return
			}
			currentFrame = s.DieFrames[len(s.DieFrames)-1]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(s.X, s.Y)
			screen.DrawImage(currentFrame, op)
			return
		}
	} else if s.State == Moving {
		currentFrame = s.MoveFrames[(FrameIndex/framePerSwitch)%len(s.MoveFrames)]
	} else if s.State == Attacking {
		currentFrame = s.AttackFrames[(FrameIndex/framePerSwitch)%len(s.AttackFrames)]
		s.AttackingCount++
		if s.AttackingCount >= len(s.AttackFrames)*5 {
			s.AttackingCount = 0
			s.State = Idle
		}
	} else if s.State == Idle {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]
	} else if s.State == Eating {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		img, _, err := image.Decode(bytes.NewReader(images.Wing_png))
		if err != nil {
			log.Fatal(err)
		}
		Img := ebiten.NewImageFromImage(img)

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(s.X+8, s.Y+12)

		screen.DrawImage(Img, op2)
		return
	} else if s.State == Sexing {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		// Heart_png for sexing!!!
		img, _, err := image.Decode(bytes.NewReader(images.Heart_png))
		if err != nil {
			log.Fatal(err)
		}
		heartImg := ebiten.NewImageFromImage(img)

		heartFrame := make([]*ebiten.Image, 5)
		for i := range heartFrame {
			sX, sY := frameOX+i*16, frameOY
			frame := heartImg.SubImage(image.Rect(sX, sY, sX+16, sY+16)).(*ebiten.Image)
			heartFrame[i] = frame
		}

		currentFrame2 := heartFrame[(FrameIndex/framePerSwitch)%len(heartFrame)]

		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(s.X+8, s.Y+12)
		screen.DrawImage(currentFrame2, op2)
		return
	} else if s.State == Winning {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		// Crown_png for winning!!!
		img, _, err := image.Decode(bytes.NewReader(images.Crown_png))
		if err != nil {
			log.Fatal(err)
		}
		Img := ebiten.NewImageFromImage(img)
		crownFrame := LoadFrames(Img, 4, 2)

		currentFrame2 := crownFrame[(FrameIndex/framePerSwitch)%len(crownFrame)]

		op2 := &ebiten.DrawImageOptions{}
		scaleX := 0.5
		scaleY := 0.5
		op2.GeoM.Scale(scaleX, scaleY)
		op2.GeoM.Translate(s.X+8, s.Y+16)
		screen.DrawImage(currentFrame2, op2)
		return
	}
	if currentFrame != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)
	}

}

func LoadFrames(img *ebiten.Image, frameCount, stateIdx int) []*ebiten.Image {
	frames := make([]*ebiten.Image, frameCount)
	for i := range frames {
		sX, sY := frameOX+i*frameWidth, frameOY+frameHeight*stateIdx
		frame := img.SubImage(image.Rect(sX, sY, sX+frameWidth, sY+frameHeight)).(*ebiten.Image)
		frames[i] = frame
	}
	return frames
}

func loadFramesWidthHeight(img *ebiten.Image, frameCount, stateIdx int, _frameWidth, _framHeight int) []*ebiten.Image {
	frames := make([]*ebiten.Image, frameCount)
	for i := range frames {
		sX, sY := frameOX+i*_frameWidth, frameOY+_framHeight*stateIdx
		frame := img.SubImage(image.Rect(sX, sY, sX+_frameWidth, sY+_framHeight)).(*ebiten.Image)
		frames[i] = frame
	}
	return frames
}

func NewBaseSprite(org organisme.Organisme) *Sprite {
	sprite := &Sprite{
		X:     15 * float64(org.GetPosX()+1),
		Y:     15 * float64(org.GetPosY()+1),
		Speed: 10,
		Id:    org.GetID(),

		TargetX: 15 * float64(org.GetPosX()+1),
		TargetY: 15 * float64(org.GetPosY()+1),

		//frameIndex int
		Species:           org.GetEspece().String(),
		DyingCount:        0,
		IsDying:           org.GetEtat(),
		IsInsect:          true,
		StatusCountWinner: 0,
		StatusCountLoser:  0,

		// States peculiar to insects
		IsManger:     false,
		IsReproduire: false,
		IsSeDeplacer: false,
		IsSeBattre:   false,
		IsWinner:     false,
		IsLooser:     false,
		IsNormal:     false,

		// States peculiar to plants
		NbParts: 1,
	}
	return sprite
}

func NewSpiderSprite(org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Spider_png))

	switch o := org.(type) {
	case *organisme.Insecte:
		if o.Sexe == enums.Male {
			img, _, err = image.Decode(bytes.NewReader(images.Spider_png))
		} else if o.Sexe == enums.Femelle {
			img, _, err = image.Decode(bytes.NewReader(images.Spider_female_png))
		}
		if err != nil {
			log.Fatal(err)
		}
	case *organisme.Plante:
		fmt.Errorf("error: newspider found plants")
		return nil
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 80

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = LoadFrames(sprite.image, 5, 0)
	sprite.MoveFrames = LoadFrames(sprite.image, 6, 1)
	sprite.AttackFrames = LoadFrames(sprite.image, 9, 2)
	sprite.DieFrames = LoadFrames(sprite.image, 9, 6)

	return sprite
}

func NewSnailSprite(org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Snail_png))

	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 16

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = LoadFrames(sprite.image, 4, 0)
	sprite.MoveFrames = LoadFrames(sprite.image, 6, 1)
	sprite.AttackFrames = LoadFrames(sprite.image, 6, 2)
	sprite.DieFrames = LoadFrames(sprite.image, 7, 4)

	return sprite
}

func NewCobraSprite(org organisme.Organisme) *Sprite {

	img, _, err := image.Decode(bytes.NewReader(images.Cobra_png))

	switch o := org.(type) {
	case *organisme.Insecte:
		if o.Sexe == enums.Male {
			img, _, err = image.Decode(bytes.NewReader(images.Cobra_male_png))
		} else if o.Sexe == enums.Femelle {
			img, _, err = image.Decode(bytes.NewReader(images.Cobra_png))
		}
		if err != nil {
			log.Fatal(err)
		}
	case *organisme.Plante:
		fmt.Errorf("error: newspider found plants")
		return nil
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 64

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = LoadFrames(sprite.image, 8, 0)
	sprite.MoveFrames = LoadFrames(sprite.image, 8, 1)
	sprite.AttackFrames = LoadFrames(sprite.image, 6, 2)
	sprite.DieFrames = LoadFrames(sprite.image, 6, 4)

	return sprite
}

func NewWormSprite(org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Worm_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 32

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = LoadFrames(sprite.image, 9, 0)
	sprite.MoveFrames = LoadFrames(sprite.image, 6, 1)
	sprite.AttackFrames = LoadFrames(sprite.image, 6, 2)
	sprite.DieFrames = LoadFrames(sprite.image, 6, 3)

	return sprite
}

func NewScarabSprite(org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Scarab_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 48

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = LoadFrames(sprite.image, 4, 0)
	sprite.MoveFrames = LoadFrames(sprite.image, 4, 1)
	sprite.AttackFrames = LoadFrames(sprite.image, 4, 2)
	sprite.DieFrames = LoadFrames(sprite.image, 5, 4)

	return sprite
}

func NewMushroomSprite(org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Mushroom_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = loadFramesWidthHeight(sprite.image, 5, 9, 16, 16)
	sprite.MoveFrames = loadFramesWidthHeight(sprite.image, 5, 9, 16, 16)
	sprite.AttackFrames = loadFramesWidthHeight(sprite.image, 5, 9, 16, 16)
	sprite.DieFrames = loadFramesWidthHeight(sprite.image, 5, 9, 16, 16)

	return sprite
}

func NewPetitHerbeSprite(org organisme.Organisme) *Sprite {
	randGrassInt, err := rand.Int(rand.Reader, big.NewInt(5))

	if err != nil {
		fmt.Errorf("error: %v", err)
		log.Fatal(err)
	}

	var img image.Image
	switch randGrassInt.Int64() {
	case 0:
		img, _, err = image.Decode(bytes.NewReader(images.SGrass6_png))
	case 1:
		img, _, err = image.Decode(bytes.NewReader(images.SGrass7_png))
	case 2:
		img, _, err = image.Decode(bytes.NewReader(images.SGrass8_png))
	case 3:
		img, _, err = image.Decode(bytes.NewReader(images.SGrass9_png))
	case 4:
		img, _, err = image.Decode(bytes.NewReader(images.SGrass5_png))
	}

	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = loadFramesWidthHeight(sprite.image, 1, 0, 16, 16)
	sprite.MoveFrames = loadFramesWidthHeight(sprite.image, 1, 0, 16, 16)
	sprite.AttackFrames = loadFramesWidthHeight(sprite.image, 1, 0, 16, 16)
	sprite.DieFrames = loadFramesWidthHeight(sprite.image, 1, 0, 16, 16)

	return sprite
}

func NewGrandHerbeSprite(org organisme.Organisme) *Sprite {
	img1, _, err := image.Decode(bytes.NewReader(images.Grass1_png))
	img2, _, err := image.Decode(bytes.NewReader(images.Grass2_png))
	img3, _, err := image.Decode(bytes.NewReader(images.Grass3_png))
	img4, _, err := image.Decode(bytes.NewReader(images.Grass4_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.image = ebiten.NewImageFromImage(img1)
	sprite.State = Idle
	sprite.IdleFrames = loadFramesWidthHeight(ebiten.NewImageFromImage(img1), 1, 0, 16, 16)
	sprite.MoveFrames = loadFramesWidthHeight(ebiten.NewImageFromImage(img2), 1, 0, 16, 16)
	sprite.AttackFrames = loadFramesWidthHeight(ebiten.NewImageFromImage(img3), 1, 0, 16, 16)
	sprite.DieFrames = loadFramesWidthHeight(ebiten.NewImageFromImage(img4), 1, 0, 16, 16)

	return sprite
}
