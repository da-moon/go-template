package testutils

import (
	"io"

	primitives "github.com/da-moon/go-template/internal/primitives"
)

const (
	// Magic constants (See Knuth: Seminumerical Algorithms)
	// Do not change! (They are for a maximal length LFSR)
	ranlen  = 55
	ranlen2 = 24
)

// randomReader contains the state for the random stream generator
type randomReader struct {
	blockSize        int
	randomDataBuffer []byte
	Data             []byte
	bytes            int
	pos              int
}

// NewRandomReader - make a new random stream generator
// random reader is used to generate random streams of bytes used
// for testing io operations
// https://github.com/ncw/stressdisk
func NewRandomReader(blockSize int) io.Reader {
	result := &randomReader{}
	result.blockSize = blockSize
	if blockSize <= 0 {
		result.blockSize = 2 * primitives.Mi
	}
	result.randomDataBuffer = make([]byte, result.blockSize+ranlen)
	result.Data = result.randomDataBuffer[0:result.blockSize]
	result.Data[0] = 1
	for i := 1; i < ranlen; i++ {
		result.Data[i] = 0xA5
	}
	// initial randomisation
	result.randomise()
	// start buffer empty
	result.bytes = 0
	result.pos = 0
	return result
}

// Read implements io.Reader for randomReader
func (r *randomReader) Read(p []byte) (int, error) {
	bytesToWrite := len(p)
	bytesWritten := 0
	for bytesToWrite > 0 {
		if r.bytes <= 0 {
			r.randomise()
		}
		chunkSize := bytesToWrite
		if bytesToWrite >= r.bytes {
			chunkSize = r.bytes
		}
		copy(p[bytesWritten:bytesWritten+chunkSize], r.Data[r.pos:r.pos+chunkSize])
		bytesWritten += chunkSize
		bytesToWrite -= chunkSize
		r.pos += chunkSize
		r.bytes -= chunkSize
	}
	return bytesWritten, nil
}

// randomise fills the random block up with randomness.
//
// This uses a random number generator from Knuth: Seminumerical
// Algorithms.  The magic numbers are the polynomial for a maximal
// length linear feedback shift register The least significant bits of
// the numbers form this sequence (of length 2**55).  The higher bits
// cause the sequence to be some multiple of this.
func (r *randomReader) randomise() {
	// copy the old randomness to the end
	copy(r.randomDataBuffer[r.blockSize:], r.randomDataBuffer[0:ranlen])
	// make a new random block
	d := r.randomDataBuffer
	for i := r.blockSize - 1; i >= 0; i-- {
		d[i] = d[i+ranlen] + d[i+ranlen2]
	}
	// Show we have some bytes
	r.bytes = r.blockSize
	r.pos = 0
}
