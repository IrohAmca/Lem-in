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
		if len(row)== 1{
			count ,er := strconv.Atoi(row[0])
			if er == nil {
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
	unique_paths := [][]string{}

	for _, path := range all_paths {
		if is_unique_path(path, unique_paths) {
			unique_paths = append(unique_paths, path)
		}
	}

	roads = unique_paths
}

func is_unique_path(path []string, paths [][]string) bool {
	for _, p := range paths {
		for i := 1; i < len(path)-1; i++ { // Start ve end odalarını hariç tut
			for j := 1; j < len(p)-1; j++ { // Start ve end odalarını hariç tut
				if path[i] == p[j] {
					return false
				}
			}
		}
	}
	return true
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
	for i := 0; i < ant_count; i++ {
		ant_paths[i] = roads[i%len(roads)]
	}

	for step := 0; ; step++ {
		moved := false
		var moves []string
		used_rooms := make(map[string]bool)
		for i := 0; i < ant_count; i++ {
			if step < len(ant_paths[i]) {
				room := ant_paths[i][step]
				if room == start_room || room == end_room || !used_rooms[room] {
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, room))
					used_rooms[room] = true
					moved = true
				}
			}
		}
		if moved {
			fmt.Println(strings.Join(moves, " "))
		} else {
			break
		} 
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
	fmt.Println("Paths:")
	for _, road := range roads {
		fmt.Println(road)
	}
}
