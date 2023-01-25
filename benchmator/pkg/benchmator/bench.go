package benchmator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Kochanac/coursework2023/benchmator/pkg/util"
	"github.com/influxdata/tdigest"
)

type Benchmator struct {
	results    *tdigest.TDigest
	resultsArr []time.Duration
	stop       func()

	f func()
}

type Stats struct {
	mean time.Duration
	p90  time.Duration
	p99  time.Duration

	arr []time.Duration
}

func (s Stats) String() string {
	return fmt.Sprintf("mean %d us; 0.9 %d us; 0.99 %d us",
		s.mean.Microseconds(), s.p90.Microseconds(), s.p99.Microseconds())
}

func (s Stats) Save(f io.Writer) error {
	bts, err := json.Marshal(s.arr)
	if err != nil {
		return err
	}
	_, err = f.Write(bts)
	if err != nil {
		return err
	}
	return nil
}

func Run(f func(), d time.Duration, callsPerSecond int) Stats {
	b := Benchmator{
		results: tdigest.New(),
		stop:    nil,
		f:       f,
	}
	return b.run(d, callsPerSecond)
}

func (b *Benchmator) run(d time.Duration, callsPerSecond int) Stats {
	b.bench(callsPerSecond)
	time.Sleep(d)
	b.stop()

	return Stats{
		mean: time.Duration(b.results.Quantile(0.5)),
		p90:  time.Duration(b.results.Quantile(0.90)),
		p99:  time.Duration(b.results.Quantile(0.99)),
		arr:  b.resultsArr,
	}
}

func (b *Benchmator) collect(rt chan time.Duration) {
	for t := range rt {
		b.resultsArr = append(b.resultsArr, t)
		b.results.Add(float64(t), 1)
	}
}

func (b *Benchmator) bench(callsPerSecond int) {
	ticker := time.NewTicker(time.Second / time.Duration(callsPerSecond))
	log.Printf("TARGET %s / %d f/s", time.Second/time.Duration(callsPerSecond), callsPerSecond)

	rt := make(chan time.Duration, 1e5)
	wg := sync.WaitGroup{}
	fc := func() {
		wg.Add(1)
		defer wg.Done()
		t1 := time.Now()
		b.f()
		rt <- time.Since(t1)
	}
	go b.collect(rt)

	quit := make(chan struct{})
	go util.StartAsyncTimer(ticker.C, quit, fc)

	b.stop = func() {
		ticker.Stop()
		close(quit)
		time.Sleep(time.Millisecond * 10)
		wg.Wait()
		close(rt)
	}
}
