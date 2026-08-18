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
}

type Note struct {
	Title		string `json:"title"`
	Content		string `json:"content"`
	Modified	string `json:"modified"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


func (a *App) ListNotes() ([]string, error)  {       // filenames in the notes dir
	entries, err := os.ReadDir(a.notesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []Note{}, nil
		}
		return nil, fmt.Errorf("reading notes directory: %w", err)
	}

	var notes []Note
	for _, entry := range entries {
		if entry.IsDir() || !string.HasSuffix(entry.Name(), ".md") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		title := strings.TrimSuffix(entry.Name(), ".md")
		title := strings.ReplaceAll(title, "-", " ")

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

func (a *App) notesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "Notes"
	}
	return filepath.Join(home, "Notes")
}


func (a *App) ReadNote(filename string) (Note, error) { // load one file's content
	
}


func (a *App) SaveNote(filename, content string) error // create or overwrite
func (a *App) DeleteNote(filename string) error        // os.Remove