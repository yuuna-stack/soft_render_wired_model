package tgalib

import (
	"encoding/binary"
	"errors"
	"os"
)

type Format int

type RGBColor struct {
	data [3]byte
}

func NewRGB(r byte, g byte, b byte) *RGBColor {
	return &RGBColor{[3]byte{b, g, r}}
}

type tga_header struct {
	idlength        byte
	colormaptype    byte
	datatypecode    byte
	colormaporigin  int16
	colormaplength  int16
	colormapdepth   byte
	x_origin        int16
	y_origin        int16
	width           int16
	height          int16
	bitsperpixel    byte
	imagedescriptor byte
}

type TGAImage struct {
	data   []byte
	width  int
	height int
}

func NewTGAImage(width int, height int) *TGAImage {
	data := make([]byte, width*height*3)
	return &TGAImage{data, width, height}
}

func SetRGBColor(img *TGAImage, x int, y int, c *RGBColor) error {
	if img.data == nil || x < 0 || y < 0 || x >= img.width || y >= img.height {
		return errors.New("Wrong coordinate")
	}
	copy(img.data[(x+y*img.width)*3:(x+y*img.width)*3+3], c.data[:])
	return nil
}

func (img *TGAImage) WriteTgaFile(name string, rle bool) error {
	file, err := os.Create(name)
	defer file.Close()
	if err != nil {
		return errors.New("Can't create file")
	}
	header := tga_header{}
	header.bitsperpixel = 3 << 3
	header.width = int16(img.width)
	header.height = int16(img.height)
	typecode := 0
	if rle {
		typecode = 10
	} else {
		typecode = 2
	}
	header.datatypecode = byte(typecode)
	header.imagedescriptor = 0x20 // top-left origin
	err = binary.Write(file, binary.LittleEndian, header)
	if err != nil {
		return errors.New("Can't dump the tga file")
	}
	if rle {
		err = img.writeRleData(file)
		if err != nil {
			return errors.New("Can't write raw data")
		}
	} else {
		err = binary.Write(file, binary.LittleEndian, img.data)
		if err != nil {
			return errors.New("Can't write raw data")
		}
	}
	dev_area_ref := [4]byte{0, 0, 0, 0}
	err = binary.Write(file, binary.LittleEndian, dev_area_ref)
	if err != nil {
		return errors.New("Can't dump the tga file")
	}
	ext_area_ref := [4]byte{0, 0, 0, 0}
	err = binary.Write(file, binary.LittleEndian, ext_area_ref)
	if err != nil {
		return errors.New("Can't dump the tga file")
	}
	footer := [18]byte{0x54, 0x52, 0x55, 0x45, 0x56, 0x49, 0x53, 0x49, 0x4F, 0x4E, 0x2D, 0x58, 0x46, 0x49, 0x4C, 0x45, 0x2E, 0x00}
	err = binary.Write(file, binary.LittleEndian, footer)
	if err != nil {
		return errors.New("Can't dump the tga file")
	}
	return nil
}

func (img *TGAImage) writeRleData(file *os.File) error {
	var max_chunk_length byte = 128
	var npixels uint64 = uint64(img.width) * uint64(img.height)
	var curpix uint64 = 0
	for curpix < npixels {
		chankstart := curpix * 3
		curbyte := curpix * 3
		var run_length byte = 1
		raw := true
		for curpix+uint64(run_length) < npixels && run_length < max_chunk_length {
			succ_eq := true
			for t := 0; succ_eq && t < 3; t++ {
				succ_eq = img.data[curbyte+uint64(t)] == img.data[curbyte+uint64(t+3)]
			}
			curbyte += 3
			if run_length == 1 {
				raw = !succ_eq
			}
			if raw && succ_eq {
				run_length--
				break
			}
			if !raw && !succ_eq {
				break
			}
			run_length++
		}
		curpix += uint64(run_length)
		if raw {
			err := binary.Write(file, binary.LittleEndian, run_length-1)
			if err != nil {
				return errors.New("Can't dump the tga file")
			}
		} else {
			err := binary.Write(file, binary.LittleEndian, run_length+127)
			if err != nil {
				return errors.New("Can't dump the tga file")
			}
		}
		if raw {
			err := binary.Write(file, binary.LittleEndian, img.data[chankstart:chankstart+uint64(run_length*3)])
			if err != nil {
				return errors.New("Can't dump the tga file")
			}
		} else {
			err := binary.Write(file, binary.LittleEndian, img.data[chankstart:chankstart+3])
			if err != nil {
				return errors.New("Can't dump the tga file")
			}
		}
	}
	return nil
}

func (img *TGAImage) FlipVertically() bool {
	if img.data == nil {
		return false
	}
	bytes_per_line := img.width * 3
	line := make([]byte, bytes_per_line)
	half := img.height >> 1
	for j := 0; j < half; j++ {
		l1 := j * bytes_per_line
		l2 := (img.height - 1 - j) * bytes_per_line
		copy(line, img.data[l1:l1+bytes_per_line])
		copy(img.data[l1:], img.data[l2:l2+bytes_per_line])
		copy(img.data[l2:], line[:bytes_per_line])
	}
	return true
}
