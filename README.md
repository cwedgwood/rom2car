# rom2car

Convert Atari 400/800/XL/XE ROM files to `.car` format.

Emulators such as [atari800](https://github.com/atari800/atari800) and
[Altirra](https://www.virtualdub.org/altirra.html) use the `.car`
format to automatically configure the bank-switching mechanism.

See
https://github.com/atari800/atari800/blob/ATARI800_5_0_0/DOC/cart.txt
for more details.

The input size must be a power of 2 between 8KiB and 1024KiB.

8KiB and 16KiB files are assumed to be standard 8KiB and 16KiB
cartridges, respectively.

Files 32KiB and larger are assumed to be XEGS banked cartridges.

A typical use case for this would be to convert the generated output
files from [llvm-mos](https://github.com/llvm-mos/)
[atari8-stdcard](https://github.com/llvm-mos/llvm-mos-sdk/tree/main/mos-platform/atari8-stdcart)
and
[atari8-xegs](https://github.com/llvm-mos/llvm-mos-sdk/tree/main/mos-platform/atari8-xegs)
build targets.

## Usage

## Atari 400/800 "Right-Cartridge" Work-Around

The Atari 800 has two cartridge slots, a left and a right, for this
reason the Atari 400/800 "OS A" ROM checks 0x9ffc (right-cartridge
present) for a zero-value indicating a right-cartridge is present, at
which point it will initialize and optionally start this cartridge.

This means when using a 16KiB (or larger bank-switched) cartridge if
we have a zero byte in this location, it will not boot properly on
Atari 400/800s.  XL and later models no longer have a ROM which
contains this check.

These bytes are at the end the banks, often unused so can be set to
0xff to avoid this.  If the byte do happen to be used, they will
likely be non-zero, an exception to this is where one or more banks
use these bytes and require this to be non-zero.
