package main

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/gonum/stat"
	"gonum.org/v1/gonum/floats"
)

func downsampleSimilar(d [][]float64, selectionRate float32, selectionStandardDeviation float64) [][]float64 {
	rand.Seed(time.Now().UnixNano())
	dest := floats2DToDense(d)
	rlen, _ := dest.Dims()
	pcModel := PrincipalComponentTrain(dest) // get trained PC
	transformed := PrincipalComponentTransform(pcModel, dest, 1)
	v := transformed.RawMatrix().Data

	mean, std := stat.MeanStdDev(v, nil)
	keepIndices := []int{}
	for i := 0; i < rlen; i++ {
		if v[i] < mean+selectionStandardDeviation*std && v[i] > mean-selectionStandardDeviation*std {
			if rand.Float32() < selectionRate {
				keepIndices = append(keepIndices, i)
			}
		} else {
			keepIndices = append(keepIndices, i)
		}
	}
	newDest := [][]float64{}
	for ii := range keepIndices {
		pos := keepIndices[ii]
		newDest = append(newDest, d[pos])
	}
	return newDest
}

func findPaths(start []float64, dest [][]float64, numPaths int) {
	rand.Seed(time.Now().UTC().UnixNano())
	newDest := downsampleSimilar(dest, 0.2, 2) //downsample those within 2 STD to 20%
	sort.Slice(newDest, func(i, j int) bool {
		return newDest[i][0] < newDest[j][0]
	})
	subtract2D(start, newDest)
	partitions := optimizerAnglePartition(newDest, numPaths)
	optimizerEvenDistanceForWalks(partitions, newDest)
}

type AnglePartition struct {
	start     float64
	end       float64 //end must always > start. if >360, keep as such
	pointsPos []int
}

func optimizerAnglePartition(points [][]float64, partitionCount int) []AnglePartition {
	// Reduce total distance
	wg := sync.WaitGroup{}
	intervalAngle := fullRad / float64(partitionCount)
	pointsAngle := _getAnglePoints(points)
	partitions := []AnglePartition{}
	for i := 0; i < partitionCount; i++ {
		partitions = append(partitions, AnglePartition{start: float64(i) * intervalAngle, end: (float64(i) + 1) * intervalAngle})
	}
	for i := range partitions {
		wg.Add(1)
		go _getPartitionPoints(&partitions[i], pointsAngle, &wg)
	}
	wg.Wait()
	k := 1.0
	for iii := 0; iii < 100; iii++ { //converges in ~100
		for i := range partitions {
			nextP := i + 1
			if i+1 == len(partitions) {
				nextP = 0
			}
			ptdiff := len(partitions[nextP].pointsPos) - len(partitions[i].pointsPos)
			avchng := 0.0
			if ptdiff > 0 { //eats into next
				nextAvRng := partitions[nextP].end - partitions[nextP].start
				avchng = float64(ptdiff) / 2.0 / float64(len(partitions[nextP].pointsPos)) * nextAvRng

			} else if ptdiff < 0 { //eats into cur
				curAvRng := partitions[i].end - partitions[i].start
				avchng = float64(ptdiff) / 2.0 / float64(len(partitions[i].pointsPos)) * curAvRng
			}
			partitions[nextP].start += k * avchng
			if partitions[nextP].start >= fullRad {
				partitions[nextP].start -= fullRad
				partitions[nextP].end -= fullRad
			} else if partitions[nextP].start < 0 {
				partitions[nextP].start += fullRad
				partitions[nextP].end += fullRad
			}
			partitions[i].end += k * avchng
		}
		k *= 0.99
		for i := range partitions {
			print(len(partitions[i].pointsPos))
			print(",")
		}

		println("iter:", iii)
		for i := range partitions {
			wg.Add(1)
			go _getPartitionPoints(&partitions[i], pointsAngle, &wg)
		}
		wg.Wait()
	}
	return partitions
}

