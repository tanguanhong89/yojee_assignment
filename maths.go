package main

import (
	"math"
	"math/rand"
	"reflect"
	"strings"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

var fullRad = 2 * math.Pi

type PrincipalComponent struct {
	pc     stat.PC
	offset []float64
}

// GetDenseFrom2DFloatArr returns 2d dense matrix, gen as generator function
func GetDenseFrom2DFloatArr(arr [][]float64) *mat.Dense {
	row := len(arr)
	col := len(arr[0])
	var arr2 []float64
	for i := range arr {
		arr2 = append(arr2, arr[i]...)
	}
	return mat.NewDense(row, col, arr2)
}

// GetDense2D returns 2d dense matrix, gen as generator function
func GetDense2D(x, y int, gen func(x, y int) float64) *mat.Dense {
	data := make([]float64, x*y)
	for i := range data {
		data[i] = gen(i/y, i/x)
	}
	return mat.NewDense(x, y, data)
}

// GetRandomDense2D returns 2d dense matrix
func GetRandomDense2D(x, y int) *mat.Dense {
	return GetDense2D(x, y, func(x, y int) float64 { return rand.Float64() })
}

// VectorizeString returns an arr of *mat.Dense
func VectorizeString(toBeVectorized string, letters string, vecLetters *mat.Dense) *mat.Dense {
	var arr []float64
	for le := range toBeVectorized {
		l := toBeVectorized[le : le+1]
		pos := strings.Index(letters, l)
		if pos < 0 {
			pos = 0
			println("New char found:", l)
		}
		arr = append(arr, mat.Row(nil, pos, vecLetters)...)
	}
	_, y := vecLetters.Dims()
	return mat.NewDense(len([]rune(toBeVectorized)), y, arr)
}

// FastSigmoid x / (1 + math.Abs(x))
func FastSigmoid(x float64) float64 {
	return x / (1 + math.Abs(x))
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// DenseMapFunction value mapper, eg DenseMapFunction(seed, FastSigmoid)
func DenseMapFunction(inMat *mat.Dense, boo func(float64) float64) {
	foo := func(i, j int, v float64) float64 {
		return boo(inMat.At(i, j))
	}
	inMat.Apply(foo, inMat)
}

func StackDense(darr []*mat.Dense) *mat.Dense {
	_, c := darr[0].Dims()
	base := mat.NewDense(len(darr), c, nil)
	for i := range darr {
		base.SetRow(i, darr[i].RawRowView(0))
	}
	return base
}

func GetFloatArrFromDense(myDense mat.Matrix) [][]float64 {
	var r [][]float64
	row, _ := myDense.Dims()
	for i := 0; i < row; i++ {
		r = append(r, mat.Row(nil, i, myDense))
	}
	return r
}

// PrincipalComponentTrain returns a trained principal component model
func PrincipalComponentTrain(strDenseVec mat.Matrix) PrincipalComponent {
	modelPrincipleComponentMain := stat.PC{}
	if status := modelPrincipleComponentMain.PrincipalComponents(strDenseVec, nil); !status {
		panic(status)
	}

	row, col := strDenseVec.Dims()
	var offset = make([]float64, col)

	for ii := 0; ii < col; ii++ {
		offset[ii] = floats.Sum(mat.Col(nil, ii, strDenseVec)) / float64(row)
	}
	return PrincipalComponent{pc: modelPrincipleComponentMain, offset: offset}
}

func PrincipalComponentGetNComp(trainedPC PrincipalComponent, noOfComp int) *mat.Dense {
	v := reflect.ValueOf(trainedPC.pc)
	n := int(v.FieldByName("n").Int())
	d := int(v.FieldByName("d").Int()) // col
	_r := n
	if d < n {
		_r = d
	}

	pcVec := mat.NewDense(d, _r, nil)
	trainedPC.pc.VectorsTo(pcVec)
	if noOfComp == 0 {
		return pcVec
	}
	return pcVec.Slice(0, d, 0, noOfComp).(*mat.Dense)
}

// PrincipalComponentTransform is same as sklearn.decomposition.PCA.transform
func PrincipalComponentTransform(trainedPC PrincipalComponent, inMat *mat.Dense, ncom int) *mat.Dense {
	row, col := inMat.Dims()
	inMat2 := mat.DenseCopyOf(inMat)
	pcVec := PrincipalComponentGetNComp(trainedPC, 0)

	explainedVar := trainedPC.pc.VarsTo(nil)[:ncom]
	floats.ScaleTo(explainedVar, 1/floats.Sum(explainedVar), explainedVar)

	var offset = trainedPC.offset

	for j := 0; j < row; j++ {
		_r := inMat2.RawRowView(j)
		floats.Sub(_r, offset)
		inMat2.SetRow(j, _r)
	}

	result := mat.NewDense(row, ncom, nil)
	pcVec2 := pcVec.Slice(0, col, 0, ncom)
	result.Mul(inMat2, pcVec2)
	for i := 0; i < row; i++ {
		_r := result.RawRowView(i)
		floats.Mul(_r, explainedVar)
		result.SetRow(i, _r)
	}
	return result
}

func NormalizeCol(dat *mat.Dense, maxMin [][]float64) {
	_, col := dat.Dims()
	for i := 0; i < col; i++ {
		maxMinC := maxMin[i]
		rng := maxMinC[0] - maxMinC[1]
		c := mat.Col(nil, i, dat)
		floats.AddConst(-maxMinC[1], c)
		floats.Scale(1.0/rng, c)
		dat.SetCol(i, c)
	}
}
func GetAngle(a []float64, b []float64) float64 {
	x := floats.Dot(a, b)
	cosTheta := x / _l2(a) / _l2(b)
	if cosTheta > 1.0 { // round off errors may result in this
		cosTheta = 1.0
	} else if cosTheta < -1.0 {
		cosTheta = -1.0
	}
	return math.Acos(cosTheta)
}

func _l2(s []float64) float64 {
	return floats.Norm(s, 2)
}

func subtract2D(a []float64, b [][]float64) {
	for i := range b {
		floats.SubTo(b[i], b[i], a)
	}
}

func floats2DToDense(a [][]float64) *mat.Dense {
	rlen, clen := len(a), len(a[0])
	d := mat.NewDense(rlen, clen, make([]float64, rlen*clen))
	for r := range a {
		d.SetRow(r, a[r])
	}
	return d
}
