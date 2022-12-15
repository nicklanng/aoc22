package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math"
)

type Sensor struct {
	Position   lib.Vector2i
	KnownRange int
}

//go:embed input
var input []byte

func main() {
	sensors, beacons, min, max := parseInput()

	part1Count := countSpacesThatCannotHaveABeacon(min, max, beacons, sensors, 2000000)
	fmt.Printf("Part 1 output: %d\n", part1Count)

	part2Position := searchForPossibleLocation(sensors, 4000000)
	fmt.Printf("Part 2 output: %v\n", part2Position.X*4000000+part2Position.Y)
}

func countSpacesThatCannotHaveABeacon(min lib.Vector2i, max lib.Vector2i, beacons []lib.Vector2i, sensors []Sensor, yLevel int) int {
	var part1Count int
	for x := min.X; x <= max.X; x++ {
		pos := lib.Vector2i{X: x, Y: yLevel}
		if _, ok := lib.Find(beacons, pos); ok {
			// beacon here
			continue
		}
		for _, s := range sensors {
			if s.Position.ManhattanDistance(pos) <= s.KnownRange {
				part1Count++
				break
			}
		}
	}
	return part1Count
}

func searchForPossibleLocation(sensors []Sensor, max int) lib.Vector2i {
	// for each sensor, find all the adjacent points to the range diamond and search those
	for _, s := range sensors {
		for _, pos := range getDiamondEdgePoints(s.KnownRange+1, s.Position) {
			if pos.X < 0 || pos.Y < 0 || pos.X > max || pos.Y > max {
				continue
			}
			found := true
			for _, s := range sensors {
				if s.Position.ManhattanDistance(pos) <= s.KnownRange {
					found = false
					break
				}
			}
			if found {
				return pos
			}
		}
	}
	return lib.Vector2i{}
}

func getDiamondEdgePoints(radius int, center lib.Vector2i) []lib.Vector2i {
	points := make([]lib.Vector2i, 0, radius*2)

	top := center.Add(lib.Vector2i{Y: -radius})
	bottom := center.Add(lib.Vector2i{Y: radius})

	points = append(points, top)
	points = append(points, bottom)
	for i := 1; i <= radius; i++ {
		points = append(points, top.Add(lib.Vector2i{X: i, Y: i}))
		points = append(points, top.Add(lib.Vector2i{X: -i, Y: i}))
		points = append(points, bottom.Add(lib.Vector2i{X: i, Y: -i}))
		points = append(points, bottom.Add(lib.Vector2i{X: -i, Y: -i}))
	}

	return points
}

func parseInput() ([]Sensor, []lib.Vector2i, lib.Vector2i, lib.Vector2i) {
	var sensors []Sensor
	var beacons []lib.Vector2i
	var min, max = lib.Vector2i{X: math.MaxInt, Y: math.MaxInt}, lib.Vector2i{X: math.MinInt, Y: math.MinInt}

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		// read line
		var sensorPos, beaconPos lib.Vector2i
		fmt.Sscanf(scanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensorPos.X, &sensorPos.Y, &beaconPos.X, &beaconPos.Y)
		knownRange := sensorPos.ManhattanDistance(beaconPos)
		sensors = append(sensors, Sensor{Position: sensorPos, KnownRange: knownRange})
		beacons = append(beacons, beaconPos)
		if sensorPos.X-knownRange < min.X {
			min.X = sensorPos.X - knownRange
		}
		if sensorPos.Y-knownRange < min.Y {
			min.Y = sensorPos.Y - knownRange
		}
		if beaconPos.X < min.X {
			min.X = beaconPos.X
		}
		if beaconPos.Y < min.Y {
			min.Y = beaconPos.Y
		}
		if sensorPos.X+knownRange > max.X {
			max.X = sensorPos.X + knownRange
		}
		if sensorPos.Y+knownRange > max.Y {
			max.Y = sensorPos.Y + knownRange
		}
		if beaconPos.X > max.X {
			max.X = beaconPos.X
		}
		if beaconPos.Y > max.Y {
			max.Y = beaconPos.Y
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return sensors, beacons, min, max
}
