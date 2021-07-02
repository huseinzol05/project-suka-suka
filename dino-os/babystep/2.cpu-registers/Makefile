build:
	nasm -f bin -o boot.bin boot.asm
	dd if=boot.bin of=boot.flp
	qemu-system-i386 -fda boot.flp