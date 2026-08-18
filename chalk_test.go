package main

import (
	"testing"
	"path/filepath"
	"os"
)


func TestListEmptyNotes(t *testing.T) {
    app := &App{notesDir: t.TempDir()}
	notes, err := app.ListNotes()
	if err != nil {
		t.Fatalf("expected no error on empty dir, got: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("expected 0 notes, got %d", len(notes))
	}
}

func TestListingMultipleNotes(t *testing.T) {
    dir := t.TempDir()
	app := &App{notesDir: dir}
	name1 := "note1.md"
	name2 := "note2.md"
	path1 := filepath.Join(dir, name1)
	path2 := filepath.Join(dir, name2)
	content1 := "This is note 1"
	content2 := "This is note 2"
	err := os.WriteFile(path1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("writing to file failed")
	}
	err = os.WriteFile(path2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("writing to file failed")
	}

	notes, err := app.ListNotes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
}

