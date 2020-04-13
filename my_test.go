package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestPrincipleComponent(t *testing.T) {
	xx := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	xx1 := mat.NewDense(3, 5, xx)

	pcModel := PrincipalComponentTrain(xx1) // get trained PC
	transformed := PrincipalComponentTransform(pcModel, xx1, 3)
	// NOTE! The signs are flipped against python scikit module
	if !floats.EqualWithinAbs(floats.Sum(floats.SubTo(make([]float64, 3), mat.Row(nil, 0, transformed), []float64{-11.180339887498949, 0, 0})), 0.001, 0.001) {
		panic("PCA values wrong")
	}
	if !floats.EqualWithinAbs(floats.Sum(floats.SubTo(make([]float64, 3), mat.Row(nil, 1, transformed), []float64{0, 0, 0})), 0.001, 0.001) {
		panic("PCA values wrong")
	}
	if !floats.EqualWithinAbs(floats.Sum(floats.SubTo(make([]float64, 3), mat.Row(nil, 2, transformed), []float64{11.180339887498949, 0, 0})), 0.001, 0.001) {
		panic("PCA values wrong")
	}

	transformed = PrincipalComponentTransform(pcModel, mat.NewDense(1, 5, []float64{1, 2, 3, 4, 5}), 3)
	if !floats.EqualWithinAbs(floats.Sum(floats.SubTo(make([]float64, 3), mat.Row(nil, 0, transformed), []float64{-11.180339887498949, 0, 0})), 0.001, 0.001) {
		panic("PCA values wrong")
	}
	//panic("Transform new not done!")
}

func TestGetAngle(t *testing.T) {
	c := func(x float64) float64 { return x / 2 / math.Pi * 360 }
	v := func(a, b, lim float64) bool {
		if math.Abs(a-b) < lim {
			return true
		}
		return false
	}
	t.Run("45 degrees", func(t *testing.T) {
		if a := c(GetAngle([]float64{1, 0}, []float64{1, 1})); !v(a, 45, 0.1) {
			t.Errorf("Wrong angle. Got " + fmt.Sprintf("%f", a))
		}
	})

	t.Run("90 degrees", func(t *testing.T) {
		if a := c(GetAngle([]float64{1, 0}, []float64{0, 1})); !v(a, 90, 0.1) {
			t.Errorf("Wrong angle. Got " + fmt.Sprintf("%f", a))
		}
	})

	t.Run("135 degrees", func(t *testing.T) {
		if a := c(GetAngle([]float64{0, 1}, []float64{1, -1})); !v(a, 135, 0.1) {
			t.Errorf("Wrong angle. Got " + fmt.Sprintf("%f", a))
		}
	})

	t.Run("360 degrees", func(t *testing.T) {
		if a := c(GetAngle([]float64{1, 0}, []float64{1, 0})); !v(a, 0, 0.1) {
			t.Errorf("Wrong angle. Got " + fmt.Sprintf("%f", a))
		}
	})
}

func TestPartitionPoints(t *testing.T) {
	partitionCount := 3

	intervalAngle := fullRad / float64(partitionCount)
	dest := parse2DStringAs2DFloats(readCSV("cleaned.csv"))
	newDest := downsampleSimilar(dest, 0.2, 2) //downsample those within 2 STD to 20%
	pointsAngle := _getAnglePoints(newDest)

	partitions := []AnglePartition{}

	for i := 0; i < partitionCount; i++ {
		partitions = append(partitions, AnglePartition{start: float64(i) * intervalAngle, end: (float64(i) + 1) * intervalAngle})
	}
	for i := range partitions {
		wg.Add(1)
		go _getPartitionPoints(&partitions[i], pointsAngle, &wg)
	}
	wg.Wait()
	s := 0
	for i := range partitions {
		s += len(partitions[i].pointsPos)
	}
	if s != len(newDest) {
		t.Errorf("Data is wrongly partitioned.")
	}
}

func TestWalkNearest(t *testing.T) {
	pos := [][]float64{
		[]float64{0, 1},
		[]float64{0, 10},
		[]float64{1, 0.5},
		[]float64{11, 0},
	}
	paths2 := _findShortestDistanceWalk([]float64{0, 0}, pos)
	paths1 := append([][]float64{[]float64{0, 0}}, pos...)
	d1 := _getTotalTravelledDistance(paths1)
	d2 := _getTotalTravelledDistance(paths2)
	if d2 >= d1 {
		t.Errorf("Optimized distance longer than original distance")
	}
}
