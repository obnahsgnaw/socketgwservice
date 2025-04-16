package main

import (
	"github.com/obnahsgnaw/pbhttp/pkg/modelfield"
	"github.com/obnahsgnaw/socketgwservice/internal/dal/field"
)

func main() {
	modelfield.Gen("./internal/dal/model", field.Models()...)
}
