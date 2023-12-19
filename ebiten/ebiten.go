package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"sort"
	"time"
	server "vivarium"
	"vivarium/climat"
	"vivarium/enums"

	"vivarium/ebiten/assets/images"
	sprite "vivarium/ebiten/sprites"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth        = 272
	screenHeight       = 272
	tileSize           = 16
	menuBarWidth       = 100
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
	updateInterval int
	updateCount    int
	SpriteMap      map[int]*sprite.Sprite

	// 按钮
	isPaused    bool
	buttonImage *ebiten.Image
	buttonRect  image.Rectangle

	buttonUnpressedImage *ebiten.Image
	buttonPressedImage   *ebiten.Image

	// Add fields for display
	CurrentHour   int
	CurrentClimat *climat.Climat

	meteoFrames       map[enums.Meteo][]*ebiten.Image
	meteoIndex        map[enums.Meteo]int
	randomCoordinates [][2]int
}

func (g *Game) loadMeteoFrames() {
	g.meteoFrames = make(map[enums.Meteo][]*ebiten.Image)

	g.meteoIndex = make(map[enums.Meteo]int)

	// // Load rain images
	rainFrames := make([]*ebiten.Image, 0)
	img1, _, err := image.Decode(bytes.NewReader(images.Rain_png))
	if err != nil {
		log.Fatal(err)
	}
	rainImage := ebiten.NewImageFromImage(img1)
	for i := 0; i < 3; i++ {
		rainFrames = append(rainFrames, rainImage.SubImage(image.Rect(0, i*272, 272, i*272+272)).(*ebiten.Image))
	}
	g.meteoFrames[enums.Pluie] = rainFrames

	// // Load fog images
	fogFrames := make([]*ebiten.Image, 1)
	img2, _, err := image.Decode(bytes.NewReader(images.Fog_png))
	if err != nil {
		log.Fatal(err)
	}
	fogFrames[0] = ebiten.NewImageFromImage(img2)
	g.meteoFrames[enums.Brouillard] = fogFrames

	// // Load dry season images
	// drySeasonFrames := make([]*ebiten.Image, 0)
	// for i := 1; i <= 4; i++ {
	// 	img, _, err := image.Decode(bytes.NewReader(images.DrySeason_png))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	drySeasonFrames = append(drySeasonFrames, ebiten.NewImageFromImage(img))
	// }
	// g.meteoFrames[enums.SaisonSeche] = drySeasonFrames

	// Load fire images
	img4, _, err := image.Decode(bytes.NewReader(images.Fire_png))
	if err != nil {
		log.Fatal(err)
	}
	fireImage := ebiten.NewImageFromImage(img4)
	g.meteoFrames[enums.Incendie] = sprite.LoadFrames(fireImage, 5, 0)

	// Load thunder images
	img5, _, err := image.Decode(bytes.NewReader(images.Thunder_png))
	if err != nil {
		log.Fatal(err)
	}
	thunderImage := ebiten.NewImageFromImage(img5)
	g.meteoFrames[enums.Tonnerre] = sprite.LoadFrames(thunderImage, 5, 0)
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
	g.buttonRect = image.Rect(160+screenWidth/2, 120, 200+screenWidth/2, 150)
}

