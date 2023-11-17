# Gen
Quality of life cli tool for generating random stuff I've been using for my personal dev life.

## Installation
1. Build the project
    - Navigate through the project folder and run `go build -o gen` 
2. Add the binary to your path
    - Windows: Set the path to the folder containing the binary in your `PATH` environment variable
    - Linux: Move the binary to `/usr/local/bin` or add the folder containing the binary to your `PATH` environment variable
3. Run `gen` in your terminal to see the available commands

## Notes:
1. Avoid building and running `gen framework` within its own directory for you to avoid deleting `framework.20231117.zip` file.
2. If you accidentally deleted the zipped folder file just go through the `rust-actix-sqlx-graphql-boilerplate` folder, zip it as `framework.20231117.zip` and put the zipped file under your root directory folder and build it.

## Current Features
- Generate a random token of any given length using the `-b` flag
    - `gen token` will generate a random key of length 32 bytes
        -  Master key will be converted to un-padded base64 url encoding
        -  Another string will be generated that can be used as your bearer token matching the master key
    - `gen framework $SET_PROJECT_PATH` will generate the latest rust boilerplate code.
      - see [https://github.com/enigs/rust-actix-sqlx-graphql-boilerplate](https://github.com/enigs/rust-actix-sqlx-graphql-boilerplate) for more details