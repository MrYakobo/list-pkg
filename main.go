package main

import (
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type el struct {
	name       string
	size       float64
	sizeString string
	date       time.Time
	// compare    func(other el) float64
	timeSort bool
}

func create(name string, size string, date string, timeSort bool) el {

	multiplier := 1e0
	if strings.Contains(size, "MiB") {
		multiplier = 1e6
	} else if strings.Contains(size, "KiB") {
		multiplier = 1e3
	}
	val, _ := strconv.ParseFloat(strings.Trim(size[0:len(size)-3], " "), 64)
	val *= multiplier

	// RFC: Jan 2 15:04:05 2006 MST
	//Fri 20 Apr 2018 01:17:49 PM CEST
	d, _ := time.Parse("Mon 02 Jan 2006 15:04:05 PM CEST", strings.Trim(date, " "))

	el := el{strings.Trim(name, " "), val, strings.Trim(size, " "), d, timeSort}

	return el
}

func (e *el) DateComp(other el) float64 {
	return (float64)(-e.date.Unix() + other.date.Unix())
}

func (e *el) SizeComp(other el) float64 {
	return -e.size + other.size
}

type arr []el

func (a arr) Len() int {
	return len(a)
}

func (a arr) Less(i, j int) bool {
	// return a[i].SizeComp(a[j]) < 0
	// return a[i].compare(a[j]) < 0
	if a[i].timeSort {
		return a[i].DateComp(a[j]) < 0
	}
	return a[i].SizeComp(a[j]) < 0
}

func (a arr) Swap(i, j int) {
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

func main() {
	flag.Usage = func() {

	}
	t := flag.Bool("t", false, "Sort programs by install/update time instead of size (default)")
	flag.Parse()

	fmtPrint := "Sorted by size\n"
	if *t {
		fmtPrint = "Sorted by install/update time\n"
	}

	fmt.Println(fmtPrint)
	str, _ := exec.Command("sh", "-c", "pacman -Qei | grep -E 'Install Date|Name|Installed Size' | cut -d : -f 2,3,4 ").Output()
	Arr := strings.Split((string)(str), "\n")

	aarr := (arr)(make([]el, len(Arr)/3))

	for i := 0; i < len(Arr)-2; i += 3 {
		a := create(Arr[i], Arr[i+1], Arr[i+2], *t)
		aarr[i/3] = a
	}

	sort.Sort(aarr)

	for i := 0; i < len(aarr); i++ {
		fmt.Printf("%s\n\t%s\n\t%s\n", aarr[i].name, aarr[i].date.Format("2 Jan 2006 15:04:05"), aarr[i].sizeString)
	}
	fmt.Printf("\n%s %d\n", "Total amount of installed programs:", len(aarr))
}
