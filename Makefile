
default: test

rom2car: rom2car.go
	go build -o $@ $^
	go vet $^

rom8: a8rom.cc
	mos-atari8-stdcart-clang -Wall -Werror -Os -DROMSIZE=8 -Wl,--whole-archive -o $@ $^

rom16: a8rom.cc
	mos-atari8-stdcart-clang -Wall -Werror -Os -DROMSIZE=16 -Wl,--whole-archive -o $@ $^

rom32: a8rom.cc
	mos-atari8-xegs-clang -Wall -Werror -Os -DROMSIZE=32 -Wl,--whole-archive -o $@ $^

rom512: a8rom.cc
	mos-atari8-xegs-clang -Wall -Werror -Os -DROMSIZE=512 -Wl,--whole-archive -o $@ $^

# debugging/testing, local-test.sh should be something created for your specific env
test:	rom2car rom8 rom16 rom32 rom512
	if [ -x ./local-test.sh ] ; then ./local-test.sh $^ ; fi

clean:
	rm -f *~
