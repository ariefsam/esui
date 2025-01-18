package idgenerator

import "github.com/teris-io/shortid"

func Generate() string {
	return shortid.MustGenerate()
}
