package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, c := range []struct {
		argv []string
		args arguments
	}{
		{
			argv: []string{"file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, false, false, false},
		},
		{
			argv: []string{"-c", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, nil, false, false, false},
		},
		{
			argv: []string{"--concurrency", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, nil, false, false, false},
		},
		{
			argv: []string{"-d", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, nil, false, false, false},
		},
		{
			argv: []string{"--document-root", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, nil, false, false, false},
		},
		{
			argv: []string{"-r", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, true, false, false},
		},
		{
			argv: []string{"--recursive", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, true, false, false},
		},
		{
			argv: []string{"-t", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, nil, false, false, false},
		},
		{
			argv: []string{"--timeout", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, nil, false, false, false},
		},
		{
			argv: []string{"-x", "^.*$", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, regexp.MustCompile(`^.*$`), false, false, false},
		},
		{
			argv: []string{"--exclude", "^.*$", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, regexp.MustCompile(`^.*$`), false, false, false},
		},
		{
			argv: []string{"-v", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, false, true, false},
		},
		{
			argv: []string{"--verbose", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, false, true, false},
		},
		{
			argv: []string{"--local-only", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, false, false, true},
		},
		{
			argv: []string{"-l", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, nil, false, false, true},
		},
	} {
		args, err := getArguments(c.argv)

		assert.Equal(t, nil, err)
		assert.Equal(t, args, c.args)
	}
}

func TestGetArgumentsWithInvalidArgv(t *testing.T) {
	for _, argv := range [][]string{
		{"-c", "3.14", "file"},
		{"-t", "foo", "file"},
		{"-c", "-t", "file"},
	} {
		_, err := getArguments(argv)
		assert.NotEqual(t, nil, err)
	}
}
