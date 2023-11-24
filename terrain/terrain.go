package terrain

import "vivarium/enums"

type CellInfo struct {
	OrganismID   int
	OrganismType string // "Escargot", "Grillons", "Lombric", "PetitSerpent", "AraignéeSauteuse", "PetitHerbe", "GrandHerbe", "Champignon"
}

type Terrain struct {
	Width, Length int
	Grid          [][][]CellInfo // Updated to store CellInfo
	CurrentHour   int
	Meteo         enums.Meteo
}

func NewTerrain(width, length int) *Terrain {
	grid := make([][][]CellInfo, length)
	for i := range grid {
		grid[i] = make([][]CellInfo, width)
		for j := range grid[i] {
			grid[i][j] = []CellInfo{} // Initialize each cell with an empty slice of CellInfo
		}
	}
	return &Terrain{
		Width:  width,
		Length: length,
		Grid:   grid,
	}
}

func (t *Terrain) AddOrganism(id int, organismType string, x, y int) {
	t.Grid[y][x] = append(t.Grid[y][x], CellInfo{OrganismID: id, OrganismType: organismType})
}

func (t *Terrain) RemoveOrganism(id int, x, y int) {
	for i, cell := range t.Grid[y][x] {
		if cell.OrganismID == id {
			t.Grid[y][x] = append(t.Grid[y][x][:i], t.Grid[y][x][i+1:]...)
			break
		}
	}
}

func (t *Terrain) UpdateOrganismPosition(id int, organismType string, oldX, oldY, newX, newY int) {
	t.RemoveOrganism(id, oldX, oldY)
	t.AddOrganism(id, organismType, newX, newY)
}
