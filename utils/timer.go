package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Timer struct {
	durations []int64
	name      string
	silent    bool
}

func NewTimer(name string) *Timer {
	return &Timer{
		durations: []int64{},
		name:      name,
	}
}

// Stops TimeIt from printing the time taken each call.
func (t *Timer) SetSilent() *Timer {
	t.silent = true
	return t
}

// Runs the function and notes down its runtime
func (t *Timer) TimeIt(fun func()) {
	now := time.Now()
	fun()
	timeTaken := time.Since(now)
	t.durations = append(t.durations, timeTaken.Milliseconds())
	if !t.silent {
		fmt.Printf("%s time taken: %v\n", t.name, timeTaken)
	}
}

// Profiles a function. Since profiling will affect the performance of the
// function, timing info will not be recorded.
func (t *Timer) ProfileIt(fun func(), filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	fun()

	pprof.StopCPUProfile()
	// mark as invalid point so we can skip it when graphing, but don't use
	// the wrong x coordinate.
	t.durations = append(t.durations, -1)
	fmt.Println("profile", filename, "done")
}

// Saves the timing info recorded by all the runs of `TimeIt` in a file.
// It's just the durations in miliseconds, separated by newlines
func (t *Timer) Save(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, duration := range t.durations {
		if duration < 0 {
			continue
		}
		_, err := f.WriteString(strconv.FormatInt(duration, 10))
		if err != nil {
			panic(err)
		}
		f.WriteString("\n")
	}
}

// Prints the timing info out in a csv like format. It's like save, except to
// stdout.
func (t *Timer) Echo() {
	var sb strings.Builder
	for i, duration := range t.durations {
		if duration < 0 {
			continue
		}
		_, err := sb.WriteString(strconv.FormatInt(duration, 10))
		if err != nil {
			panic(err)
		}
		if i != len(t.durations)-1 {
			sb.WriteString(", ")
		}
	}
	fmt.Println(t.name, "timings:")
	fmt.Println(sb.String())
}

// Adds the duration data to the line graph and returns the number of items
// added.
func (t *Timer) AddToLineGraph(line *charts.Line) int {
	data := make([]opts.LineData, len(t.durations))

	for i, duration := range t.durations {
		if duration < 0 {
			data[i] = opts.LineData{
				Value: nil,
			}
		} else {
			data[i] = opts.LineData{
				Value: duration,
			}
		}
	}

	line.AddSeries(t.name, data)

	return len(t.durations)
}

func GraphTimers(filename, title string, timers ...*Timer) {
	chart := charts.NewLine()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	maxlen := 0
	for _, timer := range timers {
		len := timer.AddToLineGraph(chart)
		if len > maxlen {
			maxlen = len
		}
	}

	x := make([]string, maxlen)
	for i := range maxlen {
		x[i] = strconv.Itoa(i)
	}
	chart.SetXAxis(x)
	chart.SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			Smooth:       opts.Bool(true),
			ConnectNulls: opts.Bool(true),
		}),
	)
	chart.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Iteration",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Time taken (ms)",
		}),
	)

	f, _ := os.Create(filename)
	chart.Render(f)
	fmt.Println("Written output to", filename)
}
