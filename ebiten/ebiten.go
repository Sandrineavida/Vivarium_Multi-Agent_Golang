package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"sort"
	"time"
	server "vivarium"

	"vivarium/ebiten/assets/images"
	sprite "vivarium/ebiten/sprites"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 272
	screenHeight = 272

	tileSize = 16

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8

	menuBarWidth = 50 // 菜单栏宽度

	// 假设按钮的尺寸为30x15（宽度x高度）
	buttonWidth  = 30
	buttonHeight = 15
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
	sprites        []*sprite.Sprite
	lastUpdateTime time.Time
	layers         [][]int
	grassLayer     []int
	updateInterval int
	updateCount    int
	SpriteMap      map[int]*sprite.Sprite // 新增精灵映射

	// 菜单栏
	menuBarImage *ebiten.Image

	// 按钮
	isPaused          bool
	pauseButtonImage  *ebiten.Image
	resumeButtonImage *ebiten.Image
	pauseButtonRect   image.Rectangle
	resumeButtonRect  image.Rectangle
}

func (g *Game) initButtons() {
	// 暂停和继续按钮的Y轴位置
	const pauseButtonY = 50
	const resumeButtonY = 80

	// X轴位置是屏幕宽度加上菜单栏宽度的一半减去按钮宽度的一半
	buttonX := screenWidth + (menuBarWidth-buttonWidth)/2

	// 初始化暂停和继续按钮的位置
	g.pauseButtonRect = image.Rect(buttonX, pauseButtonY, buttonX+buttonWidth, pauseButtonY+buttonHeight)
	g.resumeButtonRect = image.Rect(buttonX, resumeButtonY, buttonX+buttonWidth, resumeButtonY+buttonHeight)

	// 创建纯色按钮图像
	g.pauseButtonImage = ebiten.NewImage(buttonWidth, buttonHeight)
	g.pauseButtonImage.Fill(color.RGBA{R: 255, A: 255}) // 红色暂停按钮
	g.resumeButtonImage = ebiten.NewImage(buttonWidth, buttonHeight)
	g.resumeButtonImage.Fill(color.RGBA{G: 255, A: 255}) // 绿色继续按钮
}

func (g *Game) Update() error {

	// 检测鼠标点击并更新按钮状态
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if isButtonClicked(x, y, g.pauseButtonRect) {
			g.isPaused = true
			// 发送暂停命令到后端
		} else if isButtonClicked(x, y, g.resumeButtonRect) {
			g.isPaused = false
			// 发送继续命令到后端
		}
	}

	if g.isPaused {
		// 如果暂停，跳过仿真逻辑
		return nil
	}

	g.FrameIndex++

	//for smooth moving
	currentTime := time.Now()
	deltaTime := currentTime.Sub(g.lastUpdateTime).Seconds() * 4
	g.lastUpdateTime = currentTime

	// 每秒60帧，所以每30帧是0.5秒
	g.updateInterval = 30

	// 增加更新计数
	g.updateCount++

	// 每隔一段时间更新一次
	if g.updateCount >= g.updateInterval {

		if server.EcosystemForEbiten == nil {
			return nil // 或其他适当的错误处理
		}

		// 从服务器读取数据
		ecosystemForEbiten := server.EcosystemForEbiten

		for _, org := range ecosystemForEbiten.GetAllOrganisms() {
			if _, exists := g.SpriteMap[org.GetID()]; !exists {
				// 如果SpriteMap中没有这个ID，创建一个新的蜘蛛精灵
				// 后期根据org.getEspace()来确定使用那个sprite.New
				if org.GetEspece().String() == "AraignéeSauteuse" {
					g.SpriteMap[org.GetID()] = sprite.NewSpiderSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "PetitSerpent" {
					g.SpriteMap[org.GetID()] = sprite.NewCobraSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Grillons" {
					g.SpriteMap[org.GetID()] = sprite.NewScarabSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Escargot" {
					g.SpriteMap[org.GetID()] = sprite.NewSnailSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Lombric" {
					g.SpriteMap[org.GetID()] = sprite.NewWormSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Champignon" {
					g.SpriteMap[org.GetID()] = sprite.NewMushroomSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "PetitHerbe" {
					g.SpriteMap[org.GetID()] = sprite.NewPetitHerbeSprite(g.SpriteMap, org)
				} else if org.GetEspece().String() == "GrandHerbe" {
					g.SpriteMap[org.GetID()] = sprite.NewGrandHerbeSprite(g.SpriteMap, org)
				} else {
					fmt.Println("Error: Unknown Espece", org.GetEspece().String())
				}

			} else {
				if g.SpriteMap[org.GetID()].IsDead {
					//delete(g.SpriteMap, org.GetID())
					continue
				}
				// 更新生物的 Sprite 信息
				if org.GetEspece().String() == "AraignéeSauteuse" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "PetitSerpent" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Grillons" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Escargot" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Lombric" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "Champignon" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "PetitHerbe" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else if org.GetEspece().String() == "GrandHerbe" {
					sprite.UpdateOrganisme(g.SpriteMap, org)
				} else {
					fmt.Println("Error: Unknown Espece", org.GetEspece().String())
				}
			}
		}

		// 重置计数器
		g.updateCount = 0

		// if ecosystemForEbiten.Climat.Meteo == enums.Incendie {
		// 	//fmt.Println("Incendie")

		// }
	}

	for _, sprite := range g.SpriteMap {
		if sprite.IsDead {
			continue
		}
		sprite.Update(deltaTime)
	}

	// for _, sprite := range g.sprites {
	// 	sprite.Update(deltaTime)
	// }

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

func (g *Game) DrawGrass(screen *ebiten.Image) {
	w := tilesImage.Bounds().Dx()
	tileXCount := w / tileSize

	const xCount = screenWidth / tileSize
	for i, t := range g.grassLayer {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))

		sx := (t % tileXCount) * tileSize
		sy := (t / tileXCount) * tileSize
		screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	}
}

