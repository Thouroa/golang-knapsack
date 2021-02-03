package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
)

type Item struct {
	Name   string
	Weight int
	Value  int
	Count  int
}

type Knapsack struct {
	Limit int
	Items []Item
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", os.Args[0], " <data_file>")
		return
	}
	bt, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Failed to open file ", os.Args[1], ", Msg: ", err.Error())
	}
	kanpsack := Knapsack{}
	err = json.Unmarshal(bt, &kanpsack)
	if err != nil {
		log.Fatal("Json Parse Failed, Msg: ", err.Error())
	}

	Fractional(kanpsack)
	ZeroOrOne(kanpsack)
	Unboarded(kanpsack)
	Boarded(kanpsack)
}

// Fractional
func Fractional(pack Knapsack) {
	TotalValue := 0
	remainingSize := pack.Limit
	sort.SliceStable(pack.Items, func(i, j int) bool {
		return (pack.Items[i].Value * 100 / pack.Items[i].Weight) > (pack.Items[j].Value * 100 / pack.Items[j].Weight)
	})

	putItem := []string{}
	for _, item := range pack.Items {
		if item.Weight <= remainingSize {
			putItem = append(putItem, strconv.Itoa(item.Weight)+"kg "+item.Name)
			TotalValue += item.Value
			remainingSize -= item.Weight
		} else {
			putItem = append(putItem, strconv.Itoa(remainingSize)+"kg "+item.Name)
			TotalValue += item.Value * remainingSize / item.Weight
			remainingSize -= remainingSize
		}

		if remainingSize == 0 {
			break
		}
	}

	fmt.Println("Fractional:\t", TotalValue, putItem)
}

// Zero-Or-One
func ZeroOrOne(pack Knapsack) {
	cost := make([]int, pack.Limit+1)
	for i := range pack.Items {
		Value := pack.Items[i].Value
		Weight := pack.Items[i].Weight
		for j := pack.Limit; j >= Weight; j-- {
			if cost[j-Weight]+Value > cost[j] {
				cost[j] = cost[j-Weight] + Value
			}
		}
	}

	fmt.Println("Zero-Or-One:\t", cost[pack.Limit])
}

// Unboarded
func Unboarded(pack Knapsack) {
	cost := make([]int, pack.Limit+1)
	for i := range pack.Items {
		Value := pack.Items[i].Value
		Weight := pack.Items[i].Weight
		for count := 1; count*Weight <= pack.Limit; count++ {
			c_Value := Value * count
			c_Weight := Weight * count
			for j := pack.Limit; j >= c_Weight; j-- {
				if cost[j-c_Weight]+c_Value > cost[j] {
					cost[j] = cost[j-c_Weight] + c_Value
				}
			}
		}
	}

	fmt.Println("Unboarded:\t", cost[pack.Limit])
}

// Boarded
func Boarded(pack Knapsack) {
	cost := make([]int, pack.Limit+1)
	tmpPack := pack
	for _, i := range tmpPack.Items {
		for n := 1; n < i.Count; n++ {
			pack.Items = append(pack.Items, i)
		}
	}
	for i := range pack.Items {
		Value := pack.Items[i].Value
		Weight := pack.Items[i].Weight
		for j := pack.Limit; j >= Weight; j-- {
			if cost[j-Weight]+Value > cost[j] {
				cost[j] = cost[j-Weight] + Value
			}
		}
	}

	fmt.Println("Boarded:\t", cost[pack.Limit])
}
