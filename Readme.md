# aoccli
A CLI tool to help you with various Advent of Code-related manual tasks.

## Building
Run `GOOS=<os> GOARCH=<arch> go build -o aoccli` to compile the cli to a binary. For example
```
# linux amd64
GOOS=linux GOARCH=amd64 go build -o aoccli

# mac amd64
GOOS=darwin GOARCH=amd64 go build -o aoccli
```
then move the binary to somewhere that's included in your path, alternatively, add this folder
to your path.

## Configuring
Run `aoccli configure` to set up configuration needed to run the application. Here, you set things
like the root folder of your aoc workspace and the session token retrieved from the aoc website.
Run `aoccli configure --help` for all available configs.

## Scaffolding
After configuring the CLI, you can scaffold a day (make sure you set year and day when configuring). 
Run `aoccli scaffold` which creates this folder structure:
```
- <year>
    - <day>
        - input
        - s1.py
        - s2.py
```

## Solving and submitting
Once you have created a solution to your problem, you can run `aoccli solve 1/2` to run the solution and
save the stderr output to a file. You can then run `aoccli submit` to submit the solution. The CLI will
try to keep track of the current solution to submit (once the first problem is submitted and correct, the 
CLI will submit the second one the next time you run the command).