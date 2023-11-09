package environnement

// Terrain represents the spatial layout of the environment, where each cell contains a list of Organisme IDs.
type Terrain struct {
	Width, Height int       // Dimensions of the terrain
	Grid          [][][]int // A 3D slice where each cell contains a list of Organisme IDs
}

// NewTerrain creates a new Terrain with the specified width and height, initializing the grid with empty lists.
func NewTerrain(width, height int) *Terrain {
	grid := make([][][]int, height)
	for i := range grid {
		grid[i] = make([][]int, width)
		for j := range grid[i] {
			grid[i][j] = []int{} // Initialize each cell with an empty list
		}
	}
	return &Terrain{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

// AddOrganism adds an Organisme ID to a specific location in the grid.
func (t *Terrain) AddOrganism(id int, x, y int) {
	// Assuming x and y are within the bounds of the grid:
	t.Grid[y][x] = append(t.Grid[y][x], id)
}

// RemoveOrganism removes an Organisme ID from a specific location in the grid.
func (t *Terrain) RemoveOrganism(id int, x, y int) {
	// Find and remove the Organisme ID from the cell at x, y
	organismsAtCell := t.Grid[y][x]
	for i, oid := range organismsAtCell {
		if oid == id {
			t.Grid[y][x] = append(organismsAtCell[:i], organismsAtCell[i+1:]...)
			break
		}
	}
}

// UpdateOrganismPosition updates the position of an organism, moving its ID from one cell to another.
func (t *Terrain) UpdateOrganismPosition(id int, oldX, oldY, newX, newY int) {
	// Remove the Organisme ID from the old location
	t.RemoveOrganism(id, oldX, oldY)
	// Add the Organisme ID to the new location
	t.AddOrganism(id, newX, newY)
}

// Other methods for managing organisms on the terrain would be defined here.
