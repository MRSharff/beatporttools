# Beatport Tools

```shell
beatporttools -h

beatporttools is a tool to organize your music downloaded from beatport. Use at your own risk.
Usage of beatportools:
  -dest string
        destination directory, where you want the release folders to be created (default ".")
  -noprompt
        do not prompt for input, accept all prompts
  -source string
        source directory, where your Beatport downloads are located (default ".")
  -v    show info logs
  -vv
        show debug logs

```

```shell
beatporttools -noprompt

Darwin - I_ll Be There (Extended Mix).flac-------> I'll Be There (2024-01-12)\Darwin - I_ll Be There (Extended Mix).flac
Darwin - I_ll Be There (Radio Mix).flac----------> I'll Be There (2024-01-12)\Darwin - I_ll Be There (Radio Mix).flac
Darwin, 3star - Reflections (Original Mix).flac--> I'll Be There (2024-01-12)\Darwin, 3star - Reflections (Original Mix).flac

Moving files...
Files moved.

```



# TODO
- [ ] Support zip files (this is what beatport downloads to)
- [ ] Support copying files instead of renaming/moving.
- [ ] Support all Beatport file types. (AIFF and WAV do not come with tags :( )
- [ ] Support different folder naming scheme like -scheme={album_artist}_{album}_{year}  
	Unfortunately, it seems like beatport doesn't tag the release/album artist, so this may be difficult.  
    I actually hate that beatport doesn't tag that because I like being able to sort music by album artist  
    to see all compilations like Bonkers or Happy2bHardcore next to each other. (winamp 4ever)