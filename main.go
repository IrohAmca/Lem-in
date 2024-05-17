package main

import (
	"fmt"
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
)

func seperate_rows(data string) []string {
	return strings.Split(data, "\n")
}

func find_start_end_comment(data string) (int, int, int, int, []string) {
	result_start := 0
	result_end := 0
	result_comment := 0
	result_actions := 0
	actions_flag := false
	sentences := seperate_rows(data)
	for i, row := range sentences {
		if strings.HasPrefix(row, "##s") {
			result_start = i
		}
		if strings.HasPrefix(row, "##e") {
			result_end = i
		}
		if strings.HasPrefix(row, "##c") {
			result_comment = i
		}
		if !actions_flag {
			if strings.HasPrefix(row, "L") {
				result_actions = i
				actions_flag = true
			}
		}
	}
	return result_start, result_end, result_comment, result_actions, sentences
}

func save_data(start, end, comment, action int, sentences []string) {
	for i, row := range sentences {
		if i == start-1 {
			ant_count, _ = strconv.Atoi(row) // Tüm satırı al
		}
		if i == start+1 {
			words := strings.Split(row, " ")
			start_room = words[0]
		}
		if i > start+1 && i < end {
			words := strings.Split(row, " ")
			comment_rows = append(comment_rows, words[0])
		}
		if i == end+1 {
			words := strings.Split(row, " ")
			end_room = words[0]
		}
		if action != 0 {
			if i > end+1 && i < action {
				connect_rows = append(connect_rows, row)
			}
		} else {
			if i > end+1 {
				connect_rows = append(connect_rows, row)
			}
		}
	}
}

func read_file(file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data := make([]byte, 1024)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	content = string(data[:count])
	defer file.Close()
}

func find_connection(current_room string) []string {
	selected_rooms := []string{}
	for _, row := range connect_rows {
		if strings.Contains(row, current_room) {
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

func check_end_room(option_rooms []string) bool {
	for _, room := range option_rooms {
		if room == end_room {
			return true
		}
	}
	return false
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
	roads = append(roads, bfs_paths(start_room)...)
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

	ant_positions := make(map[string][]int)
	ant_paths := make([][]string, ant_count)
	for i := 0; i < ant_count; i++ {
		ant_paths[i] = roads[i%len(roads)]
	}

	for step := 1; ; step++ {
		moved := false
		var moves []string
		for i := 0; i < ant_count; i++ {
			if step < len(ant_paths[i]) {
				room := ant_paths[i][step]
				if len(ant_positions[room]) == 0 || !contains(ant_positions[room], i+1) {
					ant_positions[room] = append(ant_positions[room], i+1)
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, room))
					moved = true
				}
			}
		}
		if moved {
			fmt.Println(strings.Join(moves, " "))
		} else {
			break
		}
		clear_ant_positions(ant_positions, step, ant_paths)
	}
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func clear_ant_positions(ant_positions map[string][]int, step int, ant_paths [][]string) {
	for room := range ant_positions {
		new_positions := []int{}
		for _, ant := range ant_positions[room] {
			if step+1 < len(ant_paths[ant-1]) && ant_paths[ant-1][step+1] == room {
				new_positions = append(new_positions, ant)
			}
		}
		if len(new_positions) == 0 {
			delete(ant_positions, room)
		} else {
			ant_positions[room] = new_positions
		}
	}
}

func main() {
	read_file(os.Args[1])
	start, end, comment, action, sentences := find_start_end_comment(content)
	save_data(start, end, comment, action, sentences)

	dispatch_ants()

	fmt.Println("Start Room", start_room)
	fmt.Println("Comment Rooms", comment_rows)
	fmt.Println("Connect Room", connect_rows)
	fmt.Println("End Room", end_room)
	fmt.Println("Paths:")
	for _, road := range roads {
		fmt.Println(road)
	}
}
