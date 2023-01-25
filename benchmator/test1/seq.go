package main

import (
	"math/rand"
)

type test1 struct {
	id   int64
	data int64
}

func sequence(i int64) test1 {
	return test1{
		id:   i + 1000000000*(i%8),
		data: i + 4824246592 + 1000000000*(i%8),
	}
}

const seqMin = 6978947857221617788

//const seqMax = seqMin + ROWNUM

func getRandFromSeq() test1 {
	n := rand.Int63n(ROWNUM)
	return sequence(n + seqMin)
}
