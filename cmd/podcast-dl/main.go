package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/function61/gokit/dynversion"
	"github.com/function61/gokit/ezhttp"
	"github.com/function61/gokit/logex"
	"github.com/function61/gokit/ossignal"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"os"
	"path"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     os.Args[0],
		Short:   "Podcast downloader",
		Version: dynversion.Version,
	}

	dl := false

	cmd := &cobra.Command{
		Use:   "rss [feedId] [url]",
		Short: "Access RSS feed",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			logger := logex.StandardLogger()

			exitIfError(rssDownload(
				ossignal.InterruptOrTerminateBackgroundCtx(logger),
				args[0],
				args[1],
				dl,
				logger))
		},
	}

	cmd.Flags().BoolVarP(&dl, "dl", "", dl, "Download podcasts. Without this it's a dry-run")
	rootCmd.AddCommand(cmd)

	exitIfError(rootCmd.Execute())
}

func rssDownload(
	ctx context.Context,
	feedId string,
	feedUrl string,
	dl bool,
	logger *log.Logger,
) error {
	logl := logex.Levels(logger)

	res, err := ezhttp.Get(ctx, feedUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fp := gofeed.NewParser()
	feed, err := fp.Parse(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(feed.Title)

	for _, item := range feed.Items {
		if err := handleOneFeedItem(ctx, feedId, item, dl, logl); err != nil {
			if err == context.Canceled { // user requested cancel
				return nil
			}

			logl.Error.Printf("re-trying once after %v", err)

			if err := handleOneFeedItem(ctx, feedId, item, dl, logl); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleOneFeedItem(
	ctx context.Context,
	feedId string,
	item *gofeed.Item,
	dl bool,
	logl *logex.Leveled,
) error {
	st := &store{}

	guidSha1 := sha1.Sum([]byte(item.GUID))
	guidSha1Hex := hex.EncodeToString(guidSha1[0:8])

	ref := EpisodeRef{
		FeedId:    feedId,
		Title:     item.Title,
		Guid:      guidSha1Hex,
		Published: *item.PublishedParsed,
	}

	have, err := st.Have(ref)
	if err != nil {
		return err
	}

	if have {
		logl.Debug.Printf("already got %s", item.Title)
		return nil
	} else {
		logl.Info.Printf("downloading %s", item.Title)
	}

	if len(item.Enclosures) != 1 {
		logl.Error.Println("skipping item that does not contain exactly one enclosure")
		return nil
	}

	if !dl { // dry-run
		return nil
	}

	enclosure := item.Enclosures[0]

	// to remove any query parameters
	enclosureUrl, err := url.Parse(enclosure.URL)
	if err != nil {
		return err
	}

	req, err := ezhttp.Get(ctx, enclosure.URL)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	if err := st.Store(ref, path.Base(enclosureUrl.Path), req.Body, item); err != nil {
		return err
	}

	return nil
}

func exitIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