func optimizerEvenDistanceForWalks(adjPartitionsPoints []AnglePartition, allpoints [][]float64) {
	for p := range adjPartitionsPoints {
		points := [][]float64{}
		for i := range adjPartitionsPoints[p].pointsPos {
			points = append(points, allpoints[adjPartitionsPoints[p].pointsPos[i]])
		}
		nnpoints := _findShortestDistanceWalk([]float64{0, 0}, points)
		println("Zero centered nearest neighbour path of worker", p)
		for i := range nnpoints {
			println(nnpoints[i][0], nnpoints[i][1])
		}

		_findShortestDistanceBySwappingNClosestNeighbours([]float64{0, 0}, nnpoints, 2)
	}
}

func _findShortestDistanceWalk(startpt []float64, points [][]float64) [][]float64 {
	// this assumes walk from 0,0
	// double index, sort 1st index, find nearest neigh
	// branch-and-bound
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})
	sort.Slice(points, func(i, j int) bool {
		return points[i][0] < points[j][0]
	})

	start := [][]float64{startpt}
	points = append(start, points...)
	path := make([]int, len(points))
	used := make([]bool, len(points))

	used[0] = true
	for ii := 0; ii < len(points)-1; ii++ {
		curPathPos := ii
		curNode := path[curPathPos]
		shortestDist := -1.0

		for i := 1; i < len(points); i++ {
			if !used[i] {
				if shortestDist < 0 {
					shortestDist = _distL2(points[curNode], points[i])
					path[curPathPos+1] = i
				} else if math.Abs(points[curNode][0]-points[i][0]) < shortestDist {
					d := _distL2(points[curNode], points[i])
					if d < shortestDist {
						path[curPathPos+1] = i
						shortestDist = d
					}
				} else if points[i][0]-points[curNode][0] > shortestDist { //first index already > current shortest dist, ignore rest
					break
				}
			}
		}
		used[path[curPathPos+1]] = true
	}
	pathNodes := [][]float64{}
	for i := range path {
		pathNodes = append(pathNodes, points[path[i]])
	}
	return pathNodes[1:]
}

func _findShortestDistanceBySwappingNClosestNeighbours(startpt []float64, points [][]float64, nndistMultiplier float64) {
	// use after _findShortestDistanceWalk
	start := [][]float64{startpt}
	points = append(start, points...)
	for i := 0; i < len(points)-2; i++ {
		nei := []int{}
		rdist := _distL2(points[i+1], points[i])
		for j := i + 2; j < len(points); j++ {
			if _distL2(points[j], points[i]) < nndistMultiplier*rdist {
				nei = append(nei, j)
			}
		}
		// found possible good neighs
	}
}

func _getTotalTravelledDistance(nodes [][]float64) float64 {
	s := 0.0
	for i := 0; i < len(nodes)-1; i++ {
		s += _distL2(nodes[i], nodes[i+1])
	}
	return s
}

func _distL2(a, b []float64) float64 {
	v := []float64{0, 0}
	floats.SubTo(v, a, b)
	dist := _l2(v)
	return dist
}

func _getAnglePoints(allpoints [][]float64) []float64 {
	ref := []float64{1, 0}
	angles := []float64{}
	for i := range allpoints {
		av := GetAngle(ref, allpoints[i])
		if allpoints[i][1] < 0 {
			av += math.Pi
		}
		angles = append(angles, av)
	}
	return angles
}

func _getPartitionPoints(partition *AnglePartition, allangles []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	partition.pointsPos = make([]int, 0)
	for i := range allangles {
		av := allangles[i]
		if av >= partition.start && av < partition.end {
			partition.pointsPos = append(partition.pointsPos, i)
		} else if partition.end >= fullRad {
			if av >= partition.start || av < partition.end-fullRad {
				partition.pointsPos = append(partition.pointsPos, i)
			}
		}

	}
}
