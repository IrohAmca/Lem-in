package main

import (
	"fmt"
	"os"
	"strings"
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

func save_data(start, end, comment, action int, sentences []string) ([]string, []string, []string, []string) {
	start_rows := []string{}
	end_rows := []string{}
	comment_rows := []string{}
	connect_rows := []string{}
	for i, row := range sentences {
		if i > start && i < end {
			row = string(row[0])
			start_rows = append(start_rows, row)
		}
		if comment != 0 {
			if i > comment && i < end {
				row = string(row[0])
				comment_rows = append(comment_rows, row)
			}
		}
		if i == end+1 {
			row = string(row[0])
			end_rows = append(end_rows, row)
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

	return start_rows, comment_rows, end_rows, connect_rows
}

func read_file(file_path string) string {
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

	content := string(data[:count])
	defer file.Close()
	return content
}

func find_connection(current_room string, connect_rows []string) []string {
	selected_rooms := []string{}
	result := []string{}
	for _, row := range connect_rows {
		if strings.Contains(row, string(current_room[0])) {
			words := strings.Split(row, " ")
			for _, word := range words {
				if word != current_room {
					selected_rooms = append(selected_rooms, word)
				}
			}
		}
	}
	for _, row := range selected_rooms {
		for _, word := range row {
			if word != '-' && word != ' ' && word != rune(current_room[0]) {
				result = append(result, string(word))
			}
		}
	}
	return result
}

func check_end_room(option_rooms []string, end_rows []string) bool {
	for _, room := range option_rooms {
		for _, end_room := range end_rows {
			if room == end_room {
				return true
			}
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
func find_road_recursive(room string, start_rows, end_rows, connect_rows, comment_rows []string, road []string, roads *[][]string) {
	if check_end_room([]string{room}, end_rows) {
		*roads = append(*roads, append(road, room))
		return
	}
	connect_rooms := find_connection(room, connect_rows)
	for _, next_room := range connect_rooms {
		if loop_handler(next_room, road) {
			find_road_recursive(next_room, start_rows, end_rows, connect_rows, comment_rows, append(road, room), roads)
		}
	}
}

func find_road_options_recursive(start_rows, end_rows, comment_rows, connect_rows []string) [][]string {
	roads := [][]string{}
	for _, room := range start_rows {
		find_road_recursive(room, start_rows, end_rows, connect_rows, comment_rows, []string{}, &roads)
	}
	return roads
}

func main() {
	content := read_file("map_example_1.txt")
	var start_rows, comment_rows, end_rows, connect_rows []string
	start_rows, comment_rows, end_rows, connect_rows = save_data(find_start_end_comment(content))
	/*
		fmt.Println("Start Rows: ")
		for _, row := range start_rows {
			fmt.Println(row)
		}
		fmt.Println("Comment Rows: ")
		for _, row := range comment_rows {
			fmt.Println(row)
		}
		fmt.Println("End Rows: ")
		for _, row := range end_rows {
			fmt.Println(row)
		}
		fmt.Println("Connect Rows: ")
		for _, row := range connect_rows {
			fmt.Println(row)
		}
	*/
	// fmt.Println(find_connection("2", connect_rows))
	fmt.Println(find_road_options_recursive(start_rows, end_rows, comment_rows, connect_rows))
}
