# Ant Pathfinding and Distribution Project

This project enables ants to find the shortest path from a starting point to an endpoint and distributes them along the paths. The project includes various functions and unit tests.

## Project Structure

The project consists of the following files and folder:

- `main.go`: Main file. Contains the main function and the functions used in the project.
- `main_test.go`: Unit test file. Contains tests to ensure the functions work correctly.
- `functions.go` : Contains the functions used in the main file.
- `errors_handler.go` : Contains the error handling functions.
- `maps`: Contains the example files used in the project.

## Usage

### Functions in `main.go`

1. `main`: Reads the file, finds the paths, and distributes the ants.

### Functions in `main_test.go`

1. `TestFindConnection`: Tests the `find_connection` function.
2. `TestLoopHandler`: Tests the `loop_handler` function.
3. `TestDfsPaths`: Tests the `find_road_recursive` function.
4. `TestFindAllPaths`: Tests the `find_all_paths` function.
5. `TestIsOverlapping`: Tests the `is_overlapping` function.
6. `TestFindMaxNonOverlappingPaths`: Tests the `find_max_non_overlapping_paths` function.
7. `TestSortPathsByLength`: Tests the `sort_paths_by_length` function.
8. `TestDispatchAnts`: Tests the `dispatch_ants` function and verifies the outputs.

### Functions in `functions.go`

1. `find_connection`: Finds the connections of a room.
2. `loop_handler`: Checks if a room is in a path.
3. `find_road_recursive`: Uses Breadth-First Search (BFS) algorithm to find all paths from the start room to the end room.
4. `find_all_paths`: Finds all paths and selects the maximum number of non-overlapping paths.
5. `find_max_non_overlapping_paths`: Finds the maximum number of non-overlapping paths.
6. `is_overlapping`: Checks if two paths overlap.
7. `sort_paths_by_length`: Sorts paths by their lengths.
8. `dispatch_ants`: Distributes and moves ants along the paths.

### Functions in `errors_handler.go`

1. `Contains_Error`: Checks if input file contains wrong format.
2. `True_Format_Error`: Formats the error message.

## Setup and Running

### Requirements

- Go (Golang) 1.16 or higher

### Steps

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>

### Usage

```
go run . <file_name>
```

## Example File Format

```
<number_of_ants>
##start
<start_room_name> <x_coordinate> <y_coordinate>
...
##end
<end_room_name> <x_coordinate> <y_coordinate>
...
<connection_room1-room2>
...
```

## Example

```
20
##start
0 2 0
1 4 1
2 6 0
##end
3 5 3
0-1
0-3
1-2
3-2
```

## Output

```
L1-3 L2-1 L3-3 L4-3 L6-3 L7-3 L9-3 L10-3 L12-3 L13-3 L15-3 L16-3 L18-3 L19-3
L2-2 L5-1
L2-3 L5-2 L8-1
L5-3 L8-2 L11-1
L8-3 L11-2 L14-1
L11-3 L14-2 L17-1
L14-3 L17-2 L20-1
L17-3 L20-2
L20-3
```

