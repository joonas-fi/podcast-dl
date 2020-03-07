package main

import (
	"fmt"
	"github.com/function61/gokit/atomicfilewrite"
	"github.com/function61/gokit/fileexists"
	"github.com/function61/gokit/jsonfile"
	"github.com/mmcdole/gofeed"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	metaJsonFile = "meta.json"
)

type EpisodeRef struct {
	FeedId    string
	Title     string
	Published time.Time
	Guid      string
}

func (e *EpisodeRef) String() string {
	return fmt.Sprintf(
		"Feed[%s] Guid[%s] Published[%s] Title[%s]",
		e.FeedId,
		e.Guid,
		e.Published.Format(time.RFC3339),
		e.Title)
}

type store struct{}

// TODO: check for guid match in directory name to allow for publish date to change?
func (s *store) Have(ref EpisodeRef) (bool, error) {
	return fileexists.Exists(s.Join(ref, metaJsonFile))
}

func (s *store) Store(ref EpisodeRef, filename string, content io.Reader, meta *gofeed.Item) error {
	if err := os.MkdirAll(s.Join(ref), 0700); err != nil {
		return err
	}

	if err := atomicfilewrite.Write(s.Join(ref, filename), func(sink io.Writer) error {
		_, err := io.Copy(sink, content)
		return err
	}); err != nil {
		return err
	}

	if err := jsonfile.Write(s.Join(ref, metaJsonFile), meta); err != nil {
		return err
	}

	return nil
}

func (s *store) Join(ref EpisodeRef, comps ...string) string {
	comps2 := append([]string{ref.FeedId, ref.Published.Format("2006-01-02") + " " + ref.Guid}, comps...)
	return filepath.Join(comps2...)
}
