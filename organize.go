package main

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

func organizeIntoReleaseFolders(source, dest string, noPrompt bool) {
	var moves []move
	var moveFiles []moveFile

	dirfs := os.DirFS(source)

	dirs, err := os.ReadDir(source)
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		slog.Debug("checking file", "dir_name", dir.Name(), "is_dir", dir.IsDir())
		path := dir.Name()
		f, err := dirfs.Open(path)
		if err != nil {
			slog.Warn("Error opening file", "path", path, "error", err)
			continue
		}

		metadata, err := flac.ParseMetadata(f)
		if err != nil {
			slog.Warn("Error parsing metadata", "path", path, "error", err)
			continue
		}
		m := map[string]string{}
		for i, md := range metadata.Meta {
			if md.Type == flac.VorbisComment {
				cmt, err := flacvorbis.ParseFromMetaDataBlock(*md)
				if err != nil {
					slog.Warn("Error parsing vorbis comment", "path", path, "comment", cmt, "index", i, "error", err)
					continue
				}

				for _, tag := range cmt.Comments {
					parts := strings.SplitN(tag, "=", 2)
					if len(parts) != 2 {
						slog.Warn("tag did not have 2 parts", "tag", tag)
						continue
					}
					m[parts[0]] = parts[1]
				}
				break
			}
		}

		slog.Debug(dir.Name())
		for k, v := range m {
			slog.Debug(fmt.Sprintf("%s: %s\n", k, v))
		}

		f.Close()

		releaseName := m["album"]
		releaseTime := m["release_time"]
		newDir := filepath.Join(dest, fmt.Sprintf("%s - (%s)", releaseName, releaseTime))

		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			slog.Error("Error creating directory", "dir", newDir, "error", err)
			continue
		}

		oldPath := filepath.Join(source, path)
		newPath := filepath.Join(newDir, path)
		moves = append(moves, move{from: oldPath, to: newPath})
		moveFiles = append(moveFiles, moveFile{path: path, fromDir: oldPath, toDir: newPath})
	}

	slices.SortFunc(moves, func(a, b move) int {
		return strings.Compare(a.from, b.from)
	})

	printMoves(moves)
	//printMovesFiles(moveFiles)

	for !noPrompt {
		switch prompt("move files? y/N") {
		case "y":
			noPrompt = true
		case "N":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Unknown response, please enter 'y' or 'N'")
		}
	}

	fmt.Println("Moving files...")
	for _, move := range moves {
		if err := os.Rename(move.from, move.to); err != nil {
			slog.Warn("Error moving file %s: %s", move.from, err)
		}
	}
	fmt.Println("Files moved.")
}

type move struct {
	from, to string
}

type moveFile struct {
	path           string
	fromDir, toDir string
}

func prompt(s string) string {
	fmt.Println(s)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

// printMoves prints info about where files will move
// Example Output:
// dir1\file1.flac-------> dir1\release1 (2025-03-14)\file1.flac
// dir1\file2.flac----------> dir1\release1 (2025-03-14)\file2.flac
func printMoves(moves []move) {
	// todo: add lowest common parent dir?
	maxFromLength := math.MinInt
	k := len("> ")
	for _, m := range moves {
		if len(m.from)+k > maxFromLength {
			maxFromLength = len(m.from) + k
		}
	}
	var sb strings.Builder
	for _, m := range moves {
		sb.WriteString(m.from)
		sb.WriteString(strings.Repeat("-", maxFromLength-len(m.from)))
		sb.WriteString("> ")
		sb.WriteString(m.to)
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}

// printMoveFiles prints info about where files will move
// Example output:
// file1.flac dir1 -> dir1\release1 (2025-03-14)
// file2.flac dir -> dir1\release1 (2025-03-14)
func printMovesFiles(moveFiles []moveFile) {
	var sb strings.Builder
	for _, m := range moveFiles {
		sb.WriteString(m.path)
		sb.WriteString(":")
		sb.WriteString(m.fromDir)
		sb.WriteString(" --> ")
		sb.WriteString(m.toDir)
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}
