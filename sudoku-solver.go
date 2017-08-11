package main

import (
   "fmt"
   "bufio"
   "os"
   "strings"
   "strconv"
)

type Cell struct {
	index int
	fixed bool
	value int
	column int
	row int
	possibles []int
}

func main() {
	boardInput := InputBoard()
        cellvals := strings.Split(boardInput,".")
	var board = [81]Cell{}        
	BuildInitialBoard(cellvals, &board)	
	PrintBoard(&board)
	SetAllCellsPossibles(&board)
	SolveBoard(&board)
}

func InputBoard() string {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Create a sudoku board by entering cell values separated by a '.' Use 0 for empty cells: ")
        boardInput, _ := reader.ReadString('\n')
	boardInput = strings.Replace(boardInput, "\n", "", -1)
	return boardInput
}

func BuildInitialBoard(cellVals []string, cells *[81]Cell) {
	index := 0
	for i := range cellVals {
		num, _ := strconv.Atoi(cellVals[i])
		p := make([]int, 0, 9)
		row := i/9 
		column := i - (row * 9)
		if num != 0 {
			s := Cell{index, true, num, column, row, p}
			cells[i] = s
		} else {
			s := Cell{index, false, 0, column, row, p}
			cells[i] = s
		}
	index++
	}
} 

func SolveBoard(cells *[81]Cell) bool {
	//If there's only one possible value for a cell, fix it's value
	//Next, try to find a cell's value by elmination based on it's 
	//cell group possibles
	//Finally do a depth-first search using all of the possibles

	rollBack := make([]int, 0, 0)	
	for {
		var unchanged = true
        	for i := range cells {
        		if !(cells[i].fixed) && len(cells[i].possibles) == 1 {				
				rollBack = append(rollBack, i)
                		SetFixedValue(cells[i].possibles[0], &cells[i])
				if !SetAllCellsPossibles(cells) {
					Rollback(rollBack, cells)
					return false
				}
				unchanged = false
        		}
		}
		if unchanged == true {
			for i := range cells {
				if !(cells[i].fixed) && EliminateByGroupPossibles(i, cells) {
					rollBack = append(rollBack, i)
					if !SetAllCellsPossibles(cells) {
						Rollback(rollBack, cells)
						return false
					}
                        		unchanged = false
					break
				}
			}	
		}

		if unchanged == true {
			break
		}
	}
	
	if IsSolved(cells) {
		fmt.Println("Solved")
		PrintBoard(cells)
		
		return true
	}	

	for i := range cells {
		if !cells[i].fixed {
			for p := range cells[i].possibles {
				if Search(i, p, cells) {
					return true
				}
			}
		}
	}
	Rollback(rollBack, cells)
	return false
}

func Search (cellIndex int, possIndex int, cells *[81]Cell) bool {
	//Try a possible value, and see if a solution exists
	//If any of the cells have no possibles, a solution is not possible
	
	possibles := cells[cellIndex].possibles
	SetFixedValue(cells[cellIndex].possibles[possIndex], &cells[cellIndex])
	if !SetAllCellsPossibles(cells) {
		cells[cellIndex].fixed = false
		cells[cellIndex].value = 0
		cells[cellIndex].possibles = possibles
		return false
	}	
	
	if IsSolved(cells) {
		fmt.Println("Solved")
		PrintBoard(cells)
		return true
	}	
	
	if SolveBoard(cells) {
		return true
	} else {
		cells[cellIndex].fixed = false
		cells[cellIndex].value = 0
		cells[cellIndex].possibles = possibles
		SetAllCellsPossibles(cells)
		return false
	}
	
}	

func IsSolved(cells *[81]Cell) bool {
	for i := range cells {
		if !cells[i].fixed {
			 return false
		}
	}
	return true
}

