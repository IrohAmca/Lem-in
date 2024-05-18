package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	content      string
	start_room   string
	end_room     string
	connect_rows []string
	comment_rows []string
	roads        [][]string
	ant_count    int
	rows         [][]string
	step         int
)

func seperate_rows() {
	var temp string
	for _, row := range strings.Split(content, "\n") {
		if row != "" {
			temp = strings.TrimSpace(row)
			rows = append(rows, strings.Split(temp, " "))
		}
	}
}

func save_data() {
	seperate_rows()
	for i, row := range rows {
		if len(row) == 1  && row[0] != "##start" && row[0] != "##end"{
			count, err := strconv.Atoi(row[0])
			if err == nil {
				ant_count = count
			}
		}
		if row[0] == "##start" {
			if len(rows[i+1]) == 3 {
				start_room = rows[i+1][0]
			}
		}
		if row[0] == "##end" {
			if len(rows[i+1]) == 3 {
				end_room = rows[i+1][0]
			}
		}
		if len(row) == 3 && row[0] != start_room && row[0] != end_room {
			comment_rows = append(comment_rows, row[0])
		}
		if len(row) == 1 && strings.Contains(row[0], "-") {
			connect_rows = append(connect_rows, row[0])
		}
	}
	ErrorHandler(nil)
}
func ErrorHandler(err error) {
}
func read_file(file_path string) {
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	content = string(data)
}

func find_connection(current_room string) []string {
	selected_rooms := []string{}
	for _, row := range connect_rows {
		if strings.Contains(row, current_room+"-") || strings.Contains(row, "-"+current_room) {
			words := strings.Split(row, "-")
			for _, word := range words {
				if word != current_room {
					selected_rooms = append(selected_rooms, word)
				}
			}
		}
	}
	return selected_rooms
}

func loop_handler(room string, road []string) bool {
	for _, r := range road {
		if r == room {
			return false
		}
	}
	return true
}

func bfs_paths(start_room string) [][]string {
	queue := [][]string{{start_room}}
	var paths [][]string

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		room := path[len(path)-1]

		if room == end_room {
			paths = append(paths, path)
			continue
		}

		for _, next_room := range find_connection(room) {
			if loop_handler(next_room, path) {
				new_path := make([]string, len(path))
				copy(new_path, path)
				new_path = append(new_path, next_room)
				queue = append(queue, new_path)
			}
		}
	}
	return paths
}

func find_all_paths() {
	all_paths := bfs_paths(start_room)
	roads = find_max_non_overlapping_paths(all_paths)
}

func find_max_non_overlapping_paths(paths [][]string) [][]string {
	max_paths := [][]string{}

	var backtrack func(current_paths [][]string, index int)
	backtrack = func(current_paths [][]string, index int) {
		if index == len(paths) {
			if len(current_paths) > len(max_paths) {
				max_paths = append([][]string(nil), current_paths...)
			}
			return
		}

		if !is_overlapping(paths[index], current_paths) {
			backtrack(append(current_paths, paths[index]), index+1)
		}
		backtrack(current_paths, index+1)
	}

	backtrack([][]string{}, 0)
	return max_paths
}

func is_overlapping(path []string, paths [][]string) bool {
	for _, p := range paths {
		for i := 1; i < len(path)-1; i++ {
			for j := 1; j < len(p)-1; j++ {
				if path[i] == p[j] {
					return true
				}
			}
		}
	}
	return false
}

func sort_paths_by_length(paths [][]string) {
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if len(paths[i]) > len(paths[j]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}
}

func dispatch_ants() {
	find_all_paths()
	sort_paths_by_length(roads)

	ant_paths := make([][]string, ant_count)
	path_usage := make([]int, len(roads))

	for i := 0; i < ant_count; i++ {
		min_usage_index := 0
		min_usage_value := path_usage[0] * len(roads[0])

		for j := 1; j < len(roads); j++ {
			usage_value := path_usage[j] * len(roads[j])
			if usage_value < min_usage_value {
				min_usage_index = j
				min_usage_value = usage_value
			}
		}

		ant_paths[i] = roads[min_usage_index]
		path_usage[min_usage_index]++
	}

	ant_positions := make(map[int]int)
	for i := 0; i < ant_count; i++ {
		ant_positions[i] = 0
	}

	step = 0
	for {
		moves := []string{}
		occupied_rooms := map[string]bool{}
		for i := 0; i < ant_count; i++ {
			if ant_positions[i] < len(ant_paths[i])-1 {
				next_position := ant_positions[i] + 1
				next_room := ant_paths[i][next_position]

				if !occupied_rooms[next_room] || next_room == end_room {
					occupied_rooms[next_room] = true
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, next_room))
					ant_positions[i] = next_position
				}
			}
		}

		if len(moves) == 0 {
			break
		}

		fmt.Println(strings.Join(moves, " "))
		step++
	}
}

func main() {
	read_file(os.Args[1])
	save_data()
	fmt.Println("Ant Count", ant_count)
	fmt.Println("Start Room", start_room)
	fmt.Println("Comment Rooms", comment_rows)
	fmt.Println("Connect Room", connect_rows)
	fmt.Println("End Room", end_room)
	dispatch_ants()
	fmt.Println("Step:", step)
	fmt.Println("Paths:")
	for _, road := range roads {
		fmt.Println(road)
	}
}

/*
L1-3 L2-1
L2-2 L3-3 L4-1
L2-3 L4-2 L5-3 L6-1
L4-3 L6-2 L7-3 L8-1
L6-3 L8-2 L9-3 L10-1
L8-3 L10-2 L11-3 L12-1
L10-3 L12-2 L13-3 L14-1
L12-3 L14-2 L15-3 L16 -1
L14-3 L16-2 L17-3 L18-1
L16-3 L18-2 L19-3
L18-3 L20-3
*/