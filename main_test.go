package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestFindConnection(t *testing.T) {
	os.Args = []string{"", "test_files/test1.txt"}
	read_file(os.Args[1])
	save_data()
	test_flag := true
	tests := []struct {
		room     string
		expected []string
	}{
		{"0", []string{"1", "3"}},
		{"1", []string{"0", "2"}},
		{"2", []string{"1", "3"}},
		{"3", []string{"0", "2"}},
	}

	for _, test := range tests {
		result := find_connection(test.room)
		if !equalSlices(result, test.expected) {
			t.Errorf("find_connection(%s) = %v; want %v", test.room, result, test.expected)
			test_flag = false
		}
	}
	if test_flag {
		fmt.Println("find_connection() Passed")
	}
}

func TestLoopHandler(t *testing.T) {
	tests := []struct {
		room     string
		road     []string
		expected bool
	}{
		{"1", []string{"0", "2"}, true},
		{"2", []string{"0", "2"}, false},
	}
	test_flag := true
	for _, test := range tests {
		result := loop_handler(test.room, test.road)
		if result != test.expected {
			t.Errorf("loop_handler(%s, %v) = %v; want %v", test.room, test.road, result, test.expected)
			test_flag = false
		}
	}
	if test_flag {
		fmt.Println("loop_handler() Passed")
	}
}

func TestDfsPaths(t *testing.T) {
	start_room = "0"
	end_room = "3"
	connect_rows = []string{
		"0-1",
		"0-3",
		"1-2",
		"3-2",
	}

	expected := [][]string{
		{"0", "3"},
		{"0", "1", "2", "3"},
	}
	test_flag := true
	find_road_recursive(start_room, []string{}, &roads)
	if !equal2DSlices(roads, expected) {
		t.Errorf("find_road_recursive(%s) = %v; want %v", start_room, sort_paths_by_length(roads), expected)
		test_flag = false
	}
	if test_flag {
		fmt.Println("find_road_recursive() Passed")
	}
}

func TestFindAllPaths(t *testing.T) {
	start_room = "0"
	end_room = "3"
	connect_rows = []string{
		"0-1",
		"0-3",
		"1-2",
		"3-2",
	}

	find_all_paths()

	expected := [][]string{
		{"0", "3"},
		{"0", "1", "2", "3"},
	}
	roads = delete_same_roads(roads)
	if !equal2DSlices(roads, expected) {
		t.Errorf("find_all_paths() = %v; want %v", roads, expected)
	} else {
		fmt.Println("find_all_paths() Passed")
	}
}

func TestIsOverlapping(t *testing.T) {
	tests := []struct {
		path     []string
		paths    [][]string
		expected bool
	}{
		{
			[]string{"0", "1", "2", "3"},
			[][]string{{"0", "3"}},
			false,
		},
		{
			[]string{"0", "1", "2", "3"},
			[][]string{{"0", "1", "3"}},
			true,
		},
	}
	test_flag := true
	for _, test := range tests {
		result := is_overlapping(test.path, test.paths)
		if result != test.expected {
			t.Errorf("is_overlapping(%v, %v) = %v; want %v", test.path, test.paths, result, test.expected)
			test_flag = false
		}
	}
	if test_flag {
		fmt.Println("is_overlapping() Passed")
	}
}

func TestFindMaxNonOverlappingPaths(t *testing.T) {
	paths := [][]string{
		{"0", "3"},
		{"0", "1", "2", "3"},
		{"0", "1", "3"},
	}

	expected := [][]string{
		{"0", "3"},
		{"0", "1", "2", "3"},
	}

	result := find_max_non_overlapping_paths(paths)
	if !equal2DSlices(result, expected) {
		t.Errorf("find_max_non_overlapping_paths(%v) = %v; want %v", paths, result, expected)
	} else {
		fmt.Println("find_max_non_overlapping_paths() Passed")
	}
}

func TestSortPathsByLength(t *testing.T) {
	paths := [][]string{
		{"0", "1", "2", "3"},
		{"0", "3"},
	}

	expected := [][]string{
		{"0", "3"},
		{"0", "1", "2", "3"},
	}

	sort_paths_by_length(paths)
	if !equal2DSlices(paths, expected) {
		t.Errorf("sort_paths_by_length() = %v; want %v", paths, expected)
	} else {
		fmt.Println("sort_paths_by_length() Passed")

	}
}

func TestDispatchAnts(t *testing.T) {
	find_all_paths()
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	dispatch_ants()

	w.Close()
	os.Stdout = origStdout

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		output.Write(buf[:n])
	}

	outputStr := output.String()
	outputLines := strings.Split(strings.TrimSpace(outputStr), "\n")
	if len(outputLines) > 11 {
		t.Errorf("dispatch_ants() output too long: %s", outputStr)

	} else {
		fmt.Println("dispatch_ants() Passed")
	}
}

func TestQualified(t *testing.T) {
	map_paths := []string{
		"maps/example00.txt",
		"maps/example01.txt",
		"maps/example02.txt",
		"maps/example03.txt",
		"maps/example04.txt",
		"maps/example05.txt"}
	excepted := []int{6, 8, 11, 6, 6, 8}
	for i, path := range map_paths {
		os.Args = []string{"", path}
		read_file(os.Args[1])
		save_data()
		find_all_paths()
		origStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		dispatch_ants()

		w.Close()
		os.Stdout = origStdout

		var output strings.Builder
		buf := make([]byte, 1024)
		for {
			n, _ := r.Read(buf)
			if n == 0 {
				break
			}
			output.Write(buf[:n])
		}
		if step > excepted[i] {
			t.Errorf("%s map number of steps is %d, expected %d", path, step, excepted[i])
		} else {
			fmt.Printf("Qualified Tests Map %s Passed", path[5:14])
			fmt.Println()
		}
		remove_all_variables()
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equal2DSlices(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return strings.Join(a[i], ",") < strings.Join(a[j], ",")
	})
	sort.Slice(b, func(i, j int) bool {
		return strings.Join(b[i], ",") < strings.Join(b[j], ",")
	})
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func remove_all_variables() {
	roads = nil
	ant_count = 0
	start_room = ""
	end_room = ""
	connect_rows = nil
	comment_rows = nil
	rows = nil
	step = 0
}
