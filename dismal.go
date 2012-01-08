package main

import (
	"fmt"
	"rand"
	"time"
	"reflect"
)

const dismalVersion = 0.1

type Agent struct {
	id          int
	maxCsmp     float64
	maxProd     float64
	money       float64
	moneyGained float64
	unsoldProd  float64
	csmp        float64
	totCsmp     float64
	totProd     float64
	price       float64
	adjust      float64
}

var (
	_agents []Agent
	_prdrs  []int
	_csmrs  []int
)

func main() {
	fmt.Printf("DISMAL ECONOMIC MODEL (Version %.2f)\n", dismalVersion)
	_cfg.Load()
	_cfg.Print()
	initSim()
	fmt.Printf("%8s%7s%7s%7s%7s%7s%7s%7s%7s%7s%7s%7s\n",
		"t", "mx $", "mn $", "av PP", "mx PP", "mn PP", "av C", "mx C", "mn C", 
		"av P", "mx P", "mn P")

	iterStep := *_cfg.numIters / 25
	if iterStep == 0 {
		iterStep = 1
	}
	startTime := time.Nanoseconds()
	for iters := 0; iters < *_cfg.numIters; iters++ {
		_prdrs = _prdrs[:0]
		_csmrs = _csmrs[:0]
		for i := range _agents {
			_agents[i].update()
		}

		for len(_csmrs) > 0 && len(_prdrs) > 0 {
			csmrI := rand.Intn(len(_csmrs))
			prdrI, found := findCheapestPrdr(csmrI)
			if found {
				consume(csmrI, prdrI)
			} else {
				if len(_prdrs) == 1 && len(_csmrs) == 1 && _prdrs[0] == _csmrs[0] {
					continue
				}
			}
		}

		for i := range _agents {
			computePrice(iters + 1, &_agents[i])
		}

		if iters % iterStep == 0 {
			computeStats(iters + 1, true)
		}

//		for i := range _agents {
//			fmt.Println(_agents[i].String())
//		}
	}

	fmt.Printf("Elapsed time %.2fs\n", (float64)(time.Nanoseconds()-startTime)/1e9)
}

func initSim() {
	rand.Seed(*_cfg.rseed)
	fmt.Printf("num agents %d\n", *_cfg.numAgents)
	_agents = make([]Agent, *_cfg.numAgents)
	_prdrs = make([]int, 0, *_cfg.numAgents)
	_csmrs = make([]int, 0, *_cfg.numAgents)

	for i := range _agents {
		_agents[i].init(i)
	}

}

func (agent *Agent) init(id int) {
	agent.id = id
	agent.maxCsmp = *_cfg.avMaxCsmp
	agent.maxProd = *_cfg.avMaxProd
	agent.money = 1.0
	agent.moneyGained = 0
	agent.unsoldProd = 1.0
	agent.csmp = 0
	agent.totCsmp = 0
	agent.totProd = 0
	agent.price = agent.money
	agent.adjust = 0.001
}

func (agent *Agent) String() string {
	return fmt.Sprintf("%5d", agent.id) +
		fmt.Sprintf("%8.2f", agent.maxCsmp) + 
		fmt.Sprintf("%8.2f", agent.maxProd) +
		fmt.Sprintf("%8.2f", agent.money) +
		fmt.Sprintf("%8.2f", agent.moneyGained) +
		fmt.Sprintf("%8.2f", agent.unsoldProd) +
		fmt.Sprintf("%8.2f", agent.csmp) + 
		fmt.Sprintf("%8.2f", agent.totCsmp) + 
		fmt.Sprintf("%8.2f", agent.totProd) +
		fmt.Sprintf("%8.2f", agent.price) + 
		fmt.Sprintf("%8.2f", agent.adjust)
}

func (agent *Agent) update() {
	agent.csmp = 0
	if agent.money > 0 {
		_csmrs = append(_csmrs, agent.id)
	}
	_prdrs = append(_prdrs, agent.id)
	agent.unsoldProd = agent.maxProd
	agent.money += agent.moneyGained
	agent.moneyGained = 0
}

func findCheapestPrdr(csmrI int) (prdrI int, found bool) {
	minPrice := (float64)(1e9)
	found = false
	for i := 0; i < *_cfg.prdrSampleSize; i++ {
		rndPrdr := rand.Intn(len(_prdrs))
		if _prdrs[rndPrdr] == _csmrs[csmrI] {
			continue
		}
		prdr := &_agents[_prdrs[rndPrdr]]
		if prdr.price < minPrice {
			minPrice = prdr.price
			prdrI = rndPrdr
			found = true
		}
	}
	return
}

func consume(csmrI, prdrI int) {
	csmr := &_agents[_csmrs[csmrI]]
	prdr := &_agents[_prdrs[prdrI]]
	if csmr.id == prdr.id {
		return
	}
	csmp := csmr.maxCsmp - csmr.csmp
	csmpCost := csmp * prdr.price
	if csmpCost > csmr.money {
		csmp = csmr.money / prdr.price
	}
	if prdr.unsoldProd < csmp {
		csmp = prdr.unsoldProd
	}
	csmpCost = csmp * prdr.price

	prdr.unsoldProd -= csmp
	if prdr.unsoldProd < 0.00001 {
		prdr.unsoldProd = 0
	}
	prdr.totProd += csmp
	prdr.moneyGained += csmpCost

	csmr.money -= csmpCost
	if csmr.money < 0.00001 {
		csmr.money = 0
	}
	csmr.csmp += csmp
	csmr.totCsmp += csmp

	if prdr.unsoldProd == 0 {
		_prdrs[prdrI] = _prdrs[len(_prdrs) - 1]
		_prdrs = _prdrs[:len(_prdrs) - 1]
	}
	
	if csmr.money == 0 || csmr.csmp >= csmr.maxCsmp {
		_csmrs[csmrI] = _csmrs[len(_csmrs) - 1]
		_csmrs = _csmrs[:len(_csmrs) - 1]
	}
}

func computePrice(iters int, agent *Agent) {
	exptdProd := agent.totProd / (float64)(iters)
	priceChange := rand.Float64() * (exptdProd - agent.unsoldProd) / agent.maxProd * agent.adjust
	agent.price += priceChange
	minPrice := 0.00001
	if agent.price < minPrice {
		agent.price = minPrice
	}
}

func computeStats(iters int, showRound bool) {
	
//	fmt.Printf("%8s%7s%7s%7s%7s%7s%7s%7s%7s%7s%7s%7s\n",
//		"t", "mx $", "mn $", "av PP", "mx PP", "mn PP", "av C", "mx C", "mn C", 
//		"av P", "mx P", "mn P")

	avMoney, minMoney, maxMoney := computeStat("money")
	fmt.Printf("%7.2f%7.2f%7.2f\n", avMoney, minMoney, maxMoney)
}

func computeStat(stat string) (av, mn, mx float64) {
	av = 0.0
	mn = 1e9
	mx = 0.0
	for _, agent := range _agents {
		s := reflect.ValueOf(&agent).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			if stat == typeOfT.Field(i).Name {
				fmt.Printf("***")
			}
			fmt.Printf("%s\n", typeOfT.Field(i).Name)
//			f := s.Field(i)
//			fmt.Printf("%s \n", typeOfT.Field(i).Name)//, f.Interface())
		}
	}
	return
}


