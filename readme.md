# Go Help!
This repo contains a bunch of small apps that I find useful for my day-to-day work. I hope you find them useful too!
Some of them already exists on Linux, but since I'm stuck on Windows, I decided to create my own. :D 

# List of utilities
## guid
Generates a GUID and (optionally) copies it to the clipboard.
```text
---------------------------------------------------
Usage: guid [-c]
  -c    Execute with -c to copy the guid to the clipboard
---------------------------------------------------
Version:  1.0.0
```

## touch
Mimics the functionality of the 'touch' command in Linux.
Good app to have on Windows, if you're used to Linux.
```text
---------------------------------------------------
Usage: touch [-t] <filename1> <filename2> ... <filenameN>
  -t    Execute with -t to display the current time after touching the file
---------------------------------------------------
version:  1.0.0
```

## cat
Mimics the behavior of the Linux 'cat' command.
Probably not really useful if you have the Linux subsystem installed (or if you're actually using linux), but it's a 
fun project to work on.
```text
---------------------------------------------------
Usage: gocat [-i] [filename]...
  -i Read from stdin if no files are provided
---------------------------------------------------
version: 1.0.0
  -i    Execute with -i to read from stdin if no files are provided
```

### TODO
1. Add > redirection
2. Add >> redirection
3. Add -n (show line numbers) support

## ts
A simple utility to convert Unix timestamps to date time and vice versa. It also supports the conversion of date time.
```text
---------------------------------------------------
Usage:
  no parameters          : Prints the current time as a Unix timestamp
  <unix timestamp>       : Converts the Unix timestamp to UTC and local date time
  <YYYY-MM-DD HH:MM:SS>  : Converts the date time to a Unix timestamp
  -h, --help             : Prints this help message
---------------------------------------------------
```

## ncsv
CSV file normalizer. This utility reads a CSV file and tries to normalize the data in a shape-wise manner, making the
data columns match the header columns. This normalization is naive, and just considers position and not the actual data.
Any mismatched data will be saved in a separate file.
You can also split the output into multiple files using the `--split` option, and speed up the process by increasing the
concurrency level using the `--concurrency` option.

Why use this?
Maybe you have a huge CSV file with mismatched columns, and you want to import it to a database (like Postgres). With
this tool you don't have to create a script to prepare the data, this app will do a position normalization of the data 
so you can import and analyze it later.

(I might create a csv 2 json app later, so you can also import the data to a NoSQL database)
```text
---------------------------------------------------
Usage: ncsv [options]
  --header, -e <value>    : Use value as the header (columns separated by commas)
  --file, -f <path>       : Path to the input CSV file (required)
  --split, -s <number>    : Number of lines per output file
  --print-every, -p <num> : Print status every X lines (default: 500)
  --concurrency, -c <num> : Concurrency level (buffer sizes) (default: 1000)
---------------------------------------------------
```
