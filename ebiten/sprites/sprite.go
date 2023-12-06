package sprites

import (
	"bytes"
	"image"
	"log"
	"math"
	"vivarium/ebiten/assets/images"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/rand"
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
}

func (s *Sprite) Update(deltaTime float64) {
	s.frameIndex++
	// Calculate the distance to move this frame
	distX := s.TargetX - s.X
	distY := s.TargetY - s.Y
	//fmt.Println("distX:", distX, "distY:", distY, s.Speed*deltaTime)
	// Move the sprite towards the target position
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
	if s.State == Moving {
		currentFrame = s.MoveFrames[(FrameIndex/framePerSwitch)%len(s.MoveFrames)]
	} else if s.State == Attacking {
		currentFrame = s.AttackFrames[(FrameIndex/framePerSwitch)%len(s.AttackFrames)]
	} else if s.State == Dying {
		currentFrame = s.DieFrames[(FrameIndex/framePerSwitch)%len(s.DieFrames)]
	} else if s.State == Idle {
		currentFrame = s.IdleFrames[(FrameIndex/framePerSwitch)%len(s.IdleFrames)]
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.X, s.Y)
	screen.DrawImage(currentFrame, op)
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

func NewSpiderSprite(X, Y float64, state SpriteState) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Spider_png))
	if err != nil {
		log.Fatal(err)
	}
	spiderImage := ebiten.NewImageFromImage(img)

	s := &Sprite{
		image: spiderImage,
		X:     X,
		Y:     Y,
		State: Idle,
		Id:    rand.Intn(100000),
		Speed: 10,

		frameIndex: 0,
	}

	s.IdleFrames = loadFrames(spiderImage, 5, 0)
	s.MoveFrames = loadFrames(spiderImage, 6, 1)
	s.AttackFrames = loadFrames(spiderImage, 9, 2)
	s.DieFrames = loadFrames(spiderImage, 9, 6)

	s.TargetX = X + 50
	s.TargetY = Y + 50
	return s
}

func NewSnailSprite(X, Y float64) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Snail_png))
	if err != nil {
		log.Fatal(err)
	}
	snailImage := ebiten.NewImageFromImage(img)

	return &Sprite{
		image: snailImage,
		X:     X,
		Y:     Y,
		State: Idle,
		Id:    rand.Intn(100000),
		Speed: 10,

		frameIndex: 0,

		IdleFrames:   loadFrames(snailImage, 1, 0),
		MoveFrames:   loadFrames(snailImage, 4, 1),
		AttackFrames: loadFrames(snailImage, 1, 2),
		DieFrames:    loadFrames(snailImage, 4, 3),
	}
}
