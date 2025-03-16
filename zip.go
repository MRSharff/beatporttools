package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dhowden/tag"
)

func unzip(source, dest string, noPrompt bool, format string) error {
	formatter := buildFormatter(format)

	zr, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zr.Close()

	newDirs := map[string][]*zip.File{}
	for _, fHandle := range zr.File {
		path := fHandle.Name
		if fHandle.FileInfo().IsDir() {
			slog.Debug("skipping directory", "path", path)
			continue
		}
		file, err := fHandle.Open()
		if err != nil {
			slog.Error("error opening file", "path", path, "error", err)
			continue
		}

		b, err := io.ReadAll(file)
		if err != nil {
			slog.Error("error reading file", "path", path, "error", err)
		}
		file.Close()

		md, err := tag.ReadFrom(bytes.NewReader(b))
		if err != nil {
			slog.Warn("Error reading tag", "path", path, "error", err)
			continue
		}

		if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
			slog.Debug(path)
			for k, v := range md.Raw() {
				slog.Debug(fmt.Sprintf("%s: %s\n", k, v))
			}
		}

		newDir := filepath.Join(dest, formatter(md))
		newDirs[newDir] = append(newDirs[newDir], fHandle)
	}

	for _, dir := range newDirs {
		slices.SortFunc(dir, func(a, b *zip.File) int {
			return strings.Compare(a.Name, b.Name)
		})
	}

	printNewDirs(newDirs)

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
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	fmt.Println("Extracting files...")
	for newDir, files := range newDirs {
		// TODO: do not overwrite files. Ask to overwrite, make a copy, or skip.
		// 	Allow this to be set as a flag so that this can be run without prompts.
		for _, handle := range files {
			newPath := filepath.Join(newDir, handle.Name)
			newFile, err := os.Create(newPath)
			if err != nil {
				slog.Error("Error creating file", "filename", newPath, "error", err)
				continue
			}
			oldFile, err := handle.Open()
			if err != nil {
				slog.Error("Error opening file", "filename", handle.Name, "error", err)
				continue
			}
			if _, err := io.Copy(newFile, oldFile); err != nil {
				slog.Error("error copying file", "filename", handle.Name, "error", err)
			}
		}
	}
	fmt.Println("Done")
	return nil
}

// printNewDirs prints info about where files will be extracted to
// Example Output:
// file1.flac--> dir1\release1 (2025-03-14)\file1.flac
// file2.flac--> dir1\release1 (2025-03-14)\file2.flac
func printNewDirs(newDirs map[string][]*zip.File) {
	// todo: add lowest common parent dir?
	maxFromLength := math.MinInt
	k := len("> ")
	for _, files := range newDirs {
		for _, f := range files {
			if len(f.Name)+k > maxFromLength {
				maxFromLength = len(f.Name) + k
			}
		}
	}

	var sb strings.Builder
	for newDir, files := range newDirs {
		for _, f := range files {
			sb.WriteString(f.Name)
			sb.WriteString(strings.Repeat("-", maxFromLength-len(f.Name)))
			sb.WriteString("> ")
			sb.WriteString(filepath.Join(newDir, f.Name))
			sb.WriteString("\n")
		}
	}

	fmt.Println(sb.String())
}
