package main

import (
	"os"
	"strings"
	"testing"
)

func Test_ReadFile(t *testing.T) {

}
func TestFindConnection(t *testing.T) {
	os.Args = []string{"", "test_files/test1.txt"}
	read_file(os.Args[1])
	save_data()

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
		}
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

	for _, test := range tests {
		result := loop_handler(test.room, test.road)
		if result != test.expected {
			t.Errorf("loop_handler(%s, %v) = %v; want %v", test.room, test.road, result, test.expected)
		}
	}
}

func TestBfsPaths(t *testing.T) {
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

	result := bfs_paths(start_room)
	if !equal2DSlices(result, expected) {
		t.Errorf("bfs_paths(%s) = %v; want %v", start_room, result, expected)
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

	if !equal2DSlices(roads, expected) {
		t.Errorf("find_all_paths() = %v; want %v", roads, expected)
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

	for _, test := range tests {
		result := is_overlapping(test.path, test.paths)
		if result != test.expected {
			t.Errorf("is_overlapping(%v, %v) = %v; want %v", test.path, test.paths, result, test.expected)
		}
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
	}
}

func TestDispatchAnts(t *testing.T) {
	find_all_paths()
	sort_paths_by_length(roads)

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
	if len(outputLines)>11 {
		t.Errorf("dispatch_ants() output too long: %s", outputStr)
	
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
	for i := range a {
		if !equalSlices(a[i], b[i]) {
			return false
		}
	}
	return true
}
