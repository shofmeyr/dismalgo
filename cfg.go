package main

import (
	"fmt"
	"flag"
)

type Cfg struct {
	rseed          *int64
	numIters       *int
	numAgents      *int
	avMaxCsmp      *float64
	avMaxProd      *float64
	prdrSampleSize *int
	verboseFlags   *int
}

var _cfg = Cfg{
	rseed:          flag.Int64("d", 31, "Random seed"),
	numAgents:       flag.Int("n", 100, "Number of agents"),
	numIters:        flag.Int("i", 100000, "Number of iterations"),
	avMaxCsmp:       flag.Float64("c", 10, "Max. consumption"),
	avMaxProd:       flag.Float64("p", 10, "Max. production"),
	prdrSampleSize:  flag.Int("z", 10, "Sample size for selecting producer"),
	verboseFlags:    flag.Int("v", 0, "Verbose flags"),
}

func (cfg* Cfg) Load() {
	flag.Parse()
}

func (cfg* Cfg) Print() {
	//	fmt.Printf("%-30s %v\n", "Random seed", reflect.ValueOf(Rseed).Elem().Interface())
	fmt.Printf("Random seed        %d\n", *cfg.rseed)
	fmt.Printf("Number agents      %d\n", *cfg.numAgents)
	fmt.Printf("Number iterations  %d\n", *cfg.numIters)
	fmt.Printf("Max. consumption   %.2f\n", *cfg.avMaxCsmp)
	fmt.Printf("Max. production    %.2f\n", *cfg.avMaxProd)
	fmt.Printf("Sample size        %d\n", *cfg.prdrSampleSize)
	fmt.Printf("Verbose flags      %d\n", *cfg.verboseFlags)
}
