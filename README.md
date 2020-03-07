![Build status](https://github.com/joonas-fi/podcast-dl/workflows/Build/badge.svg)
[![Download](https://img.shields.io/github/downloads/joonas-fi/podcast-dl/total.svg?style=for-the-badge)](https://github.com/joonas-fi/podcast-dl/releases)

Utility for hoarding podcasts. Downloads the mp3s and also stores each episode metadata in JSON.


How to run
----------

Dry run tells you what it would do:

```console
$ podcast-dl rss changelog https://changelog.com/podcast/feed
2020/03/07 11:41:13 downloading From open core to open source
```

Adding `--dl` will actually do what the dry run reported it would do:

```console
$ podcast-dl rss --dl changelog https://changelog.com/podcast/feed
(it now downloads them)
```


Downloaded directory structure
------------------------------

```console
$ tree on-the-metal/
on-the-metal/
├── 2019-11-15 01561a241bf5b6ba
│   ├── 8396a7b2.mp3
│   └── meta.json
├── 2019-12-02 178b1338b2c64c3b
│   ├── ee45f70b.mp3
│   └── meta.json
├── 2019-12-06 469e4246e403892a
│   ├── 4a4a1ee6.mp3
│   └── meta.json
├── 2019-12-16 5e8c0cb422cf3477
│   ├── 344925bb.mp3
│   └── meta.json
├── 2019-12-23 17e837a9e8dbb2df
│   ├── f11defc9.mp3
│   └── meta.json
├── 2019-12-30 f587e45a7abf089c
│   ├── 9d0cafea.mp3
│   └── meta.json
├── 2020-01-06 7b752ab00e9c2f6a
│   ├── 9212f8ff.mp3
│   └── meta.json
├── 2020-01-13 e73fa7a55678855b
│   ├── ceda47c6.mp3
│   └── meta.json
├── 2020-01-20 4fdcb004e5ef859c
│   ├── f037142a.mp3
│   └── meta.json
├── 2020-01-27 20cbc9567e70dc4b
│   ├── 754f684f.mp3
│   └── meta.json
└── 2020-02-03 92cbdcae249d7676
    ├── 204d8793.mp3
    └── meta.json

11 directories, 22 files
```

Summary:

- Separate directory tree for each podcast
- One directory for each episode
- Episode directory has the mp3 file and also the metadata in JSON format


Pro-tips for managing many podcasts
-----------------------------------

I have this `common.sh`:

```bash
set -eu

podcast() {
	local feedId="$1"
	local feedUrl="$2"

	podcast-dl rss "$feedId" "$feedUrl" "${@:3}"
}
```

And for each podcast I have a small entrypoint script, e.g. `changelog.sh`:

```bash
source common.sh

# this points to the RSS feed
podcast "changelog" "https://changelog.com/podcast/feed" "$@"
```

Then I can run it as:

```console
$ ./changelog.sh
(dry run results)

$ ./changelog.sh --dl
(download progress)
```
