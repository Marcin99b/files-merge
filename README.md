# files-merge

CLI tool that merges duplicated folders — e.g. `DCIM`, `DCIM(0)`, `DCIM(1)` copied from different cameras or drives — into a single folder, recursively, without touching the originals.

## Install

```
go install github.com/Marcin99b/files-merge@latest
```

Or clone and build:

```
git clone https://github.com/Marcin99b/files-merge
cd files-merge
go build -o files-merge .
```

## Usage

```
files-merge -s <source> -d <destination>
```

- `-s, --source` — folder to scan (default: current directory)
- `-d, --destination` — folder to write the merged result to (required)

## Example

Given:

```
source/
  DCIM/photo1.jpg
  DCIM(0)/photo1.jpg
  DCIM(0)/photo2.jpg
  DCIM(1)/photo3.jpg
```

`files-merge -s source -d destination` produces:

```
destination/
  DCIM/
    photo1.jpg
    photo1(1).jpg
    photo2.jpg
    photo3.jpg
```

Any top-level folder matching `<name>(<n>)` is treated as a duplicate of `<name>`. Subfolders with matching names (e.g. `Pictures/Screenshots` present in more than one duplicate) are merged the same way. Files with clashing names get a `(1)`, `(2)`, ... suffix instead of overwriting each other.
