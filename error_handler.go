package main

import (
	"fmt"
	"os"
	"strings"
)

func Contains_Error(start_flag, end_flag, connect_flag, comment_flag, ant_flag bool) {
	if !start_flag {
		fmt.Println("No start room")
		os.Exit(1)
	}
	if !end_flag {
		fmt.Println("No end room")
	}
	if !connect_flag {
		fmt.Println("No connection")
		os.Exit(1)
	}
	if !comment_flag {
		fmt.Println("No comment room")
		os.Exit(1)
	}
	if !ant_flag {
		fmt.Println("No ant count")
		os.Exit(1)
	}
}

func True_Format_Error() {
	all_rooms := append(comment_rows, start_room, end_room)

	for _, rows := range all_rooms {
		if strings.HasPrefix(rows, "L")  || strings.HasPrefix(rows, "#") {
			fmt.Println("Room name cannot start with L or #")
			os.Exit(1)
		}
	}
}
