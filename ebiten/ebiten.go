// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	server "vivarium"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"vivarium/ebiten/assets/images"
	sprite "vivarium/ebiten/sprites"
)

const (
	screenWidth  = 240
	screenHeight = 240
)
const (
	tileSize = 16
)
const (
	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var (
	runnerImage *ebiten.Image
	tilesImage  *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

type Game struct {
	FrameIndex     int
	sprites        []sprite.Sprite
	layers         [][]int
	updateInterval int // 新增：更新间隔计数
	updateCount    int // 新增：当前更新计数
}

func (g *Game) Update() error {
	g.FrameIndex++

	// 每秒60帧，所以每30帧是0.5秒
	g.updateInterval = 30

	// 增加更新计数
	g.updateCount++

	if g.updateCount >= g.updateInterval {

		if server.EcosystemForEbiten == nil {
			return nil // 或其他适当的错误处理
		}

		// 从服务器读取数据
		ecosystemForEbiten := server.EcosystemForEbiten

		// 更新所有生物的 Sprite 信息
		for _, org := range ecosystemForEbiten.GetAllOrganisms() {
			sprite.UpdateOrganisme(org)
		}

		// 更新 spriteMap 中所有 Sprite 的状态，用于渲染
		for _, sprite := range sprite.SpriteMap {
			sprite.Update()
		}

		// 重置计数器
		g.updateCount = 0
	}

	return nil
}

func (g *Game) DrawBackground(screen *ebiten.Image) {
	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	const xCount = screenWidth / tileSize
	for _, l := range g.layers {
		for i, t := range l {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))

			sx := (t % tileXCount) * tileSize
			sy := (t / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) DrawSprite(screen *ebiten.Image, sprite *sprite.Sprite) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	i := (g.FrameIndex / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy+frameHeight, sx+frameWidth, sy+frameHeight*2)).(*ebiten.Image), op)
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.DrawBackground(screen)

	for _, s := range g.sprites {
		s.Draw(screen, g.FrameIndex)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	// 初始化服务器
	go server.StartServer()

	g := &Game{
		layers: [][]int{
			{
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
			},
			{
				188, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 190,
				213, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 215,

				213, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,

				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 215,
				238, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 240,
			},
		},
		sprites: []sprite.Sprite{
			*sprite.NewSpiderSprite(screenWidth/2+20, screenHeight/2, sprite.Idle),
			*sprite.NewSpiderSprite(screenWidth/2-20, screenHeight/2, sprite.Moving),
			*sprite.NewSpiderSprite(screenWidth/2+20, screenHeight/2-20, sprite.Attacking),
			*sprite.NewSpiderSprite(screenWidth/2-20, screenHeight/2-20, sprite.Dying),
		},
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Multi agent system")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
