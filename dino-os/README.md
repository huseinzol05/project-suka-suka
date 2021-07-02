# Dino-OS

An Operating system that only run Chrome Dino, kernel based on https://github.com/FRosner/FrOS/tree/minimal-c-kernel

Inspired by [I made an entire OS that only runs Tetris](https://www.youtube.com/watch?v=FaILnmUYS_U&t=97s).

## Installation

### For Mac

1. Install QEMU, NASM, and cross-compiler GCC,

```bash
brew install nasm qemu
brew tap nativeos/i386-elf-toolchain
brew install i386-elf-binutils i386-elf-gcc
```

After that should able,

```bash
i386-elf-gcc --version
```