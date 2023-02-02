package model

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Vector3f struct {
	data [3]float32
}

func NewVector3f(x float32, y float32, z float32) *Vector3f {
	return &Vector3f{[3]float32{x, y, z}}
}

func (v *Vector3f) GetX() float32 {
	return v.data[0]
}

func (v *Vector3f) GetY() float32 {
	return v.data[1]
}

func (v *Vector3f) GetZ() float32 {
	return v.data[2]
}

type Face struct {
	idx []int
}

type Model struct {
	verts []Vector3f
	faces []Face
}

func ReadModel(filename string) (*Model, error) {
	readFile, err := os.Open(filename)
	defer readFile.Close()
	if err != nil {
		return nil, errors.New("Can't open file")
	}
	m := Model{}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		s := fileScanner.Text()
		if strings.HasPrefix(s, "v ") {
			strs := strings.Split(s, " ")
			xs, ys, zs := strs[1], strs[2], strs[3]
			x, err := strconv.ParseFloat(xs, 32)
			if err != nil {
				x = 0
			}
			y, err := strconv.ParseFloat(ys, 32)
			if err != nil {
				y = 0
			}
			z, err := strconv.ParseFloat(zs, 32)
			if err != nil {
				z = 0
			}
			m.verts = append(m.verts, *NewVector3f(float32(x), float32(y), float32(z)))
		} else if strings.HasPrefix(s, "f ") {
			strs := strings.Split(s, " ")
			f1 := strings.Split(strs[1], "/")
			f2 := strings.Split(strs[2], "/")
			f3 := strings.Split(strs[3], "/")
			idx1, err := strconv.ParseInt(f1[0], 10, 32)
			if err != nil {
				idx1 = 0
			}
			idx1--
			idx2, err := strconv.ParseInt(f2[0], 10, 32)
			if err != nil {
				idx2 = 0
			}
			idx2--
			idx3, err := strconv.ParseInt(f3[0], 10, 32)
			if err != nil {
				idx3 = 0
			}
			idx3--
			f := []int{int(idx1), int(idx2), int(idx3)}
			m.faces = append(m.faces, Face{f})
		}
	}
	readFile.Close()
	return &m, nil
}

func (model *Model) VertexCount() int {
	return len(model.verts)
}

func (model *Model) FacesCount() int {
	return len(model.faces)
}

func (model *Model) GetVertex(idx int) *Vector3f {
	return &model.verts[idx]
}

func (model *Model) GetFace(idx int) *[]int {
	return &model.faces[idx].idx
}
