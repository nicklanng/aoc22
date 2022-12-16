package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math"
	"strings"
)

//go:embed input
var input []byte

type Searcher struct {
	position          int
	destination       int
	timeToDestination int
	waiting           bool
}

type State struct {
	Me            Searcher
	Elephant      Searcher
	visitedValves lib.Bitfield1[uint64]
	time          int
	flow          int
	totalReleased int
}

func (s *State) Tick() {
	s.time++
	s.totalReleased += s.flow
	s.Me.timeToDestination--
	s.Elephant.timeToDestination--
}

func (s *State) SetMyDestination(i, destination, distance int) {
	s.Me.timeToDestination = distance
	s.Me.destination = destination
	s.visitedValves.Set(i)
}

func (s *State) SetElephantDestination(i, destination, distance int) {
	s.Elephant.timeToDestination = distance
	s.Elephant.destination = destination
	s.visitedValves.Set(i)
}

func main() {
	adjacency, valveToIndex, flowRates := parseInput()
	distanceMatrix := createDistanceMatrix(adjacency)

	var valves []int
	for k := range flowRates {
		valves = append(valves, k)
	}

	indexToValve := map[int]string{}
	for k, v := range valveToIndex {
		indexToValve[v] = k
	}

	partOneState := State{
		Me: Searcher{
			position:    valveToIndex["AA"],
			destination: -1,
		},
	}

	partOneState = search(partOneState, false, valves, distanceMatrix, flowRates)
	fmt.Println(partOneState)

	partTwoState := State{
		time: 4,
		Me: Searcher{
			position:    valveToIndex["AA"],
			destination: -1,
		},
		Elephant: Searcher{
			position:    valveToIndex["AA"],
			destination: -1,
		},
	}
	partTwoState = search(partTwoState, true, valves, distanceMatrix, flowRates)
	fmt.Println(partTwoState)
}

func search(initialState State, useElephant bool, valves []int, distanceMatrix [][]int, flowRates map[int]int) State {
	bestState := initialState

	stateQueue := lib.Queue[State]{}
	stateQueue.Push(initialState)

	for {
		state, ok := stateQueue.Pop()
		if !ok {
			break
		}

		// if I don't know where to go
		if state.Me.destination == -1 && !state.Me.waiting {
			var branched bool
			for i := range valves {
				// if we've already been to this valve, ignore
				if state.visitedValves.Has(i) {
					continue
				}

				branched = true
				branchState := state

				// set new destination and update travel time
				branchState.SetMyDestination(i, valves[i], distanceMatrix[branchState.Me.position][valves[i]])

				// search that branch of possibility for the best pressure released
				stateQueue.Push(branchState)
			}
			if branched {
				continue
			}
		}

		// if the elephant doesnt know where to go next
		if useElephant && state.Elephant.destination == -1 && !state.Elephant.waiting {
			var branched bool
			for i := range valves {
				// if we've already been to this valve, ignore
				if state.visitedValves.Has(i) {
					continue
				}

				branchState := state
				branched = true

				// set new destination and update travel time
				branchState.SetElephantDestination(i, valves[i], distanceMatrix[branchState.Me.position][valves[i]])

				// search that branch of possibility for the best pressure released
				stateQueue.Push(branchState)
			}
			if branched {
				continue
			}
		}

		if state.Me.destination == -1 {
			state.Me.waiting = true
		}
		if state.Elephant.destination == -1 {
			state.Elephant.waiting = true
		}

		state.Tick()

		if !state.Me.waiting {
			if state.Me.position == state.Me.destination {
				state.Me.destination = -1
				state.flow += flowRates[state.Me.position]
			}

			if state.Me.timeToDestination == 0 {
				state.Me.position = state.Me.destination
			}
		}

		if useElephant && !state.Elephant.waiting {
			if state.Elephant.position == state.Elephant.destination {
				state.Elephant.destination = -1
				state.flow += flowRates[state.Elephant.position]
			}

			if state.Elephant.timeToDestination == 0 {
				state.Elephant.position = state.Elephant.destination
			}
		}

		if state.time < 30 {
			stateQueue.Push(state)
			continue
		}

		if state.totalReleased > bestState.totalReleased {
			bestState = state
		}
	}

	return bestState
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

		// breadth-first search
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

func parseInput() ([][]int, map[string]int, map[int]int) {
	var adjacency [][]int
	nameToValveIndex := map[string]int{}
	flowRates := map[int]int{}

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
			flowRates[valveIndex] = flowRate
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
