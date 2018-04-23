# Notes

Some notes about making emulators.

## Links

* NES
 * [NES Dev Wiki - 6502 Instructions](https://wiki.nesdev.com/w/index.php/6502_instructions)
 * [NES Dev Wiki - 6502 Microprocessor Specification](http://nesdev.com/6502.txt)

* Assembly Utilities
 * [Easy 6502 - Introduction to 6502 Assembly Programming](http://skilldrick.github.io/easy6502/index.html)


## Example 6502 Program

This is a simple example program showing how different assembly
programs get compiled & assembled. This program just sets three
pixels on the screen to three different colors. The pixels are
the three top, left-most pixels in the screen.


### Program:
```
LDA #$01
STA $0200
LDA #$05
STA $0201
LDA #$08
STA $0202
```

### Hex Dump:
```
0600: a9 01 8d 00 02 a9 05 8d 01 02 a9 08 8d 02 02 
```

### Disassembly:
```
Address  Hexdump   Dissassembly
-------------------------------
$0600    a9 01     LDA #$01
$0602    8d 00 02  STA $0200
$0605    a9 05     LDA #$05
$0607    8d 01 02  STA $0201
$060a    a9 08     LDA #$08
$060c    8d 02 02  STA $0202
```