func (g *Game) Update() error {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= g.buttonRect.Min.X && x <= g.buttonRect.Max.X && y >= g.buttonRect.Min.Y && y <= g.buttonRect.Max.Y {
			g.isPaused = !g.isPaused
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

	// 60 frames per second, every 30 frames is 0.5 seconds
	g.updateInterval = 30

	// Count for updating
	g.updateCount++

	// Updated every 0.5s
	if g.updateCount >= g.updateInterval {

		if server.EcosystemForEbiten == nil {
			return nil
		}

		// Read data from server
		ecosystemForEbiten := server.EcosystemForEbiten

		for _, org := range ecosystemForEbiten.GetAllOrganisms() {
			if _, exists := g.SpriteMap[org.GetID()]; !exists {
				if org.GetEspece().String() == "AraignéeSauteuse" {
					g.SpriteMap[org.GetID()] = sprite.NewSpiderSprite(org)
				} else if org.GetEspece().String() == "PetitSerpent" {
					g.SpriteMap[org.GetID()] = sprite.NewCobraSprite(org)
				} else if org.GetEspece().String() == "Grillons" {
					g.SpriteMap[org.GetID()] = sprite.NewScarabSprite(org)
				} else if org.GetEspece().String() == "Escargot" {
					g.SpriteMap[org.GetID()] = sprite.NewSnailSprite(org)
				} else if org.GetEspece().String() == "Lombric" {
					g.SpriteMap[org.GetID()] = sprite.NewWormSprite(org)
				} else if org.GetEspece().String() == "Champignon" {
					g.SpriteMap[org.GetID()] = sprite.NewMushroomSprite(org)
				} else if org.GetEspece().String() == "PetitHerbe" {
					g.SpriteMap[org.GetID()] = sprite.NewPetitHerbeSprite(org)
				} else if org.GetEspece().String() == "GrandHerbe" {
					g.SpriteMap[org.GetID()] = sprite.NewGrandHerbeSprite(org)
				} else {
					fmt.Println("Error: Unknown Espece", org.GetEspece().String())
				}
			} else {
				if g.SpriteMap[org.GetID()].IsDead {
					//delete(g.SpriteMap, org.GetID())
					continue
				}
				// Update the information of sprite
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

		// Update data from server.EcosystemForEbiten
		ecosystemData := server.EcosystemForEbiten
		if ecosystemData != nil {
			g.CurrentHour = ecosystemData.Hour
			g.CurrentClimat = ecosystemData.Climat
		}

		// Reset counter
		g.updateCount = 0
	}

	for _, sprite := range g.SpriteMap {
		if sprite.IsDead {
			continue
		}
		sprite.Update(deltaTime)
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

// generateRandomCoordinate  generates a random coordinate
func generateRandomCoordinate() (int, int) {
	return rand.Intn(15) + 1, rand.Intn(15) + 1
}

func (g *Game) DrawWeather(screen *ebiten.Image) {
	if g.CurrentClimat != nil {
		switch g.CurrentClimat.Meteo {
		case enums.Pluie:
			// Draw rain
			if g.meteoIndex[enums.Pluie] < len(g.meteoFrames[enums.Pluie])*100 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(0, 0)
				screen.DrawImage(g.meteoFrames[enums.Pluie][(g.meteoIndex[enums.Pluie]/10)%len(g.meteoFrames[enums.Pluie])], op)
				g.meteoIndex[enums.Pluie]++
			} else {
				g.meteoIndex[enums.Pluie] = 0
				//g.CurrentClimat.Meteo = enums.Rien
			}
		case enums.Brouillard:
			// Draw fog
			if g.meteoIndex[enums.Brouillard] < len(g.meteoFrames[enums.Brouillard])*100 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(0, 0)
				screen.DrawImage(g.meteoFrames[enums.Brouillard][0], op)
				g.meteoIndex[enums.Brouillard]++
			} else {
				g.meteoIndex[enums.Brouillard] = 0
				//g.CurrentClimat.Meteo = enums.Rien
			}
		case enums.SaisonSeche:
			// ToDo： Draw dry season
		case enums.Incendie:
			// Draw fire
			if g.meteoIndex[enums.Incendie] == 0 {
				rand.Seed(time.Now().UnixNano())
				numCoordinates := 100
				for i := 0; i < numCoordinates; i++ {
					x, y := generateRandomCoordinate()
					g.randomCoordinates = append(g.randomCoordinates, [2]int{x, y})
				}
			}

			if g.meteoIndex[enums.Incendie] < len(g.meteoFrames[enums.Incendie])*100 {
				for i := 0; i < len(g.randomCoordinates); i++ {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(g.randomCoordinates[i][0]*16), float64(g.randomCoordinates[i][1]*16))
					screen.DrawImage(g.meteoFrames[enums.Incendie][(g.meteoIndex[enums.Incendie]/10)%len(g.meteoFrames[enums.Incendie])], op)
				}
				g.meteoIndex[enums.Incendie]++
			} else {
				g.randomCoordinates = nil
				g.meteoIndex[enums.Incendie] = 0
				//g.CurrentClimat.Meteo = enums.Rien
			}

		case enums.Tonnerre:
			// Draw thunder
			if g.meteoIndex[enums.Tonnerre] == 0 {
				rand.Seed(time.Now().UnixNano())
				numCoordinates := 100
				for i := 0; i < numCoordinates; i++ {
					x, y := generateRandomCoordinate()
					g.randomCoordinates = append(g.randomCoordinates, [2]int{x, y})
				}
			}

			if g.meteoIndex[enums.Tonnerre] < len(g.meteoFrames[enums.Tonnerre])*100 {
				for i := 0; i < len(g.randomCoordinates); i++ {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(g.randomCoordinates[i][0]*16), float64(g.randomCoordinates[i][1]*16))
					screen.DrawImage(g.meteoFrames[enums.Tonnerre][(g.meteoIndex[enums.Tonnerre]/10)%len(g.meteoFrames[enums.Tonnerre])], op)
				}
				g.meteoIndex[enums.Tonnerre]++
			} else {
				g.randomCoordinates = nil
				g.meteoIndex[enums.Tonnerre] = 0
				//g.CurrentClimat.Meteo = enums.Rien
			}
		case enums.Rien:
			// Draw nothing
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.DrawBackground(screen)

	// Define priority mapping
	priorityMap := map[string]int{
		"PetitHerbe":       1,
		"GrandHerbe":       2,
		"Champignon":       3,
		"Escargot":         4,
		"Grillons":         5,
		"Lombric":          6,
		"PetitSerpent":     7,
		"AraignéeSauteuse": 8,
	}

	// Create a sprite slice sorted by priority
	sortedSprites := make([]*sprite.Sprite, 0, len(g.SpriteMap))
	for _, sprite := range g.SpriteMap {
		sortedSprites = append(sortedSprites, sprite)
	}

	// Use custom sorting
	sort.Slice(sortedSprites, func(i, j int) bool {
		return priorityMap[sortedSprites[i].Species] < priorityMap[sortedSprites[j].Species]
	})

	// Loop through the sorted sprites and draw them
	for _, sprite := range sortedSprites {
		sprite.Draw(screen, g.FrameIndex)
	}

	// Draw button
	buttonOp := &ebiten.DrawImageOptions{}
	buttonOp.GeoM.Translate(float64(g.buttonRect.Min.X), float64(g.buttonRect.Min.Y))
	if g.isPaused {
		screen.DrawImage(g.buttonUnpressedImage, buttonOp)
	} else {
		screen.DrawImage(g.buttonPressedImage, buttonOp)
	}

	// Show datas of climat and hour
	if g.CurrentClimat != nil {
		climatInfo := fmt.Sprintf("Hour: %d\nMeteo: %s\nTemperature: %d°C\nHumidity: %.2f%%\nCO2: %.2f%%\nO2: %.2f%%",
			g.CurrentHour, g.CurrentClimat.Meteo, g.CurrentClimat.Temperature, g.CurrentClimat.Humidite, g.CurrentClimat.Co2, g.CurrentClimat.O2)

		// Calculate text position
		x := float64(screenWidth) // The starting X coordinate of the text
		y := 10.0                 // The starting Y coordinate of the text

		// Draw text on the right side of the screen
		ebitenutil.DebugPrintAt(screen, climatInfo, int(x), int(y))
	}

	g.DrawWeather(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth + menuBarWidth, screenHeight
}

func main() {

	// Launch server
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
	}

	g.initMenuBarAndButton()
	g.loadMeteoFrames()

	ebiten.SetWindowSize((screenWidth+menuBarWidth)*2, screenHeight*2)
	ebiten.SetWindowTitle("Multi-agent system")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
