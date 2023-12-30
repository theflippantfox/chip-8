# Chip-8 in Go
_My implimentation of Chip-8 interpreter in golang in order for me to get familiar with workings of a CPU and basics of emulation_

## Table Of Contents

  * [Introduction](#introduction)
    * [CHIP-8](#chip-8)
    * [Features](#features)
    * [CHIP-8 Instruction Set](#chip-8-instruction-set)
  * [Build Guide](#build-guide)
  * [TODO](#todo)
  * [Known Issues](#known-issues)

# Introduction

## CHIP-8

Chip-8 is a simple, interpreted, programming language which was first used on some do-it-yourself computer systems in the late 1970s and early 1980s. The COSMAC VIP, DREAM 6800, and ETI 660 computers are a few examples. These computers typically were designed to use a television as a display, had between 1 and 4K of RAM, and used a 16-key hexadecimal keypad for input. The interpreter took up only 512 bytes of memory, and programs, which were entered into the computer in hexadecimal, were even smaller.

In the early 1990s, the Chip-8 language was revived by a man named Andreas Gustafsson. He created a Chip-8 interpreter for the HP48 graphing calculator, called Chip-48. The HP48 was lacking a way to easily make fast games at the time, and Chip-8 was the answer. Chip-48 later begat Super Chip-48, a modification of Chip-48 which allowed higher resolution graphics, as well as other graphical enhancements.
## Features
### Memory

The Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM, from location 0x000 (0) to 0xFFF (4095). The first 512 bytes, from 0x000 to 0x1FF, are where the original interpreter was located, and should not be used by programs.

Most Chip-8 programs start at location 0x200 (512), but some begin at 0x600 (1536). Programs beginning at 0x600 are intended for the ETI 660 computer.

Memory Map:
```
+---------------+= 0xFFF (4095) End of Chip-8 RAM
|               |
|               |
|               |
|               |
|               |
| 0x200 to 0xFFF|
|     Chip-8    |
| Program / Data|
|     Space     |
|               |
|               |
|               |
+- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
|               |
|               |
|               |
+---------------+= 0x200 (512) Start of most Chip-8 programs
| 0x000 to 0x1FF|
| Reserved for  |
|  interpreter  |
+---------------+= 0x000 (0) Start of Chip-8 RAM
```
### Registers 


Chip-8 has 16 general purpose 8-bit registers, usually referred to as Vx, where x is a hexadecimal digit (0 through F). There is also a 16-bit register called I. This register is generally used to store memory addresses, so only the lowest (rightmost) 12 bits are usually used.
The VF register should not be used by any program, as it is used as a flag by some instructions.
The VF register should not be used by any program, as it is used as a flag by some instructions.
The stack is an array of 16 16-bit values, used to store the address that the interpreter shoud return to when finished with a subroutine. Chip-8 allows for up to 16 levels of nested subroutines.
The program counter (PC) should be 16-bit, and is used to store the currently executing address. The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.


### Keyboard
The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad

### Display
The original implementation of the Chip-8 language used a 64x32-pixel monochrome display
Chip-8 draws graphics on screen through the use of sprites. A sprite is a group of bytes which are a binary representation of the desired picture. Chip-8 sprites may be up to 15 bytes, for a possible sprite size of 8x15.

### Timer & Sound

Chip-8 provides 2 timers, a delay timer and a sound timer.

The delay timer is active whenever the delay timer register (DT) is non-zero. This timer does nothing more than subtract 1 from the value of DT at a rate of 60Hz. When DT reaches 0, it deactivates.

The sound timer is active whenever the sound timer register (ST) is non-zero. This timer also decrements at a rate of 60Hz, however, as long as ST's value is greater than zero, the Chip-8 buzzer will sound. When ST reaches zero, the sound timer deactivates.

The sound produced by the Chip-8 interpreter has only one tone. The frequency of this tone is decided by the author of the interpreter.

## CHIP-8 Instruction Set

The original implementation of the Chip-8 language includes 36 different instructions, including math, graphics, and flow control functions.

Super Chip-48 added an additional 10 instructions, for a total of 46.

All instructions are 2 bytes long and are stored most-significant-byte first. In memory, the first byte of each instruction should be located at an even addresses. If a program includes sprite data, it should be padded so any instructions following it will be properly situated in RAM.

This document does not yet contain descriptions of the Super Chip-48 instructions. They are, however, listed below.

In these listings, the following variables are used:

nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
n or nibble - A 4-bit value, the lowest 4 bits of the instruction
x - A 4-bit value, the lower 4 bits of the high byte of the instruction
y - A 4-bit value, the upper 4 bits of the low byte of the instruction
kk or byte - An 8-bit value, the lowest 8 bits of the instruction

# Build Guide

# TODO

# Known Issues
