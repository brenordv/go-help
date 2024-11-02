# Go Help!
This repo contains a bunch of utilities that I find useful for my day-to-day work. I hope you find them useful too!
Some of them already exists on Linux, but since I'm stuck on Windows, I decided to create my own. :D 

# List of utilities
## guid
```text
---------------------------------------------------
Usage: guid [-c]
  -c    Execute with -c to copy the guid to the clipboard
---------------------------------------------------
Version:  1.0.0
```

## touch
```text
Mimics the functionality of the 'touch' command in Linux.
---------------------------------------------------
Usage: touch [-t] <filename1> <filename2> ... <filenameN>
  -t    Execute with -t to display the current time after touching the file
---------------------------------------------------
version:  1.0.0
```

## cat
```text
Mimics the behavior of the Linux 'cat' command.
---------------------------------------------------
Usage: gocat [-i] [filename]...
  -i Read from stdin if no files are provided
---------------------------------------------------
version: 1.0.0
  -i    Execute with -i to read from stdin if no files are provided
```

## ts
```text
---------------------------------------------------
Usage:
  no parameters          : Prints the current time as a Unix timestamp
  <unix timestamp>       : Converts the Unix timestamp to UTC and local date time
  <YYYY-MM-DD HH:MM:SS>  : Converts the date time to a Unix timestamp
  -h, --help             : Prints this help message
---------------------------------------------------
```

### TODO
1. Add > redirection
2. Add >> redirection
3. Add -n (show line numbers) support