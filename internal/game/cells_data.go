package game

import "fmt"

func (s *Session) newCellsData() *CellsData {
	data := &CellsData{
		rootCount: s.cfg.RootCount,
		rows:      [][]*Cell{},
	}

	for i := 0; i < s.cfg.RootCount; i++ {
		row := make([]*Cell, s.cfg.RootCount)

		for j := 0; j < s.cfg.RootCount; j++ {
			row[j] = &Cell{
				Value:     "",
				Highlight: false,
			}
		}

		data.rows = append(data.rows, row)
	}

	return data
}

func (cd *CellsData) cellExists(x, y int) bool {
	if x < 0 || x >= len(cd.rows) {
		return false
	}

	if y < 0 || y >= len((cd.rows)[x]) {
		return false
	}

	return true
}

func (cd *CellsData) getCell(x, y int) *Cell {
	if !cd.cellExists(x, y) {
		return nil
	}

	return cd.rows[x][y]
}

func (cd *CellsData) setValue(x, y int, value string) error {
	if !cd.cellExists(x, y) {
		return fmt.Errorf("cell %d/%d doesn't exist", x, y)
	}

	cd.rows[x][y].Value = value

	return nil
}

func (cd *CellsData) setHighlight(x, y int, status bool) error {
	if !cd.cellExists(x, y) {
		return fmt.Errorf("cell %d/%d doesn't exist", x, y)
	}

	cd.rows[y][x].Highlight = status

	return nil
}

// x and y tell the position of the cell to start pattern checking from.
func (cd *CellsData) getPattern(x, y int) [][]int {
	cell := cd.getCell(x, y)

	if cell == nil || cell.Value == "" {
		return nil
	}

	// Todo: implment direction check.

	return nil
}
