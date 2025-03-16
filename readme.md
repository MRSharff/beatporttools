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

### Example:
```shell
beatporttools organize -source ./testdir -dest ./testdir -y
2025/03/15 22:54:42 WARN Error reading tag path=Shox.Well_St..(Original_Mix).C_Minor.127.2014.aiff error="no tags found"
testdir\Darwin - I_ll Be There (Extended Mix).flac--------------------------> testdir\I'll Be There (2024)\Darwin - I_ll Be There (Extended Mix).flac
testdir\Darwin - I_ll Be There (Radio Mix).flac-----------------------------> testdir\I'll Be There (2024)\Darwin - I_ll Be There (Radio Mix).flac
testdir\Darwin, 3star - Reflections (Original Mix).flac---------------------> testdir\I'll Be There (2024)\Darwin, 3star - Reflections (Original Mix).flac
testdir\Sub_Focus,_Kele.Turn_It_Around.(Original_Mix).B_Major.118.2013.mp3--> testdir\Torus (2013)\Sub_Focus,_Kele.Turn_It_Around.(Original_Mix).B_Major.118.2013.mp3

Creating new directories...
Moving files...
Done

```


# TODO
- [ ] Support zip files (this is what beatport downloads to)
- [ ] Support copying files instead of renaming/moving.
- [x] Support mp3 files
- [ ] Support all Beatport file types. (AIFF and WAV do not come with tags :( )
- [ ] Support different folder naming scheme like -scheme={album_artist}_{album}_{year}  
	Unfortunately, it seems like beatport doesn't tag the release/album artist, so this may be difficult.  
    I actually hate that beatport doesn't tag that because I like being able to sort music by album artist  
    to see all compilations like Bonkers or Happy2bHardcore next to each other. (winamp 4ever)