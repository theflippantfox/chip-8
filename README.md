# Chip-8 in Go
_My implimentation of Chip-8 interpreter in golang in order for me to get familiar with workings of a CPU and basics of emulation_

## Features
 - Uses terminal itself as the display

## Known Issues
 - For some reason the topmost row of pixels is displayed at the bottom

## Getting Started
Firstly, clone the repo to your desired location
```shell

git clone https://github.com/theflippantfox/chip-8
```

After that change your current directory to the chip-8. For most it should be
```shell

cd chip-8
```

Then sync the libraries using
```shell

go mod tidy
```

You can now test the program using go run chip8 tests/[name of the test]
Example:
```shell
go run chip8 roms/1-chip8-logo.ch8
```

To build the program and get an executable file run
```shell
go build chip8
```
