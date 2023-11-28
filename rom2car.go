// rom2car
//
// Convert a '.ROM' file to a '.CAR' file suitable for emulators.  This means adding a header and appropriate checksum.

// >=16KiB carts on 400/800 will not work if '9ffc' is zero, this is because the Atari 400/800 "OS A" supports a "right-cartridge"
// and zero there creates a false detection, attempting to boot it which usually will hand.  To avoid this if we detect zeros in
// anything but the last bank (which is the left cartridge, and will be mapped to 0xa000) path to 0xff.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	fDebug   = flag.Bool("debug", false, "Show debugging")
	fOutput  = flag.String("output", "", "Output filename (guess)")
	fRCfix   = flag.String("rcfix", "auto", "Fix 'OS A' right-cartrige detection (auto/never/always)")
	fVerbose = flag.Bool("verbose", false, "Be verbose")
)

type rcfix int

const (
	rcfixAuto rcfix = iota
	rcfixNever
	rcfixAlways
)

var rcfixMode rcfix

func parseFixup() {
	switch strings.ToLower(*fRCfix) {
	case "auto":
		rcfixMode = rcfixAuto
	case "no":
		fallthrough
	case "never":
		rcfixMode = rcfixNever
	case "yes":
		fallthrough
	case "always":
		rcfixMode = rcfixAlways
	default:
		fatalf("Unvalid rcfix '%s'\n", *fRCfix)
	}
	debugf("rcfix: %s\n", *fRCfix)
}

func validSize(n int) bool {
	// Ensure n is a power of 2
	if (n & (n - 1)) != 0 {
		return false
	}
	n >>= 10 // KiB
	if n < 8 || n > 1024 {
		return false
	}
	return true
}

func warnf(m string, args ...any) {
	fmt.Fprintf(os.Stdout, m, args...)
}

func fatalf(m string, args ...any) {
	warnf(m, args...)
	os.Exit(1)
}

func debugf(m string, args ...any) {
	if !*fDebug {
		return
	}
	warnf(m, args...)
}

func verbosef(m string, args ...any) {
	if !*fVerbose {
		return
	}
	warnf(m, args...)
}

func errorf(err error) {
	if err != nil {
		fatalf("%v", err)
	}
}

// deriveOutput tries to sensibly derive an output filename, either by replacing a .rom suffix with .car or appending .car to
// whatever we started with
func deriveOutput(fn string) string {

	// do we have "*.rom" ? if so strip that
	if n := strings.LastIndex(fn, "."); n > 0 {
		if strings.ToLower(fn[n:]) == ".rom" {
			fn = fn[:n]
		}
	}
	return fn + ".car"
}

func main() {
	flag.Parse()
	parseFixup()

	if *fDebug {
		*fVerbose = true
	}

	var inFn, outFn string

	if na := len(flag.Args()); na != 1 {
		if na > 1 {
			fatalf("One filename must be given")
		}
	}
	inFn = flag.Args()[0]

	if *fOutput != "" {
		outFn = *fOutput
	} else {
		outFn = deriveOutput(inFn)
	}

	b, err := os.ReadFile(inFn)
	errorf(err)
	verbosef("Read '%s' size %d\n", inFn, len(b))

	if !validSize(len(b)) {
		fatalf("Input size of %d is not appropriate", len(b))
	}

	// https://raw.githubusercontent.com/atari800/atari800/ATARI800_5_0_0/DOC/cart.txt
	types := map[int]uint8{
		8:    1,
		16:   2,
		32:   12,
		64:   13,
		128:  14,
		256:  23,
		512:  24,
		1024: 25,
	}
	t := types[len(b)/1024]
	if t == 0 {
		fatalf("Unable to determine type for size %d\n", len(b))
	}

	if rcfixMode != rcfixNever {
		fixupFailed := false
		// patch the right-cartrige detect bytes so OS A works
		for o := 0; o < (len(b) - 8*1024); o += 8 * 1024 {
			if b[o+0x1ffc] == 0 {

				// if in auto mode, check surrounding bytes are zero
				if rcfixMode == rcfixAuto {
					nZero := 0
					for i := 0x1ff0; i <= 0x1fff; i++ {
						if b[o+i] != 0 {
							nZero++
						}
					}
					if nZero > 0 {
						warnf("rxfix auto; saw ~%d surrounding zeroe near 0x%05x\n", nZero, o+0x1ffc)
						fixupFailed = true
					}
				}

				b[o+0x1ffc] = 0xff
				debugf("Unzeroed 0x%05x\n", o+0x1ffc)
			}
		}
		if fixupFailed {
			fatalf("Cannot reliably fixup/patch ROM\n")
		}
	}

	carHeader := [16]byte{0x43, 0x41, 0x52, 0x54} // "CART"
	carHeader[7] = t                              // type

	// checksum
	var cksum uint32
	for _, v := range b {
		cksum += uint32(v)
	}
	for i := 0; i < 4; i++ {
		carHeader[12-1-i] = byte(cksum)
		cksum >>= 8
	}

	err = os.WriteFile(outFn, append(carHeader[:], b...), 0644)
	errorf(err)
	verbosef("Wrote '%s' size %d\n", outFn, len(carHeader)+len(b))
}
