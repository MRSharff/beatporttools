package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
)

const (
	ReleaseYearFormat = "{{release_year}}"
	ReleaseNameFormat = "{{release_name}}"
)

func organizeIntoReleaseFolders(source, dest string, noPrompt bool, format string) {
	var moves []move
	var moveFiles []moveFile

	dirfs := os.DirFS(source)

	dirs, err := os.ReadDir(source)
	if err != nil {
		log.Fatal(err)
	}

	formatter := buildFormatter(format)

	newDirs := make(map[string]struct{})
	for _, dir := range dirs {
		slog.Debug("checking file", "dir_name", dir.Name(), "is_dir", dir.IsDir())
		if dir.IsDir() {
			slog.Debug("file is a directory, skipping")
			continue
		}
		path := dir.Name()
		f, err := dirfs.Open(path)
		if err != nil {
			slog.Warn("Error opening file", "path", path, "error", err)
			continue
		}

		readSeeker, ok := f.(io.ReadSeeker)
		if !ok {
			slog.Warn("File does not implement io.ReadSeeker", "path", path)
			continue
		}

		md, err := tag.ReadFrom(readSeeker)
		if err != nil {
			slog.Warn("Error reading tag", "path", path, "error", err)
			continue
		}

		if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
			slog.Debug(dir.Name())
			for k, v := range md.Raw() {
				slog.Debug(fmt.Sprintf("%s: %s\n", k, v))
			}
		}

		f.Close()

		newDir := filepath.Join(dest, formatter(md))
		newDirs[newDir] = struct{}{}

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

	fmt.Println("Creating new directories...")
	for newDir := range newDirs {
		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			slog.Error("Error creating directory, exiting", "dir", newDir, "error", err)
			return
		}
	}

	fmt.Println("Moving files...")
	for _, move := range moves {
		// TODO: do not overwrite files. Ask to overwrite, make a copy, or skip.
		// 	Allow this to be set as a flag so that this can be run without prompts.
		if err := os.Rename(move.from, move.to); err != nil {
			slog.Warn("Error moving file %s: %s", move.from, err)
		}
	}
	fmt.Println("Done")
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

func buildFormatter(format string) func(md tag.Metadata) string {
	var sb strings.Builder
	var grabbers []func(md tag.Metadata) any

	formatLen := len(format)
	for i := 0; i < formatLen; {
		switch {
		case strings.HasPrefix(format[i:], ReleaseNameFormat):
			sb.WriteString("%s")
			grabbers = append(grabbers, func(md tag.Metadata) any {
				return md.Album()
			})
			i += len(ReleaseNameFormat)
		case strings.HasPrefix(format[i:], ReleaseYearFormat):
			sb.WriteString("%d")
			grabbers = append(grabbers, func(md tag.Metadata) any {
				return md.Year()
			})
			i += len(ReleaseYearFormat)
		default:
			sb.WriteRune(rune(format[i:][0]))
			i++
		}
	}
	format = sb.String()
	return func(md tag.Metadata) string {
		var args []any
		for _, grabber := range grabbers {
			args = append(args, grabber(md))
		}
		return fmt.Sprintf(format, args...)
	}
}

// formatDir is an alternate formatting function, it's a bit more straight
// forward, but it has to run over the whole format string for each tag.
func formatDir(format string, md tag.Metadata) string {
	var sb strings.Builder

	formatLen := len(format)
	for i := 0; i < formatLen; {
		switch {
		case strings.HasPrefix(format[i:], ReleaseNameFormat):
			sb.WriteString(md.Album())
			i += len(ReleaseNameFormat)
		case strings.HasPrefix(format[i:], ReleaseYearFormat):
			sb.WriteString(strconv.Itoa(md.Year()))
			i += len(ReleaseYearFormat)
		default:
			sb.WriteRune(rune(format[i:][0]))
			i++
		}
	}
	return sb.String()
}
