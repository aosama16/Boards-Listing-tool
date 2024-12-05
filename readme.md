# Build system
* `make build` 
	- Builds both the CLI and web applicaiton in `./build`
	- Outputs `cli_boards_merger` & `web_boards_merger`
* `make build-cli` 
	- Builds the CLI applicaiton in `./build`
	- Outputs `cli_boards_merger`
* `make build-web` 
	- Builds the web applicaiton in `./build`
	- Outputs `web_boards_merger`
* `make clean` 
	- Removes the build directory `./build`
* `make test` 
	- Builds the project and runs go test 
* `make cov` 
	- Builds the project, runs the tests, and produces an HTML coverage report.
	- Outputs `coverage.html` in the build directory `./build`

# Running the Application
## cli-boards-merger arguments
`./build/cli_boards_merger -h`
```
  -path   string
          Path to the directory containing JSON files
  -r      Enable recursive directory traversal
  -depth  int
          Maximum depth for directory traversal, used only when recursive is set (default 10)
  -l      Enable logs
```

## web-boards-merger arguments
`./build/web_boards_merger -h`
```
  -port   string
          Port number for the web server (default "8080")
```

# Project Structure
```
 ├── build                  Build directory generated from `make build`, contains executable and coverage report
 ├── cmd
 │   ├── cli                Driver code for CLI application
 │   └── web                Driver code for web application
 └── Internal
     ├── core               Contains logic for directory searching and aggregating JSON files
     ├── model              Data structure for boards and associated logic for Marshaling, Unmarshaling & merging boards
     ├── utils
     |   ├── logger         Simple Logging library, can be enabled/disabled
     |   └── testutils      Utility functions for testing (mainly temp directory and file management)
     └── web                Contains routing logic, and HTML templates handling
         ├── static         Static files to be served in a file server		
         └── templates      HTML templates to be processed by "html/template"
```

# Assumptions and Design Decisions:
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
		1. Duplicate keys are allowed, but only one output key is generated, latest key read will determine the value chosen for that key
		2. Extra key-value properties are allowed and perserved in the final result.
		3. Property keys are trimmed of leading and trailing spaces
		4. `name`, `vendor` & `core` property values are trimmed before evaluation, a value of spaces "   " is considered missing data
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
	1. Use htmx to handle POST request to process a directory
	2. Use "html/template" to handle html rendering
	3. A directory path must available to the local server, and any relative path is relative to the current working directory
	4. Display results in a table format
		- Optional arguments `core` and `has_wifi` display `N/A` if not available
		- Additional properties are displayed as is in the "Additional Info" column
	5. Errors are displayed instead of the table result
	6. Logging is always enabled, and is written to the terminal

- Testing
	1. Cross Platform tests (Github actions)
	2. Table-Driven unit tests (Go test)

- Logging
	1. Logging can be enabled by passing flag `-l`
	2. Three level of warning exists (INFO, WARN, ERROR)