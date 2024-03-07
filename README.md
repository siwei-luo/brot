# 🍞 Brot

![Image of Brot](https://github.com/siwei-luo/brot/blob/7c26aac6a657025028d4340660ff090a7807b15d/brot.png)

Brot is a simple configurable program which does repetitive tasks I do not want to do by myself anymore.

In case you are wondering what the heck this is all about and about its limited usefulness and special use cases at
the moment because I am mainly using this project to learn Golang from scratch.

## Configuration

Brot will look for a configuration file named _brot.yaml_ at following locations in the following order (the last in
list wins prioritization):
* _/etc/brot_
* _$HOME/.config_
* In the same directory where the binary is located at.
* Passed as flag _--config_ or _-c_.

```yaml
---
apiVersion: v1
defaults:
  loglevel: debug
  logformat: text
```

__apiVersion:__ Brot uses [semantic versioning](https://semver.org/) and this value describes to which major version
the configuration file is compatible to.

__loglevel:__ Possible values are _debug_, _info_, _warn_ and _error_.

__logformat:__ Possible values are _text_ or _json_.

## Global flags

__--config, -c__ Path to configuration file to use.

__--verbosity, -v__ Integer value for verbosity, the higher the value the more verbose.
1 ~ error,
2 ~ warn,
3 ~ info,
4 ~ debug

## Sub command: relocate

Use this sub command to move or copy files around using rules. E.g. to tidy up your download directory.

Brot will not change anything in case there is already a file in the destination directory with the same name.

### Configuration: relocate

All relocation rules have to specified as YAML list with the key _relocate:_ in the root. Each relocation item needs
following values to execute properly.

```yaml
relocate:
  - name: move pdfs
    src: $HOME/Downloads
    dst: $HOME/Documents
    patterns:
      - "*.pdf"
    mode: move
  - name: copy pictures
    src: $HOME/Downloads
    dst: /media/USB/Pictures
    patterns:
      - "*.jpg"
      - "*.png"
    mode: copy
```

__name:__ Human readable alias for each rule. It must not be unique, but it helps if it actually is.

__src:__ Directory to read files from. Does not follow any symbolic links if found.

__dst:__ Directory to relocate files to. Brot will not create the destination directory for you if it does not exist.

Environment variable expansion is possible for values in _src:_ and _dst:_. E.g. _$HOME_ expands to the user's home
directory if set.

__patterns:__ Specify _patterns:_ to target only specific files in _src:_. Leave it empty to match all files.

__mode__: Specify either _move_ or _copy_.

### Flags: relocate

__--dry-run, -d__ Just print out possible matches but do not move/copy anything.

## Sub command: cleanup

Use this sub command to remove files around using rules. E.g. to tidy up your download directory.

### Configuration: cleanup

All cleanup rules have to specified as YAML list with the key _cleanup:_ in the root. Each cleanup item needs
following values to execute properly.

```yaml
cleanup:
  - name: mac os foo
    src: $HOME/Downloads
    patterns:
      - ".DS_Store"
      - "._.DS_Store"
```

__name:__ Human readable alias for each rule. It must not be unique, but it helps if it actually is.

__src:__ Directory to read files from. Does not follow any symbolic links if found.

Environment variable expansion is possible for values in _src:_. E.g. _$HOME_ expands to the user's home
directory if set.

__patterns:__ Specify _patterns:_ to target only specific files in _src:_. Leave it empty to match all files.

### Flags: cleanup

__--dry-run, -d__ Just print out possible matches but do not remove anything.
 configuration
## Sub command: completion

Use this sub command to generate shell completions for Bash, Fish, PowerShell or Zsh which can be sourced.

```sh
brot completion [bash|fish|powershell|zsh]
```

## Build from source

```sh
# get the source
$ go get github.com/siwei-luo/brot

# (optional) update dependencies
$ go get -u
$ go mod tidy

# build for your current environment
# set GOOS and GOARCH to cross compile for other targets
$ go build github.com/siwei-luo/brot
```

## Logo

The logo is derived from the original version created by [Takuya Ueda](https://twitter.com/tenntenn) licensed under
[CC 3.0 Attributions](https://creativecommons.org/licenses/by/3.0/de/deed.en) license.

## License

Copyright © 2021-2024 Siwei Luo <siwei@lu0.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
