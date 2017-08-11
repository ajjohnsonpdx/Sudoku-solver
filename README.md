# Sudoku Solver in Go

I wrote this Solver mostly as an exercise to teach myself Go.

It uses a basic algorthim that first determines if you can set a cell's value by checking if there's only one possible value based on the other known values in it's cell group (the group of cells in it's row, column, and square).  Then it determines if a cell's value can be fixed based on all of the possible values of the other cells in it's cell groups.  If any of it's possible values do not exist in the set of possible values in each of it's cell groups, by process of elimination, that value is a possible solution.  Finally, do a depth-first recursive search of all possible solutions, eliminating non-solutions using the first 2 principles.


# Running the Solver

    $ go run sudoku-solver.go

### Entering the Board

Boards are a set of 9x9 cells, with known values, where 0 represents blank cells.  Cell values are separated by a period (`.`) 

# Board files

Some sample boards are included. 
