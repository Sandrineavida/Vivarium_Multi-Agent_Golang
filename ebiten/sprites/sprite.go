package sprites

import (
	"bytes"
	"fmt"
	"image"
	"log"
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

type Sprite struct {
	X float64
	Y float64

	Id int

	image *ebiten.Image

	State        SpriteState
	IdleFrames   []*ebiten.Image
	MoveFrames   []*ebiten.Image
	AttackFrames []*ebiten.Image
	DieFrames    []*ebiten.Image

	frameIndex int
}

func (s *Sprite) Update() {
	s.frameIndex++
	fmt.Println(s.frameIndex)
}

func (s *Sprite) MoveTo(x, y float64) {
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
		State: state,
		Id:    rand.Intn(100000),

		frameIndex: 0,
	}

	s.IdleFrames = loadFrames(spiderImage, 5, 0)
	s.MoveFrames = loadFrames(spiderImage, 6, 1)
	s.AttackFrames = loadFrames(spiderImage, 9, 2)
	s.DieFrames = loadFrames(spiderImage, 9, 6)
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
	}
}
