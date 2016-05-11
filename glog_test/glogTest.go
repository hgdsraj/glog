package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	tlog "github.com/tarm/glog"
)

type d []time.Duration

func (d d) Len() int           { return len(d) }
func (d d) Less(i, j int) bool { return d[i] < d[j] }
func (d d) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d d) String() (str string) {
	sort.Sort(d)
	mult := 1
outer:
	for {
		for _, ns := range []int{10, 25, 50} {
			dur := time.Duration(ns * mult)
			idx := sort.Search(len(d), func(i int) bool {
				return d[i] > dur
			})
			rem := len(d) - idx

			strt := fmt.Sprintf("   %8d samples longer than %v\n", rem, dur)
			if rem > len(d)*9/10 {
				str = strt
				continue
			}
			str += strt
			if rem < 1 {
				break outer
			}
		}
		mult *= 10
	}
	str = fmt.Sprintf("%v samples total\n", len(d)) + str

	return str
}

func main() {
	n := 100000
	rot := 3000
	flag.IntVar(&n, "n", n, "iterations to run")
	flag.IntVar(&rot, "size", rot, "max glog size")
	flag.Parse()

	var glogTimes, tlogTimes, logTimes, fmtTimes d

	glog.MaxSize = uint64(rot)
	tlog.MaxSize = uint64(rot)

	str := strings.Repeat("x", 1000)
	for i := 0; i < n; i++ {
		t0 := time.Now()
		glog.Info(str)
		t1 := time.Now()
		tlog.Info(str)
		t2 := time.Now()
		log.Println(str)
		t3 := time.Now()
		fmt.Fprintf(ioutil.Discard, "%s", str)
		t4 := time.Now()
		glogTimes = append(glogTimes, t1.Sub(t0))
		tlogTimes = append(tlogTimes, t2.Sub(t1))
		logTimes = append(logTimes, t3.Sub(t2))
		fmtTimes = append(fmtTimes, t4.Sub(t3))
	}
	fmt.Println("Glog times", glogTimes)
	fmt.Println()
	fmt.Println("Tlog times", tlogTimes)
	fmt.Println()
	fmt.Println("Log times", logTimes)
	fmt.Println()
	fmt.Println("Fmt times", fmtTimes)
}
