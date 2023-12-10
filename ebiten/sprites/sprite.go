package sprites

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/hajimehoshi/ebiten/v2"

	//"golang.org/x/exp/rand"
	"image"
	"log"
	"vivarium/ebiten/assets/images"
	"vivarium/organisme"
)

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8

	framePerSwitch = 10 // It decides the speed of animation: the greater the slower
)

type SpriteState int

const (
	Idle SpriteState = iota
	Moving
	Attacking
	Dying
	Eating
	Sexing
	Winning
	Losing
)

type SpriteType int

const (
	Spider SpriteType = iota
	Snail
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

//var SpriteMap = make(map[int]*Sprite)

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

	//-----------------------------------------------------
	Species string

	IsDead            bool
	DyingCount        int
	EatingCount       int
	IsDying           bool
	StatusCountWinner int
	StatusCountLoser  int

	IsInsect bool

	// 昆虫特有的状态
	IsManger     bool
	IsReproduire bool
	IsSeDeplacer bool
	IsSeBattre   bool
	IsWinner     bool
	IsLooser     bool
	IsNormal     bool

	// 植物特有的状态
	NbParts int
}

// 每次update请求后，都会根据agent更新精灵状态，如果该id不在map中，自动生成精灵
func UpdateOrganisme(spriteMap map[int]*Sprite, org organisme.Organisme) {
	switch o := org.(type) {
	case *organisme.Insecte:
		UpdateInsecte(spriteMap, o) // o 是 *organisme.Insecte 类型
	case *organisme.Plante:
		UpdatePlante(spriteMap, o)
	}
}

func UpdateInsecte(spriteMap map[int]*Sprite, org *organisme.Insecte) {
	spriteInfo := spriteMap[org.GetID()]
	// 15 ？
	// hotfix-1210: 新增一个0到0.5的随机数，防止精灵重叠
	// 生成一个介于 0 和 1 之间的加密安全的随机数
	randBigInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		// 处理错误
	}
	// 将随机数转换为浮点数，并缩小到 0 到 0.5 的范围
	randNum := float64(randBigInt.Int64()) / 2000000.0
	spriteInfo.TargetX = 15 * (float64(org.GetPosX()) + randNum)
	spriteInfo.TargetY = 15 * (float64(org.GetPosY()) + randNum)

	spriteInfo.Species = org.GetEspece().String()
	spriteInfo.DyingCount = 0
	spriteInfo.IsDying = org.GetEtat()
	spriteInfo.IsInsect = true
	spriteInfo.StatusCountWinner = 0
	spriteInfo.StatusCountLoser = 0

	// 昆虫特有的状态
	spriteInfo.IsManger = org.IsManger
	spriteInfo.IsReproduire = org.IsReproduire
	spriteInfo.IsSeDeplacer = org.IsSeDeplacer
	spriteInfo.IsSeBattre = org.IsSeBattre
	spriteInfo.IsWinner = org.IsWinner
	spriteInfo.IsLooser = org.IsLooser
	spriteInfo.IsNormal = org.IsNormal

	// 植物特有的状态
	spriteInfo.NbParts = 1

	spriteMap[org.GetID()] = spriteInfo
}

func UpdatePlante(spriteMap map[int]*Sprite, org *organisme.Plante) {
	spriteInfo := spriteMap[org.GetID()]
	spriteInfo.X = 15 * float64(org.GetPosX())
	spriteInfo.Y = 15 * float64(org.GetPosY())

	spriteInfo.Species = org.GetEspece().String()
	spriteInfo.DyingCount = 0
	spriteInfo.IsDying = org.GetEtat()
	spriteInfo.IsInsect = false
	spriteInfo.StatusCountWinner = 0
	spriteInfo.StatusCountLoser = 0

	// 植物特有的状态
	spriteInfo.NbParts = org.NbParts

	spriteMap[org.GetID()] = spriteInfo

}

