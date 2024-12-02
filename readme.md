### Build system
* `make build` builds the project in `./build`
* `make clean` 
* `make test` builds the projects and runs go test 
* `make cov` builds the projects, runs the tests and produces a HTML coverage report in the build directory

### Program Arguments
```
  -depth int
        Maximum depth for directory traversal, used only when recursive is set (default 10)
  -l    Enable logs
  -path string
        Path to the directory containing JSON files
  -r    Enable recursive directory traversal
```

### Assumptions and Design Decisions:
- Process a path to any directory.
	1. Only directories are accepted; invalid paths, file paths, no permissions are reported as errors.
	2. Process one directory
	3. Process files with `.json` extension only (case insensitive, `.JSON` is allowed for example)
	4. Optional recursive directory walking with a max depth option (depth 0 refers to direct children).
	5. Symbolic links are skipped to avoid recursion issues.
	6. Directory path can be provided via arguments or user input if not specified.

- Combine all board lists inside the JSON files into a single JSON output
	- Validity:
		1. JSON files may contain an array of boards or a single board object. Single objects are normalized into an array.
		2. Only JSON objects with `name` and `vendor` are valid (case-sensitive).
		3. `core` and `has_wifi` are optional, and are ommitted if not privided in the original data.
	- JSON properties processing:
		1. Extra key-value properties are allowed and perserved in the final result.
		2. Property keys are trimmed of leading and trailing spaces
		3. `name`, `vendor` & `core` property values are trimmed before evaluation, a value of spaces "   " is considered missing data
	- Duplicate boards
		1. Boards with identical `name` and `vendor` are merged.
		2. If conflicting properties exists, a warning log is produced, and one of the values is choses (Based on read order)

- Order the board list alphabetically first by `vendor`, and then by `name`

- Include metadata in the JSON output under a `_metadata` object including: - The total number of unique vendors - The total number of boards
	- Duplicate data, Merge boards if vendor and board name is identical
		1. "Vendor-A" and "Vendor A" will be considered as different vendors
		2. "Board-1" and "Board 1" will be considered as different boards

- Print JSON
	1. Indented JSON CLI output

- Stretch goal, create a web service which will serve this JSON data over HTTP

- Testing
	1. Cross Platform tests (Github actions)
	2. Unit tests (Go test)

- Logging
	1. Logging can be enabled by passing flag `-l`
	2. Three level of warning exists (INFO, WARN, ERROR)