func PrintBoard(cells *[81]Cell) {
	for i := 0; i < 9; i++ {
		for p := 0; p < 9; p++ {
			fmt.Print(cells[(i*9)+p].value, " ")
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func Rollback(rollBacks []int, cells *[81]Cell) {
	//If a cell's value was set, roll it back
	//Recalculate all the possibles when finished

	for i := range rollBacks {
		cells[rollBacks[i]].fixed = false
		cells[rollBacks[i]].value = 0
	}
	SetAllCellsPossibles(cells)
}

func SetAllCellsPossibles(cells *[81]Cell) bool {
	var numPossibles int
	for i := range cells {
		if !cells[i].fixed {
			numPossibles = SetCellPossibles(i, cells)
			if numPossibles == 0 {return false}
		}
	}
	return true
}


func SetCellPossibles(cellIndex int, cells *[81]Cell) int {
	//Get all the fixed values in the row/column/square groups
	//Any value not in that set is a possible value for the cell

	var cellGroup = [9]Cell{}
	cellGroup = GetAllCellsByColumn(cells[cellIndex].column, cells)
	fixedColValues := GetFixedValues(&cellGroup)

	cellGroup = GetAllCellsByRow(cells[cellIndex].row, cells)
	fixedRowValues := GetFixedValues(&cellGroup)
	
	mergedFixedRCValues := Merge(fixedColValues, fixedRowValues)
	
	cellGroup = GetAllCellsBySquare(cellIndex, cells)
	fixedSquareValues := GetFixedValues(&cellGroup)
	
	mergedFixedValues := Merge(mergedFixedRCValues, fixedSquareValues)
	

	possibles := make([]int, 0, 0)
	for i := 1; i < 10; i++ {
		if !IsMember(i, mergedFixedValues) {
			possibles = append(possibles, i)
		}
	} 	
	cells[cellIndex].possibles = possibles
	return len(possibles)
}

func GetFixedValues(cellGroup *[9]Cell) []int { 
	//Get all the fixed values in a row/column/square group

	fixedValues := make([]int, 0, 0)
	for i := range cellGroup {
		if cellGroup[i].fixed == true {
			fixedValues = append(fixedValues, cellGroup[i].value)
		}
	}
	return fixedValues
}
func SetFixedValue(value int, cell *Cell) {
	cell.fixed = true
	cell.value = value
	s := make([]int,0,0)
	cell.possibles = s
}

func GetCellGroupPossibles(cellGroup []Cell) []int {
	//Return all the possibles in a row/column/square group

	possibles := make([]int, 0, 0)
        for i := range cellGroup {
                if cellGroup[i].fixed == false {
                        for p := range cellGroup[i].possibles {
				possibles = append(possibles, cellGroup[i].possibles[p])
			}
                }
        }
	return possibles
}

func GetAllCellsByColumn(column int, cells *[81]Cell) [9]Cell {
	//Return all cells in a column group

	var p int = 0 
	var columnCells = [9]Cell{}
	for i := column; i < 81; i+=9 {
		columnCells[p] = cells[i]
		p++
	}
	return columnCells
}

func GetAllCellsByRow(row int, cells *[81]Cell) [9]Cell {
	//Return all cells in a row group

	var p int = (row * 9) 
	var rowCells = [9]Cell{}
	for i := 0; i < 9; i++ {
		rowCells[i] = cells[p]
		p++
	}
	return rowCells
}

func GetAllCellsBySquare(index int, cells *[81]Cell) [9]Cell {
	//Return all cells in the square group

	var squareCells = [9]Cell{}
	
	switch index {
	case 0, 1, 2, 9, 10, 11, 18, 19, 20:
		squareCells[0] = cells[0]
		squareCells[1] = cells[1]
		squareCells[2] = cells[2]
		squareCells[3] = cells[9]
		squareCells[4] = cells[10]
		squareCells[5] = cells[11]
		squareCells[6] = cells[18]
		squareCells[7] = cells[19]
		squareCells[8] = cells[20]
	
	case 3, 4, 5, 12, 13, 14, 21, 22, 23:
		squareCells[0] = cells[3]
		squareCells[1] = cells[4]
		squareCells[2] = cells[5]
		squareCells[3] = cells[12]
		squareCells[4] = cells[13]
		squareCells[5] = cells[14]
		squareCells[6] = cells[21]
		squareCells[7] = cells[22]
		squareCells[8] = cells[23]
	
	case 6, 7, 8, 15, 16, 17, 24, 25, 26:
		squareCells[0] = cells[6]
		squareCells[1] = cells[7]
		squareCells[2] = cells[8]
		squareCells[3] = cells[15]
		squareCells[4] = cells[16]
		squareCells[5] = cells[17]
		squareCells[6] = cells[24]
		squareCells[7] = cells[25]
		squareCells[8] = cells[26]
	
	case 27, 28, 29, 36, 37, 38, 45, 46, 47:
		squareCells[0] = cells[27]
		squareCells[1] = cells[28]
		squareCells[2] = cells[29]
		squareCells[3] = cells[36]
		squareCells[4] = cells[37]
		squareCells[5] = cells[38]
		squareCells[6] = cells[45]
		squareCells[7] = cells[46]
		squareCells[8] = cells[47]
	
	case 30, 31, 32, 39, 40, 41, 48, 49, 50:
		squareCells[0] = cells[30]
		squareCells[1] = cells[31]
		squareCells[2] = cells[32]
		squareCells[3] = cells[39]
		squareCells[4] = cells[40]
		squareCells[5] = cells[41]
		squareCells[6] = cells[48]
		squareCells[7] = cells[49]
		squareCells[8] = cells[50]
	
	case 33, 34, 35, 42, 43, 44, 51, 52, 53:
		squareCells[0] = cells[33]
		squareCells[1] = cells[34]
		squareCells[2] = cells[35]
		squareCells[3] = cells[42]
		squareCells[4] = cells[43]
		squareCells[5] = cells[44]
		squareCells[6] = cells[51]
		squareCells[7] = cells[52]
		squareCells[8] = cells[53]
	
	case 54, 55, 56, 63, 64, 65, 72, 73, 74:
		squareCells[0] = cells[54]
		squareCells[1] = cells[55]
		squareCells[2] = cells[56]
		squareCells[3] = cells[63]
		squareCells[4] = cells[64]
		squareCells[5] = cells[65]
		squareCells[6] = cells[72]
		squareCells[7] = cells[73]
		squareCells[8] = cells[74]
	
	case 57, 58, 59, 66, 67, 68, 75, 76, 77:
		squareCells[0] = cells[57]
		squareCells[1] = cells[58]
		squareCells[2] = cells[59]
		squareCells[3] = cells[66]
		squareCells[4] = cells[67]
		squareCells[5] = cells[68]
		squareCells[6] = cells[75]
		squareCells[7] = cells[76]
		squareCells[8] = cells[77]
	
	case 60, 61, 62, 69, 70, 71, 78, 79, 80:
		squareCells[0] = cells[60]
		squareCells[1] = cells[61]
		squareCells[2] = cells[62]
		squareCells[3] = cells[69]
		squareCells[4] = cells[70]
		squareCells[5] = cells[71]
		squareCells[6] = cells[78]
		squareCells[7] = cells[79]
		squareCells[8] = cells[80]

	}
	return squareCells
}

func EliminateByGroupPossibles(cellIndex int, cells *[81]Cell) bool {
	//Look at the possibles for each of the three cell groups for a cell
	//and see if we can fix the value by process of elimination

	var cellGroup = [9]Cell{}
	cellGroup = GetAllCellsByColumn(cells[cellIndex].column, cells)
	newCellGroupC := RemoveCellFromGroup(cellIndex, cellGroup)
	columnPossibles := GetCellGroupPossibles(newCellGroupC)
	if (CheckPossiblesForElmination(&cells[cellIndex], columnPossibles)) {
		return true
	}
	
	cellGroup = GetAllCellsByRow(cells[cellIndex].row, cells)
        newCellGroupR := RemoveCellFromGroup(cellIndex, cellGroup)
        rowPossibles := GetCellGroupPossibles(newCellGroupR)
	if (CheckPossiblesForElmination(&cells[cellIndex], rowPossibles)) {
                return true
        }

	cellGroup = GetAllCellsBySquare(cellIndex, cells)
        newCellGroups := RemoveCellFromGroup(cellIndex, cellGroup)
        squarePossibles := GetCellGroupPossibles(newCellGroups)
        if (CheckPossiblesForElmination(&cells[cellIndex], squarePossibles)) {
                return true
        }
	
	return false	
}

func CheckPossiblesForElmination(cell *Cell, possibles []int) bool {
	//Given a set of possible values (from a cell group)
	//If any of the cell's possible values don't exist in that set,
	//then fix the cells's value to that number

	arr := GetUniques(cell.possibles, possibles)
	if len(arr) == 1 {
		SetFixedValue(arr[0], cell)
		return true
	}
	return false
}

func RemoveCellFromGroup(cellTargetIndex int, cells [9]Cell) []Cell {
	//Remove a cell from a cell group given an index

	remainingCells := make([]Cell, 0, 0)
        for p := range cells {
                if cells[p].index != cellTargetIndex {
                        remainingCells = append(remainingCells, cells[p])
                }
        }
	return remainingCells
}

func Merge(s1, s2 []int) []int {
	//Merge two arrays and return unique values
	 
       	unique := make(map[int]struct{})

        for _, v := range s1 {
                unique[v] = struct{}{}
        }
        for _, v := range s2 {
                unique[v] = struct{}{}
        }
        final := make([]int, len(unique))
        i := 0
        for k := range unique {
                final[i] = k
                i++ 
        }
        return final
}

func GetUniques(obj []int, target []int) []int {
        //Return unique values of obj that exist in target
       
	 returnArr := make([]int, 0, 0)
        found := false
        for i := 0; i < len(obj); i++ {
                for p :=0; p < len(target); p++ {
                        if (obj[i] == target[p]) {
                                found = true
                        }
                }
                if !found {
                        returnArr = append(returnArr, obj[i])
                        
                }
		found = false

        }
        return returnArr
}

func IsMember(obj int, arr []int) bool {
        for i :=0; i < len(arr); i++ {
                if (arr[i] == obj) {return true}
        }
        return false
}
