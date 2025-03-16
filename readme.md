# Beatport Tools

```shell
beatporttools -h

A tool for working with music files downloaded from Beatport
Usage:
        beatporttools <command> [arguments]

Global Flags:
  -v    show info logs
  -vv
        show debug logs
Commands:
        organize        Reorganizes music downloaded from beatport

```

# Commands

## Organize

```shell
beatporttools organize -h
usage:
        beatporttools organize [-source source] [-dest dest] [-y]
flags:
  -dest string
        destination directory, where you want the release folders to be created (default ".")
  -source string
        source directory, where your Beatport downloads are located (default ".")
  -y    do not prompt for input, accept all prompts
example:
        beatporttools organize -y -source ~/Downloads/beatport_tracks_2025_03 -dest ~/Downloads/beatport_tracks_2025_03_organized
```



# TODO
- [ ] Support zip files (this is what beatport downloads to)
- [ ] Support copying files instead of renaming/moving.
- [ ] Support all Beatport file types. (AIFF and WAV do not come with tags :( )
- [ ] Support different folder naming scheme like -scheme={album_artist}_{album}_{year}  
	Unfortunately, it seems like beatport doesn't tag the release/album artist, so this may be difficult.  
    I actually hate that beatport doesn't tag that because I like being able to sort music by album artist  
    to see all compilations like Bonkers or Happy2bHardcore next to each other. (winamp 4ever)