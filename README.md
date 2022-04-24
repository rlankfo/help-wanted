# help-wanted
Want to contribute to open source but don't know where to start? Use this program to search github issues with "help wanted" labels.

```
✗ go build -o hw main.go
✗ ./hw --help
Usage of ./help-wanted:
  -config string
        YAML config file. Command line flags will override values set in this file. (default "~/.config/help-wanted/config.yml")
  -hours int
        Hours since issue was created. (default 72)
  -label value
        Find issues with this label.
  -org value
        Github organization to search.
  -verbose
        Prints Github search string associated with query.
```

Feel free to contribute PRs to make this project better!

## Configuration
By default, the program will load it's config from `~/.config/help-wanted/config.yml`. These values can be overriden with command line args.

```yaml
verbose: true
hours: 72
labels:
- "help wanted"
orgs:
- "ohmyzsh"
- "homebrew"
```
