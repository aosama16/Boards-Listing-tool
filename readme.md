### Functional Requirements:
- Process a path to any directory.
	- **Validate path:** isDir (files are rejected), Empty Path. Incorrect path, permissions
	- Handle only one directory
	- Handle files with .json only, case insensitive (.JSON is allowed for example)
	- Recursive search through subdirectories (optional recursive flag with max depth)
	- Skip Symbolic links
	- depth 0 traversal refers to the children of the directory that will be parsed
	- depth 1 traversal refers to first recursive level of the subdirectories' children
	- processing a path can be done through program arguments and through user input if the arguments aren't used
- Combine all board lists inside the JSON files into a single JSON output
	- LOG: Empty JSON, Empty Array
	- LOG: Invalid JSON
	- Normalize single object into array: Different Structure Array vs single board
	- Error on Vital data, warn on empty fields: Most Importantly, Board Name, Vendor
	- Additional Data, add to additional data field as notes
	- Sanitize the data: trim strings, uniform casing (Start of each word is capitalized)
	- Duplicate data, Merge boards if vendor and board name is identical
		- "Vendor-A" and "Vendor A" will be considered as different vendors
		- "Board-1" and "Board 1" will be considered as different boards
		- merge data only if not existent, if extra duplication, only one will show up
	- perserve extra key-value property in JSON data
	- the only required keys for each board is "name" and "vendor", these are case-sensitive, so "NAME" for example will be treated as an extra property if the "name" property exists as well, if it doesn't, then this board is invalid
	- any fields other that boards[] are ignored
	- "core" and "wifi" are optional data, and are ommitted if not privided in the original data
	- Property keys are trimmed of leading and trailing spaces
	- "name", "vendor" & "core" property values are trimmed before evaluation, a value of spaces "   " is considered missing data
- Order the board list alphabetically first by `vendor`, and then by `name`
- Include metadata in the JSON output under a `_metadata` object including: - The total number of unique vendors - The total number of boards
	- Duplicate data, Merge boards if vendor and board name is identical
		- "Vendor-A" and "Vendor A" will be considered as different vendors
		- "Board-1" and "Board 1" will be considered as different boards
- Print JSON
	- Design: Can be converted into API
	- Config: Formatted JSON Output vs compact
- Stretch goal, create a web service which will serve this JSON data over HTTP
	- path submission in HTTP request: How to take input?
	- Filter Output JSON
- Testing
	- Build System (Cross Platform)
	- Types of test
	- Structure / Testing Framework
- Logging/Warning/Error System
- Potential extensions to this CLI.


### MakeFile commands
- build
- clean
- test
- cov
