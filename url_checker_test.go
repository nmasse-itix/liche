package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLCheckerCheck(t *testing.T) {
	c := newURLChecker(0, "", nil, newSemaphore(1024), false)

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	assert.NotEqual(t, nil, c.Check("http://www.google.com/README-I-DONT-EXIST.md", "README.md"))

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithExclude(t *testing.T) {
	c := newURLChecker(0, "", regexp.MustCompile(`^http:\/\/localhost:[13]$`), newSemaphore(1024), false)

	for _, u := range []string{"http://localhost:1", "http://localhost:3", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"http://localhost:2", "READYOU.md"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckLocal(t *testing.T) {
	c := newURLChecker(0, "", nil, newSemaphore(1024), true)

	for _, u := range []string{"https://www.google.com"} {
		assert.Equal(t, errSkipped, c.Check(u, "README.md"))
	}

	for _, u := range []string{"README.md", "file://README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}
	for _, u := range []string{"file://foo-bar-missing-file-azertyuiop"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithTimeout(t *testing.T) {
	c := newURLChecker(30*time.Second, "", nil, newSemaphore(1024), false)

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckMany(t *testing.T) {
	c := newURLChecker(0, "", nil, newSemaphore(1024), false)

	for _, us := range [][]string{{}, {"https://google.com", "README.md"}} {
		rc := make(chan urlResult, 1024)
		c.CheckMany(us, "README.md", rc)

		for r := range rc {
			assert.NotEqual(t, "", r.url)
			assert.Equal(t, nil, r.err)
		}
	}
}
func TestURLCheckerResolveURL(t *testing.T) {
	f := newURLChecker(0, "", nil, newSemaphore(1024), false)

	for _, c := range []struct {
		source, target string
		local          bool
	}{
		{"foo", "foo", true},
		{"https://google.com", "https://google.com", false},
	} {
		u, local, err := f.resolveURL(c.source, "foo.md")

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
		assert.Equal(t, c.local, local)
	}
}

func TestURLCheckerResolveURLWithAbsolutePath(t *testing.T) {
	f := newURLChecker(0, "", nil, newSemaphore(1024), false)

	u, _, err := f.resolveURL("/foo", "foo.md")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", u)
}

func TestURLCheckerResolveURLWithDocumentRoot(t *testing.T) {
	f := newURLChecker(0, "foo", nil, newSemaphore(1024), false)

	for _, c := range []struct {
		source, target string
		local          bool
	}{
		{"foo", "foo", true},
		{"https://google.com", "https://google.com", false},
		{"/foo", "foo/foo", true},
	} {
		u, local, err := f.resolveURL(c.source, "foo.md")

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
		assert.Equal(t, c.local, local)
	}
}
