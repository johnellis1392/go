package main

import (
	"fmt"
	"os"
)

// Locations of MNIST Training & Testing Datasets
const (
	TrainingImages = "./data/train-images-idx3-ubyte.gz"
	TrainingLabels = "./data/train-labels-idx1-ubyte.gz"
	TestingImages  = "./data/t10k-images-idx3-ubyte.gz"
	TestingLabels  = "./data/t10k-labels-idx1-ubyte.gz"

	ImageSetMagicNumber = int32(0x0803)
	LabelSetMagicNumber = int32(0x0801)
)

// Label represents a dataset's training label values.
type Label uint8

// LabelSet is a set of labels for training a classifier.
type LabelSet struct {
	N      int32
	Labels []Label
}

// Pixel is a single pixel value.
type Pixel uint8

// Image is a 2-dimensional array of pixels.
type Image [][]Pixel

// ImageSet represents a set of input images for the classifier.
type ImageSet struct {
	N      int32
	Rows   int32
	Cols   int32
	Images []Image
}

// MNIST represents the full MNIST dataset.
type MNIST struct {
	TrainingImages *ImageSet
	TrainingLabels *LabelSet
	TestingImages  *ImageSet
	TestingLabels  *LabelSet

	Err error
}

type fileloader struct {
	filename string
	data     []byte
	size     int
	pos      int
	err      error
}

func (fl *fileloader) readi32() int32 {
	if fl.err != nil || fl.pos+3+1 >= len(fl.data) {
		return int32(0)
	}

	b := fl.data[fl.pos : fl.pos+3+1]
	var res int32

	res |= int32(b[0]) << 24
	res |= int32(b[1]) << 16
	res |= int32(b[2]) << 8
	res |= int32(b[3]) << 0

	return res
}

func (fl *fileloader) readui8() uint8 {
	if fl.err != nil || fl.pos >= len(fl.data) {
		return uint8(0)
	}

	b := uint8(fl.data[fl.pos])
	fl.pos++

	return b
}

func (fl *fileloader) load(filename string) {
	if fl.err != nil {
		return
	}

	fl.clear()
	if _, fl.err = os.Stat(filename); fl.err != nil {
		return
	}

	var f *os.File
	f, fl.err = os.Open(filename)
	if fl.err != nil {
		return
	}
	defer f.Close()

	fl.size, fl.err = f.Read(fl.data)
	if fl.err != nil {
		return
	}
}

func (fl *fileloader) clear() {
	fl.filename = ""
	fl.data = []byte{}
	fl.size = 0
	fl.pos = 0
}

func newFileloader() *fileloader {
	fl := &fileloader{
		filename: "",
		data:     []byte{},
		size:     0,
		pos:      0,
		err:      nil,
	}
	return fl
}

func (fl *fileloader) loadImageSet(filename string) *ImageSet {
	if fl.err != nil {
		return nil
	}

	fl.load(TrainingImages)
	if fl.err != nil {
		return nil
	}

	magicnum := fl.readi32()
	if magicnum != ImageSetMagicNumber {
		fl.err = fmt.Errorf("illegal value for image set magic number: %v, %q", magicnum, magicnum)
		return nil
	}

	n := fl.readi32()
	rows := fl.readi32()
	cols := fl.readi32()
	images := []Image{}

	for i := 0; i < int(n); i++ {
		image := Image{}
		for j := 0; j < int(rows); j++ {
			for k := 0; k < int(cols); k++ {
				image[j][k] = Pixel(fl.readui8())
			}
		}
		images = append(images, image)
	}

	imageset := &ImageSet{
		N:      n,
		Rows:   rows,
		Cols:   cols,
		Images: images,
	}

	return imageset
}

func (fl *fileloader) loadLabelSet(filename string) *LabelSet {
	if fl.err != nil {
		return nil
	}

	fl.load(filename)
	if fl.err != nil {
		return nil
	}

	magicnum := fl.readi32()
	if magicnum != LabelSetMagicNumber {
		fl.err = fmt.Errorf("illegal value for image set magic number: %v, %q", magicnum, magicnum)
		return nil
	}

	n := fl.readi32()
	var labels []Label
	for i := 0; i < int(n); i++ {
		labels = append(labels, Label(fl.readui8()))
	}

	labelset := &LabelSet{
		N:      n,
		Labels: labels,
	}

	return labelset
}

// MNISTError is a class for representing errors in MNIST loading.
type MNISTError struct {
	err error
}

var _ error = (*MNISTError)(nil)

func (me *MNISTError) Error() string {
	return me.err.Error()
}

// LoadMNIST loads the whole MNIST dataset.
func LoadMNIST() (*MNIST, error) {
	mnist := &MNIST{}
	fl := newFileloader()

	mnist.TrainingImages = fl.loadImageSet(TrainingImages)
	mnist.TrainingLabels = fl.loadLabelSet(TrainingLabels)

	mnist.TestingImages = fl.loadImageSet(TestingImages)
	mnist.TestingLabels = fl.loadLabelSet(TestingLabels)

	if fl.err != nil {
		mnist.Err = fl.err
		return mnist, fl.err
	}

	return mnist, nil
}
