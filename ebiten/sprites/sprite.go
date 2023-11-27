package sprites

import (
	"bytes"
	"image"
	"log"
	"vivarium/ebiten/assets/images"

	"github.com/hajimehoshi/ebiten/v2"
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

	image *ebiten.Image

	state  SpriteState
	Idle   []*ebiten.Image
	Move   []*ebiten.Image
	Attack []*ebiten.Image
	Die    []*ebiten.Image
}

func (s *Sprite) Update() {
	// 更新Sprite的状态，例如移动位置等
}

func (s *Sprite) MoveTo(x, y float64) {

}

func NewSpiderSprite(X, Y float64) *Sprite {
	img, _, err := image.Decode(bytes.NewReader(images.Spider_png))
	if err != nil {
		log.Fatal(err)
	}
	spiderImage := ebiten.NewImageFromImage(img)

	return &Sprite{
		image: spiderImage,
		X:     X,
		Y:     Y,
	}
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
