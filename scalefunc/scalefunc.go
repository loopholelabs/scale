package scalefunc

import "github.com/loopholelabs/scale-go/scalefile"

type ScaleFunc struct {
	ScaleFile scalefile.ScaleFile
	Function  []byte
}
