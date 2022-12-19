package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
)

type State struct {
	Minute int

	Ore      int
	Clay     int
	Obsidian int
	Geode    int

	OreBot      int
	ClayBot     int
	ObsidianBot int
	GeodeRobot  int
}

func (s *State) Tick() {
	s.Minute++
	s.Ore += s.OreBot
	s.Clay += s.ClayBot
	s.Obsidian += s.ObsidianBot
	s.Geode += s.GeodeRobot
}

type Cost struct {
	Ore      int
	Clay     int
	Obsidian int
}

type Blueprint struct {
	ID int

	OreRobot    Cost
	ClayRobot   Cost
	ObsidianBot Cost
	GeodeBot    Cost

	MaxOreBots      int
	MaxClayBots     int
	MaxObsidianBots int
}

//go:embed input
var input []byte

func main() {
	blueprints := parseInput()

	partOne := TestBlueprints(blueprints, 24)
	var total int
	for i := range partOne {
		total += (i + 1) * partOne[i]
	}
	fmt.Println(total)

	partTwo := TestBlueprints(blueprints[:3], 32)
	fmt.Println(partTwo[0] * partTwo[1] * partTwo[2])
}

func TestBlueprints(blueprints []Blueprint, runs int) []int {
	results := make([]int, len(blueprints))

	for i, b := range blueprints {
		var mostGeodes int
		seenStates := map[State]struct{}{}
		states := lib.Stack[State]{}
		states.Push(State{Minute: 0, OreBot: 1})

		for {
			state, ok := states.Pop()
			if !ok {
				break
			}

			if _, ok := seenStates[state]; ok {
				continue
			}
			seenStates[state] = struct{}{}

			if state.Minute == runs {
				//fmt.Println(state)
				if mostGeodes < state.Geode {
					mostGeodes = state.Geode
				}
				continue
			}

			// build geode robot
			if state.Ore >= b.GeodeBot.Ore && state.Obsidian >= b.GeodeBot.Obsidian {
				buildGeodeRobot(&state, b)
				states.Push(state)
				continue
			}

			// build obsidian robot
			if state.Ore >= b.ObsidianBot.Ore && state.Clay >= b.ObsidianBot.Clay && state.ObsidianBot < b.MaxObsidianBots {
				nextState := state
				buildObsidianRobot(&nextState, b)
				states.Push(nextState)
			}

			// build clay robot
			if state.Ore >= b.ClayRobot.Ore && state.ClayBot < b.MaxClayBots {
				nextState := state
				buildClayRobot(&nextState, b)
				states.Push(nextState)
			}

			// build ore robot
			if state.Ore >= b.OreRobot.Ore && state.OreBot < b.MaxOreBots {
				nextState := state
				buildOreRobot(&nextState, b)
				states.Push(nextState)
			}

			// wait
			state.Tick()
			states.Push(state)
		}

		results[i] = mostGeodes
	}
	return results
}

func buildGeodeRobot(stores *State, b Blueprint) {
	stores.Ore -= b.GeodeBot.Ore
	stores.Obsidian -= b.GeodeBot.Obsidian
	stores.Tick()
	stores.GeodeRobot++
}
func buildObsidianRobot(stores *State, b Blueprint) {
	stores.Ore -= b.ObsidianBot.Ore
	stores.Clay -= b.ObsidianBot.Clay
	stores.Tick()
	stores.ObsidianBot++
}
func buildClayRobot(stores *State, b Blueprint) {
	stores.Ore -= b.ClayRobot.Ore
	stores.Tick()
	stores.ClayBot++
}
func buildOreRobot(stores *State, b Blueprint) {
	stores.Ore -= b.OreRobot.Ore
	stores.Tick()
	stores.OreBot++
}

func parseInput() []Blueprint {
	var blueprints []Blueprint

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		var blueprint Blueprint
		fmt.Sscanf(
			scanner.Text(),
			"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
			&blueprint.ID,
			&blueprint.OreRobot.Ore,
			&blueprint.ClayRobot.Ore,
			&blueprint.ObsidianBot.Ore,
			&blueprint.ObsidianBot.Clay,
			&blueprint.GeodeBot.Ore,
			&blueprint.GeodeBot.Obsidian,
		)

		blueprint.MaxOreBots = lib.IntMax(blueprint.OreRobot.Ore, blueprint.ClayRobot.Ore, blueprint.ObsidianBot.Ore, blueprint.GeodeBot.Ore)
		blueprint.MaxClayBots = blueprint.ObsidianBot.Clay
		blueprint.MaxObsidianBots = blueprint.GeodeBot.Obsidian

		blueprints = append(blueprints, blueprint)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return blueprints
}
