package main

// this is so gross

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"strings"
)

//go:embed input
var input []byte

type Searcher struct {
	position          string
	destination       string
	timeToDestination int
}

type SearchState struct {
	Me              Searcher
	Elephant        Searcher
	includeElephant bool
	time            int
	flow            int
	totalReleased   int
}

func main() {
	adjacency, nameToIndex, flowRates := parseInput()
	distanceMatrix := createDistanceMatrix(adjacency)

	var possibleTargets []string
	for k := range flowRates {
		possibleTargets = append(possibleTargets, k)
	}

	part1MaxReleased := search(SearchState{
		Me: Searcher{
			position: "AA",
		},
		includeElephant: false,
	}, 30, possibleTargets, nameToIndex, flowRates, distanceMatrix)
	fmt.Println(part1MaxReleased)

	part2MaxReleased := search(SearchState{
		Me: Searcher{
			position: "AA",
		},
		Elephant: Searcher{
			position: "AA",
		},
		includeElephant: true,
	}, 26, possibleTargets, nameToIndex, flowRates, distanceMatrix)
	fmt.Println(part2MaxReleased)
}

func search(state SearchState, maxTime int, possibleTargets []string, nameToIndex map[string]int, flowRates map[string]int, distanceMatrix [][]int) SearchState {
	var stateToReturn SearchState
	for {
		// check for me needing a new destination and set up that branch
		if state.Me.destination == "" {
			for i, destination := range possibleTargets {
				nextState := state

				nextState.Me.destination = destination
				posIndex := nameToIndex[nextState.Me.position]
				dstIndex := nameToIndex[nextState.Me.destination]
				nextState.Me.timeToDestination = distanceMatrix[posIndex][dstIndex]

				dupSlice := make([]string, len(possibleTargets))
				copy(dupSlice, possibleTargets)
				dupSlice = append(dupSlice[:i], dupSlice[i+1:]...)

				if state.includeElephant && state.Elephant.destination == "" {
					for j, destination := range dupSlice {
						nextState.Elephant.destination = destination
						posIndex := nameToIndex[nextState.Elephant.position]
						dstIndex := nameToIndex[nextState.Elephant.destination]
						nextState.Elephant.timeToDestination = distanceMatrix[posIndex][dstIndex]
						dupSlice2 := make([]string, len(dupSlice))
						copy(dupSlice2, dupSlice)
						dupSlice2 = append(dupSlice2[:j], dupSlice2[j+1:]...)

						finalState := search(nextState, maxTime, dupSlice2, nameToIndex, flowRates, distanceMatrix)
						if finalState.totalReleased > stateToReturn.totalReleased {
							stateToReturn = finalState
						}
					}
				} else {
					finalState := search(nextState, maxTime, dupSlice, nameToIndex, flowRates, distanceMatrix)
					if finalState.totalReleased > stateToReturn.totalReleased {
						stateToReturn = finalState
					}
				}
			}
		}

		if state.includeElephant && state.Elephant.destination == "" {
			nextState := state

			for i, destination := range possibleTargets {
				nextState.Elephant.destination = destination
				posIndex := nameToIndex[nextState.Elephant.position]
				dstIndex := nameToIndex[nextState.Elephant.destination]
				nextState.Elephant.timeToDestination = distanceMatrix[posIndex][dstIndex]

				dupSlice := make([]string, len(possibleTargets))
				copy(dupSlice, possibleTargets)
				dupSlice = append(dupSlice[:i], dupSlice[i+1:]...)

				finalState := search(nextState, maxTime, dupSlice, nameToIndex, flowRates, distanceMatrix)
				if finalState.totalReleased > stateToReturn.totalReleased {
					stateToReturn = finalState
				}
			}
		}

		// tick
		state.time++
		state.Me.timeToDestination--
		state.totalReleased += state.flow

		// turn valve
		if state.Me.position == state.Me.destination {
			state.Me.destination = ""
			state.flow += flowRates[state.Me.position]
		}

		if state.Me.timeToDestination == 0 {
			state.Me.position = state.Me.destination
		}

		if state.includeElephant {
			state.Elephant.timeToDestination--
			if state.Elephant.timeToDestination == 0 {
				state.Elephant.position = state.Elephant.destination
			}

			if state.Elephant.position == state.Elephant.destination {
				state.Elephant.destination = ""
				state.flow += flowRates[state.Elephant.position]
			}
		}

		if state.time == maxTime {
			if state.totalReleased > stateToReturn.totalReleased {
				stateToReturn = state
			}
			return stateToReturn
		}

		//fmt.Println(state)
	}
}

func createDistanceMatrix(adjList [][]int) [][]int {
	numVertices := len(adjList)
	distanceMatrix := make([][]int, numVertices)
	for i := 0; i < numVertices; i++ {
		distanceMatrix[i] = make([]int, numVertices)
	}

	for sourceVertex := 0; sourceVertex < numVertices; sourceVertex++ {
		// initialize distance and visited arrays
		distance := make([]int, numVertices)
		visited := make([]bool, numVertices)
		for i := 0; i < numVertices; i++ {
			distance[i] = math.MaxInt32
			visited[i] = false
		}
		distance[sourceVertex] = 0

		// breadth-first partOneSearch
		queue := []int{sourceVertex}
		for len(queue) > 0 {
			vertex := queue[0]
			queue = queue[1:]
			visited[vertex] = true

			for _, adjVertex := range adjList[vertex] {
				if !visited[adjVertex] {
					distance[adjVertex] = distance[vertex] + 1
					queue = append(queue, adjVertex)
					visited[adjVertex] = true
				}
			}
		}

		// update distance matrix
		for i := 0; i < numVertices; i++ {
			distanceMatrix[sourceVertex][i] = distance[i]
		}
	}

	return distanceMatrix
}

func parseInput() ([][]int, map[string]int, map[string]int) {
	var adjacency [][]int
	nameToValveIndex := map[string]int{}
	flowRates := map[string]int{}

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		s := scanner.Text()

		var (
			valve    string
			flowRate int
		)
		fmt.Sscanf(s, "Valve %s has flow rate=%d", &valve, &flowRate)

		valveIndex, ok := nameToValveIndex[valve]
		if !ok {
			adjacency = append(adjacency, []int{})
			valveIndex = len(adjacency) - 1
			nameToValveIndex[valve] = valveIndex
		}

		if flowRate > 0 {
			flowRates[valve] = flowRate
		}

		_, neighborsStr, _ := strings.Cut(scanner.Text(), "valve")
		for _, str := range strings.Split(strings.TrimSpace(neighborsStr[1:]), ", ") {
			dstIndex, ok := nameToValveIndex[str]
			if !ok {
				adjacency = append(adjacency, []int{})
				dstIndex = len(adjacency) - 1
				nameToValveIndex[str] = dstIndex
			}
			adjacency[valveIndex] = append(adjacency[valveIndex], dstIndex)
		}

		valveIndex++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return adjacency, nameToValveIndex, flowRates
}
