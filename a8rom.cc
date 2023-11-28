// Generate a cartridge image that print it's size

#include <stdio.h>

#define STR_HELPER(x) #x
#define STR(x) STR_HELPER(x)

asm(".weak __cart_rom_size \n __cart_rom_size = " STR(ROMSIZE));

int main() {
  puts("Cartridge " STR(ROMSIZE));
  for (;;)
    ;
}
