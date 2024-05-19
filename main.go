package main

import (
	"fmt"
	"os"
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


func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <map_file>")
		os.Exit(1)
	}
	read_file(os.Args[1])
	save_data()
	dispatch_ants()
	fmt.Println("Ant Count", ant_count)
	fmt.Println("Start Room", start_room)
	fmt.Println("Comment Rooms", comment_rows)
	fmt.Println("Connect Room", connect_rows)
	fmt.Println("End Room", end_room)
	fmt.Println("Step:", step)
	fmt.Println("Paths:")
	for _, road := range roads {
		fmt.Println(road)
	}
}