func (s *Sprite) Update(deltaTime float64) {
	// 更新精灵帧索引
	s.frameIndex++

	if s.IsDead {
		// 如果精灵已死，不进行渲染
		return
	}

	if s.IsNormal == false {
		// 如果是昆虫
		if s.IsInsect {
			if s.IsManger {
				// 执行与进食相关的逻辑 戴个恰饭图标
				s.State = Eating
				fmt.Println("please eat aaaaaaaaaaaaaaaaaaaaaaaaaa")
			} else {
				fmt.Println("please dont eat aaaaaaaaaaaaaaaaaaaaaaaaaa")
			}

			if s.IsReproduire {
				// 执行与繁殖相关的逻辑 戴个💗💗💗
				s.State = Sexing
				fmt.Println("please fuck each other aaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
			if s.IsSeBattre {
				if s.IsWinner {
					if s.StatusCountWinner <= 20 {
						s.StatusCountWinner++
						// 执行胜利者的逻辑 戴个小王冠
						s.State = Winning
					}
					s.StatusCountWinner = 0
					fmt.Println("winwinwinwinwinwinwinwinwinwinwinwinwinwinwinwinwinwinwin")
				} else if s.IsLooser {
					if s.StatusCountLoser <= 20 {
						s.StatusCountLoser++
						// 执行失败者的逻辑 显示Loser
						s.State = Losing
					}
					s.StatusCountLoser = 0
					fmt.Println("losing losinglosinglosinglosinglosinglosinglosinglosinglosinglosinglosinglosinglosing")
				} else {
					// 执行正常战斗的逻辑 戴个打架图标
					s.State = Attacking
					fmt.Println("fighting fighting fighting fighting fighting fighting fighting fighting fighting fighting fighting ")
				}
			}
		} else {
			// 如果是植物
			if s.NbParts > 0 {
				// 根据NbParts=1-4显示百分比图标
			}
		}
	} else {
		// 执行正常状态的逻辑 无图标状态
	}

	// Calculate the distance to move this frame
	distX := s.TargetX - s.X
	distY := s.TargetY - s.Y

	if (distX != 0) && (distY != 0) {
		// 如果精灵正在移动，更新精灵状态
		s.State = Moving
	} else if s.State == Moving {
		s.State = Idle
	}
	//fmt.Println("distX:", distX, "distY:", distY, s.Speed*deltaTime)
	// Move the sprite X and Y towards the target position
	if math.Abs(distX) > s.Speed*deltaTime {
		s.X += s.Speed * deltaTime * sign(distX)
		//fmt.Println(s.Speed * deltaTime * sign(distX))
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
		// 如果精灵已死，不进行渲染
		return
	}

	if s.IsDying {
		currentFrame = s.DieFrames[(FrameIndex/framePerSwitch)%len(s.DieFrames)]
		s.DyingCount++
		if s.DyingCount >= len(s.DieFrames) {
			s.IsDead = true
			return
		}
	} else if s.State == Moving {
		currentFrame = s.MoveFrames[(FrameIndex/framePerSwitch)%len(s.MoveFrames)]
	} else if s.State == Attacking {
		currentFrame = s.AttackFrames[(FrameIndex/framePerSwitch)%len(s.AttackFrames)]
	} else if s.State == Idle {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]
	}

	if s.State == Eating {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		img, _, err := image.Decode(bytes.NewReader(images.Wing_png))
		if err != nil {
			log.Fatal(err)
		}
		Img := ebiten.NewImageFromImage(img)
		//wingFrame := loadFrames(Img, 1, 0)

		//currentFrame2 := wingFrame[(FrameIndex/framePerSwitch)%len(wingFrame)]

		op2 := &ebiten.DrawImageOptions{}

		// scaleX := 0.5
		// scaleY := 0.5
		// op2.GeoM.Scale(scaleX, scaleY)
		op2.GeoM.Translate(s.X+8, s.Y+12)

		screen.DrawImage(Img, op2)
		return
	}
	if s.State == Sexing {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		//heart for sexing!!!
		img, _, err := image.Decode(bytes.NewReader(images.Heart_png))
		if err != nil {
			log.Fatal(err)
		}
		heartImg := ebiten.NewImageFromImage(img)
		//heartFrame := loadFrames(heartImg, 5, 0)

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
	}
	if s.State == Winning {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.X, s.Y)
		screen.DrawImage(currentFrame, op)

		//heart for eating!!!
		img, _, err := image.Decode(bytes.NewReader(images.Crown_png))
		if err != nil {
			log.Fatal(err)
		}
		Img := ebiten.NewImageFromImage(img)
		crownFrame := loadFrames(Img, 4, 2)

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

func loadFrames(img *ebiten.Image, frameCount, stateIdx int) []*ebiten.Image {
	frames := make([]*ebiten.Image, frameCount)
	for i := range frames {
		sX, sY := frameOX+i*frameWidth, frameOY+frameHeight*stateIdx
		frame := img.SubImage(image.Rect(sX, sY, sX+frameWidth, sY+frameHeight)).(*ebiten.Image)
		frames[i] = frame
	}
	return frames
}

func NewBaseSprite(org organisme.Organisme) *Sprite {
	sprite := &Sprite{
		X:     15 * float64(org.GetPosX()),
		Y:     15 * float64(org.GetPosY()),
		Speed: 10,
		Id:    org.GetID(),

		//frameIndex int
		Species:           org.GetEspece().String(),
		DyingCount:        0,
		IsDying:           org.GetEtat(),
		IsInsect:          true,
		StatusCountWinner: 0,
		StatusCountLoser:  0,

		// 昆虫特有的状态
		IsManger:     false,
		IsReproduire: false,
		IsSeDeplacer: false,
		IsSeBattre:   false,
		IsWinner:     false,
		IsLooser:     false,
		IsNormal:     false,

		// 植物特有的状态
		NbParts: 1,
	}
	return sprite
}

func NewSpiderSprite(spriteMap map[int]*Sprite, org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Spider_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.Speed = 80

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = loadFrames(sprite.image, 5, 0)
	sprite.MoveFrames = loadFrames(sprite.image, 6, 1)
	sprite.AttackFrames = loadFrames(sprite.image, 9, 2)
	sprite.DieFrames = loadFrames(sprite.image, 9, 6)

	return sprite
}

func NewSpiderSprite2(X, Y float64, state SpriteState) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Spider_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := &Sprite{
		X:     X,
		Y:     Y,
		Speed: 10,

		TargetX: X + 20,
		TargetY: Y - 20,
	}

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = state
	sprite.IdleFrames = loadFrames(sprite.image, 5, 0)
	sprite.MoveFrames = loadFrames(sprite.image, 6, 1)
	sprite.AttackFrames = loadFrames(sprite.image, 9, 2)
	sprite.DieFrames = loadFrames(sprite.image, 9, 6)

	return sprite
}

func NewSnailSprite(spriteMap map[int]*Sprite, org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Snail_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = loadFrames(sprite.image, 1, 0)
	sprite.MoveFrames = loadFrames(sprite.image, 4, 1)
	sprite.AttackFrames = loadFrames(sprite.image, 1, 2)
	sprite.DieFrames = loadFrames(sprite.image, 4, 3)

	return sprite
}

func NewCobraSprite(spriteMap map[int]*Sprite, org organisme.Organisme) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Cobra_png))
	if err != nil {
		log.Fatal(err)
	}

	sprite := NewBaseSprite(org)

	sprite.image = ebiten.NewImageFromImage(img)
	sprite.State = Idle
	sprite.IdleFrames = loadFrames(sprite.image, 8, 0)
	sprite.MoveFrames = loadFrames(sprite.image, 8, 1)
	sprite.AttackFrames = loadFrames(sprite.image, 6, 2)
	sprite.DieFrames = loadFrames(sprite.image, 6, 4)

	return sprite
}
