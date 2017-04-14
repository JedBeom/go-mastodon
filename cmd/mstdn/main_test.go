package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadFileFile(t *testing.T) {
	b, err := readFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatalf("should read something: %v", err)
	}
}

func TestReadFileStdin(t *testing.T) {
	f, err := os.Open("main.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	stdin := os.Stdin
	os.Stdin = f
	defer func() {
		os.Stdin = stdin
	}()

	b, err := readFile("-")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatalf("should read something: %v", err)
	}
}

func TestTextContent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "<p>foo</p>", want: "foo"},
		{input: "<p>foo<span>\nbar\n</span>baz</p>", want: "foobarbaz"},
		{input: "<p>foo<span>\nbar<br></span>baz</p>", want: "foobar\nbaz"},
	}
	for _, test := range tests {
		got := textContent(test.input)
		if got != test.want {
			t.Fatalf("want %q but %q", test.want, got)
		}
	}
}

func TestGetConfig(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "mstdn")
	if err != nil {
		t.Fatal(err)
	}
	home := os.Getenv("HOME")
	appdata := os.Getenv("APPDATA")
	os.Setenv("HOME", tmpdir)
	os.Setenv("APPDATA", tmpdir)
	defer func() {
		os.RemoveAll(tmpdir)
		os.Setenv("HOME", home)
		os.Setenv("APPDATA", appdata)
	}()

	file, config, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(file); err == nil {
		t.Fatal("should not exists")
	}
	if config.AccessToken != "" {
		t.Fatalf("should be empty: %v", config.AccessToken)
	}
	if config.ClientID == "" {
		t.Fatalf("should not be empty")
	}
	if config.ClientSecret == "" {
		t.Fatalf("should not be empty")
	}
	config.AccessToken = "foo"
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(file, b, 0700)
	if err != nil {
		log.Fatal(err)
	}
	file, config, err = getConfig()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("should exists: %v", err)
	}
	if got := config.AccessToken; got != "foo" {
		t.Fatalf("want %q but %q", "foo", got)
	}
}
