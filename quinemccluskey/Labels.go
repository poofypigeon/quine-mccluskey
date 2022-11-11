package quinemccluskey

import (
	"errors"
	"fmt"
)

type InputLabels struct {
	labels [64]string
}

func (l *InputLabels) Init() {
	l.labels = [64]string{}
}

func (l *InputLabels) Set(x int, label string) {
	if x >= len(l.labels) {
		return
	}

	l.labels[x] = label
}

func (l InputLabels) Str(x int) string {
	if l.labels[x] == "" {
		return fmt.Sprintf("x%d", x)
	}

	return l.labels[x]
}

type OutputLabels struct {
	labels   [64]string
	nOutputs int
}

func (l *OutputLabels) Init() {
	l.labels = [64]string{}
	l.nOutputs = 0
}

func (l *OutputLabels) Add(label string) {
	l.labels[l.nOutputs] = label
	l.nOutputs++
}

func (l OutputLabels) NumOutputs() int {
	return l.nOutputs
}

func (l OutputLabels) Str(x int) string {
	if x > l.nOutputs {
		panic(errors.New("quinemccluskey: x exceeds number of defined outputs"))
	}

	if l.labels[x] == "" {
		return fmt.Sprintf("f%d", x)
	}

	return l.labels[x]
}
