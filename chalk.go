package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"context"
)

// App struct
type App struct {
	ctx context.Context
	notesDir string
}

type Note struct {
	Filename	string `json:"filename"`
	Title		string `json:"title"`
	Content		string `json:"content"`
	Modified	string `json:"modified"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	home, err := os.UserHomeDir()
	dir := "Notes"
	if err == nil {
		dir = filepath.Join(home, "Notes")
	}
	return &App{notesDir: dir}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


func (a *App) ListNotes() ([]Note, error)  {       // filenames in the notes dir
	entries, err := os.ReadDir(a.notesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Note{}, nil
		}
		return nil, fmt.Errorf("reading notes directory: %w", err)
	}

	var notes []Note
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		title := strings.TrimSuffix(entry.Name(), ".md")
		title = strings.ReplaceAll(title, "-", " ")

		notes = append(notes, Note{
			Filename: entry.Name(),
			Title: title,
			Modified: info.ModTime().Format("2006-01-02 15:04"),
		})
	}
	sort.Slice(notes, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02 15:04", notes[i].Modified)
		tj, _ := time.Parse("2006-01-02 15:04", notes[j].Modified)
		return ti.After(tj)
	})

	return notes, nil
}

/*func (a *App) notesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "Notes"
	}
	return filepath.Join(home, "Notes")
}*/


func (a *App) ReadNote(filename string) (Note, error) { // load one file's content
	// 1. sanitize filename with filepath.Clean / reject path separators
	safe := filepath.Base(filepath.Clean(filename))
	if safe != filename {
		return Note{}, fmt.Errorf("invalid filename: %s", filename)
	}

	// 2. filepath.Join(a.notesDir, filename)
	path := filepath.Join(a.notesDir, safe)
	// 3. os.ReadFile(path)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Note{}, fmt.Errorf("note not found: %s", filename)
		}
		return Note{}, fmt.Errorf("reading note %s: %w", filename, err)
	}
	// 5. return string(contents), nil
	return Note{Content: string(content)}, nil
}


func (a *App) SaveNote(filename, content string) error { // create or overwrite
	// 1. sanitize filename (same as ReadNote)
	safe := filepath.Base(filepath.Clean(filename))
	if safe != filename {
		return fmt.Errorf("invalid filename: %s", filename)
	}
	// 2. filepath.Join(a.notesDir, filename)
	path := filepath.Join(a.notesDir, safe)
	// 3. os.WriteFile(path, []byte(content), 0644)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("saving note %s: %w", filename, err)
	}
	return nil
}
func (a *App) DeleteNote(filename string) error {       // os.Remove
	// 1. sanitize filename
	safe := filepath.Base(filepath.Clean(filename))
	if safe != filename {
		return fmt.Errorf("invalid filename: %s", filename)
	}

	// 2. filepath.Join(a.notesDir, filename)
	path := filepath.Join(a.notesDir, safe)
	// 3. os.Remove(path)
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("deleting note %s: %w", filename, err)
	}
	return nil
}