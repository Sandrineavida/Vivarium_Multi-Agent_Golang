package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"sort"
	"time"
	server "vivarium"
	"vivarium/climat"

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

	menuBarWidth = 100 // 菜单栏宽度

	isSimulationPaused = false
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

	// 按钮
	isPaused    bool
	buttonImage *ebiten.Image
	buttonRect  image.Rectangle

	buttonUnpressedImage *ebiten.Image
	buttonPressedImage   *ebiten.Image

	// 添加用于显示的字段
	CurrentHour   int
	CurrentClimat *climat.Climat
}

func (g *Game) initMenuBarAndButton() {
	// Initialize buttons
	unpressedImg, _, err := image.Decode(bytes.NewReader(images.ButtonUnpressed_png))
	if err != nil {
		log.Fatal(err)
	}
	g.buttonUnpressedImage = ebiten.NewImageFromImage(unpressedImg)

	pressedImg, _, err := image.Decode(bytes.NewReader(images.ButtonPressed_png))
	if err != nil {
		log.Fatal(err)
	}
	g.buttonPressedImage = ebiten.NewImageFromImage(pressedImg)

	// Set button rectangle (x, y, x+width, y+height)
	//g.buttonRect = image.Rect(screenWidth, 10, screenWidth+40, 40)
	g.buttonRect = image.Rect(160+screenWidth/2, 120, 200+screenWidth/2, 150)
}

func (g *Game) Update() error {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Println(g.buttonRect.Min.X, g.buttonRect.Max.X, g.buttonRect.Min.Y, g.buttonRect.Max.Y)
		if x >= g.buttonRect.Min.X && x <= g.buttonRect.Max.X && y >= g.buttonRect.Min.Y && y <= g.buttonRect.Max.Y {
			g.isPaused = !g.isPaused
			fmt.Println("我按按钮了我按按钮了我按按钮了我按按钮了我按按钮了我按按钮了我按按钮了我按按钮了我按按钮了")
			server.PauseSignal <- !isSimulationPaused
		}
	}

	if g.isPaused {
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

		// 从 server.EcosystemForEbiten 更新数据
		ecosystemData := server.EcosystemForEbiten
		if ecosystemData != nil {
			g.CurrentHour = ecosystemData.Hour
			g.CurrentClimat = ecosystemData.Climat
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

	// Draw button
	buttonOp := &ebiten.DrawImageOptions{}
	buttonOp.GeoM.Translate(float64(g.buttonRect.Min.X), float64(g.buttonRect.Min.Y))
	if g.isPaused {
		screen.DrawImage(g.buttonPressedImage, buttonOp)
	} else {
		screen.DrawImage(g.buttonUnpressedImage, buttonOp)
	}

	// 显示 Hour 和 Climat 数据
	if g.CurrentClimat != nil {
		climatInfo := fmt.Sprintf("Hour: %d\nMeteo: %s\nTemperature: %d°C\nHumidity: %.2f%%\nCO2: %.2f%%\nO2: %.2f%%",
			g.CurrentHour, g.CurrentClimat.Meteo, g.CurrentClimat.Temperature, g.CurrentClimat.Humidite, g.CurrentClimat.Co2, g.CurrentClimat.O2)

		// 计算文本的位置
		x := float64(screenWidth) // 文本的起始 X 坐标
		y := 10.0                 // 文本的起始 Y 坐标

		// 在屏幕右侧绘制文本
		ebitenutil.DebugPrintAt(screen, climatInfo, int(x), int(y))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//return screenWidth, screenHeight
	return screenWidth + menuBarWidth, screenHeight
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
	}

	g.initMenuBarAndButton()

	ebiten.SetWindowSize((screenWidth+menuBarWidth)*2, screenHeight*2)
	ebiten.SetWindowTitle("Multi-agent system")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
