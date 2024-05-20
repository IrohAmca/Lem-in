package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func separate_rows() {
	var temp string
	for _, row := range strings.Split(content, "\n") {
		if row != "" {
			temp = strings.TrimSpace(row)
			rows = append(rows, strings.Split(temp, " "))
		}
	}
}

func save_data() {
	separate_rows()
	start_flag := false
	end_flag := false
	connect_flag := false
	comment_flag := false
	ant_flag := false
	for i, row := range rows {
		if len(row) == 1 && row[0] != "##start" && row[0] != "##end" {
			count, err := strconv.Atoi(row[0])
			if err == nil {
				ant_count = count
				if ant_count >= 1 {
					ant_flag = true
				}
			}
		}
		if row[0] == "##start" {
			if len(rows[i+1]) == 3 {
				start_room = rows[i+1][0]
				start_flag = true
			}
		}
		if row[0] == "##end" {
			if len(rows[i+1]) == 3 {
				end_room = rows[i+1][0]
				end_flag = true
			}
		}
		if len(row) == 3 && row[0] != start_room && row[0] != end_room {
			comment_rows = append(comment_rows, row[0])
			comment_flag = true

		}
		if len(row) == 1 && strings.Contains(row[0], "-") {
			connect_rows = append(connect_rows, row[0])
			connect_flag = true
		}
	}
	Contains_Error(start_flag, end_flag, connect_flag, comment_flag, ant_flag)
	True_Format_Error()
	Connection_Error()
}

func read_file(file_path string) {
	data, err := os.ReadFile(file_path)
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
func find_road_recursive(room string, road []string, roads *[][]string) {
	if room == end_room {
		*roads = append(*roads, append(road, room))
		return
	}
	connect_rooms := find_connection(room)
	for _, next_room := range connect_rooms {
		if loop_handler(next_room, road) {
			find_road_recursive(next_room, append(road, room), roads)
		}
	}
}
func is_same_road(road []string, roads [][]string) bool {
	for _, r := range roads {
		if len(r) == len(road) {
			same := true
			for i := 0; i < len(road); i++ {
				if road[i] != r[i] {
					same = false
					break
				}
			}
			if same {
				return true
			}
		}
	}
	return false
}
func delete_same_roads(roads [][]string) [][]string {
	unique_roads := [][]string{}
	for _, road := range roads {
		if !is_same_road(road, unique_roads) {
			unique_roads = append(unique_roads, road)
		}
	}
	return unique_roads
}
func find_all_paths() {
	find_road_recursive(start_room, []string{}, &roads)
	if roads == nil {
		fmt.Println("ERROR: No path found")
		os.Exit(1)
	}
	roads = find_max_non_overlapping_paths(sort_paths_by_length(roads))
	roads = delete_same_roads(roads)
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

func sort_paths_by_length(paths [][]string) [][]string {
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if len(paths[i]) > len(paths[j]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}
	return paths
}

func dispatch_ants() {
	find_all_paths()

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
