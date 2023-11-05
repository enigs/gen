# Gen
Quality of life cli tool for generating random stuff I've been using for my personal dev life.

## Installation
1. Build the project
    - Navigate through the project folder and run `go build -o gen` 
2. Add the binary to your path
    - Windows: Set the path to the folder containing the binary in your `PATH` environment variable
    - Linux: Move the binary to `/usr/local/bin` or add the folder containing the binary to your `PATH` environment variable
3. Run `gen` in your terminal to see the available commands


## Current Features
- Generate a random token of any given length using the `-b` flag
    - `gen token -b 32` will generate a random key of length 32 bytes
        -  Token will be converted to un-padded base64 url encoding
        -  Another version will be converted to hex