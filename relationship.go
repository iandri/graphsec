package main

import (
	"io"
)

type Relationship struct {
	Name string
	get  func() interface{}
}

type Relationships []Relationship

func NewRelationship() *Relationships {
	r := new(Relationships)
	return r
}

func (r *Relationships) Add(relationship Relationship) {
	*r = append(*r, relationship)
}

func (r *Relationships) Traverse(fullOutput bool, output io.Writer) {
	longStdoutlast := "\n\t-> %s \n\t\t- %v\n"
	longStdout := "\n\t-> %s \n\t\t- %v "
	longStdoutFirst := "%s \n - %v "

	shortStdoutlast := "%s\n"
	shortStdout := "%s -> "

	last := len(*r)

	for i, v := range *r {
		if v.get != nil {
			if i == 0 {
				if fullOutput {
					fprintf(output, longStdoutFirst, v.Name, v.get())
				} else {
					fprintf(output, shortStdout, v.Name)
				}
				continue
			}
			if i == last-1 {
				if fullOutput {
					fprintf(output, longStdoutlast, v.Name, v.get())
				} else {
					fprintf(output, shortStdoutlast, v.Name)
				}
			} else {
				if fullOutput {
					fprintf(output, longStdout, v.Name, v.get())
				} else {
					fprintf(output, shortStdout, v.Name)
				}
			}
		} else {
			fprintf(output, v.Name)
		}
	}
}
