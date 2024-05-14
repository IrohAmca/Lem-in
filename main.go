package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
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

func save_data(start, end, comment, action int, sentences []string) ([]string, []string, []string, []string, []string) {
	start_rows := []string{}
	end_rows := []string{}
	comment_rows := []string{}
	connect_rows := []string{}
	action_rows := []string{}
	for i, row := range sentences {
		if i > start && i < end {
			start_rows = append(start_rows, row)
		}
		if comment != 0 {
			if i > comment && i < end {
				comment_rows = append(comment_rows, row)
			}
		}
		if i == end+1 {
			end_rows = append(end_rows, row)
		}
		if action != 0 {
			if i > end+1 && i < action {
				connect_rows = append(connect_rows, row)
			}
			if i >= action {
				action_rows = append(action_rows, row)
			}
		} else {
			if i > end+1 {
				connect_rows = append(connect_rows, row)
			}
		}
	}
	return start_rows, comment_rows, end_rows, connect_rows, action_rows
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

func find_next_rooms(current_room string, connect_rows, comment_rows []string) []string {
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
	for _, row := range comment_rows {
		if strings.Contains(selected_rooms[0], string(current_room[0])) {
			result = append(result, row)
		}
	}
	return result
}

func select_room(rooms []string, target_room string) string {
	rooms_cor := map[string][2]int{}
	for _, room := range rooms {
		room_info := strings.Split(room, " ")
		room_name := room_info[0]
		room_x, _ := strconv.Atoi(room_info[1])
		room_y, _ := strconv.Atoi(room_info[2])
		rooms_cor[room_name] = [2]int{room_x, room_y}
	}
	target_cor := map[string][2]int{}
	target_info := strings.Split(target_room, " ")
	target_name := target_info[0]
	target_x, _ := strconv.Atoi(target_info[1])
	target_y, _ := strconv.Atoi(target_info[2])
	target_cor[target_name] = [2]int{target_x, target_y}

	min_distance := math.MaxInt32
	selected_room := ""
	for room, cor := range rooms_cor {
		distance := calculate_distance(cor, target_cor[target_name])
		if distance < min_distance {
			min_distance = distance
			selected_room = room
		}
	}
	fmt.Println("Selected Room: ", selected_room)
	return selected_room
}

func calculate_distance(cor1, cor2 [2]int) int {	
	x1, y1 := cor1[0], cor1[1]
	x2, y2 := cor2[0], cor2[1]
	return int(math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
}

func find_road(start_rows, end_rows, connect_rows, comment_rows []string) {
	for {
		current_room := start_rows[0]
		selected_room := select_room(find_next_rooms(current_room, connect_rows, comment_rows), end_rows[0]) // Birden fazla end için seçim ile eklenti yapılmalı
		current_room = selected_room
		if current_room == end_rows[0] {
			fmt.Println("End Room: ", current_room)
			break
		}
	}
}

func main() {
	content := read_file("map_example_1.txt")
	var start_rows, comment_rows, end_rows, connect_rows, actions_rows []string
	start_rows, comment_rows, end_rows, connect_rows, actions_rows = save_data(find_start_end_comment(content))
	/*fmt.Println("Start Rows: ")
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
	fmt.Println("Actions Rows: ")
	for _, row := range actions_rows {
		fmt.Println(row)
	}*/
	find_road(start_rows, end_rows, connect_rows, comment_rows)
	fmt.Print("Actions: ", actions_rows)
}