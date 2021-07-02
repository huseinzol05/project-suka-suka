build:
	nasm -f bin -o hello-world.bin hello-world.asm
	dd if=hello-world.bin of=hello-world.flp
	qemu-system-i386 -fda hello-world.flp