func (g *Game) drawMenuBar(screen *ebiten.Image) {
	menuBarOp := &ebiten.DrawImageOptions{}
	menuBarOp.GeoM.Translate(float64(screenWidth), 0)
	screen.DrawImage(g.menuBarImage, menuBarOp)

	pauseButtonOp := &ebiten.DrawImageOptions{}
	pauseButtonOp.GeoM.Translate(float64(g.pauseButtonRect.Min.X), float64(g.pauseButtonRect.Min.Y))
	screen.DrawImage(g.pauseButtonImage, pauseButtonOp)

	resumeButtonOp := &ebiten.DrawImageOptions{}
	resumeButtonOp.GeoM.Translate(float64(g.resumeButtonRect.Min.X), float64(g.resumeButtonRect.Min.Y))
	screen.DrawImage(g.resumeButtonImage, resumeButtonOp)
}

func isButtonClicked(x, y int, buttonRect image.Rectangle) bool {
	return x >= buttonRect.Min.X && x <= buttonRect.Max.X &&
		y >= buttonRect.Min.Y && y <= buttonRect.Max.Y
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.DrawBackground(screen)

	// 定义优先级映射
	priorityMap := map[string]int{
		"PetitHerbe":       1,
		"GrandHerbe":       2,
		"Champignon":       3,
		"Escargot":         4,
		"Grillons":         5,
		"Lombric":          6,
		"PetitSerpent":     7,
		"AraignéeSauteuse": 8,
		// PetitHerbe
		// GrandHerbe
		// Champignon
		// Escargot
		// Grillons
		// Lombric
		// PetitSerpent
		// AraignéeSauteuse
	}

	// 创建一个按优先级排序的精灵切片
	sortedSprites := make([]*sprite.Sprite, 0, len(g.SpriteMap))
	for _, sprite := range g.SpriteMap {
		sortedSprites = append(sortedSprites, sprite)
	}

	// 使用自定义排序
	sort.Slice(sortedSprites, func(i, j int) bool {
		return priorityMap[sortedSprites[i].Species] < priorityMap[sortedSprites[j].Species]
	})

	// 遍历排序后的精灵并绘制它们
	for _, sprite := range sortedSprites {
		sprite.Draw(screen, g.FrameIndex)
	}

	// 遍历所有精灵并绘制它们
	// for _, sprite := range g.SpriteMap {
	// 	sprite.Draw(screen, g.FrameIndex)
	// }

	// 绘制菜单栏和按钮
	g.drawMenuBar(screen)

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
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,

				247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247, 247,
			},
			{
				188, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 189, 190,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,

				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,

				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 215,
				238, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 239, 240,
			},
		},
		lastUpdateTime: time.Now(),
		SpriteMap:      make(map[int]*sprite.Sprite),

		grassLayer: []int{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},

		menuBarImage: ebiten.NewImage(menuBarWidth, screenHeight),
	}

	g.initButtons()

	ebiten.SetWindowSize(screenWidth*2+menuBarWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Multi agent system")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
