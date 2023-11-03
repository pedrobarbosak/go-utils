package main

import (
	"log"

	"github.com/pedrobarbosak/go-utils/defaults"
)

type MyStruct struct {
	Name        string  `default:"John Smith"`
	Age         uint    `default:"42"`
	Float       float32 `default:"42.55"`
	Ignore      string  `default:"-"`
	IgnoreEmpty string  `default:""`
	Struct2     Struct2 `default:""`
	Slice       []*Struct2
	Slice2      []Struct2 `default:""`
}

type Struct2 struct {
	Name string `default:"ABC"`
}

func (s *Struct2) String() string {
	return s.Name
}

func main() {
	defaulter := defaults.New()

	s1 := MyStruct{
		Name:        "",
		Age:         0,
		Float:       0,
		Ignore:      "A",
		IgnoreEmpty: "",
		Struct2:     Struct2{Name: "F"},
		Slice:       []*Struct2{{Name: "AAA"}, {}},
		Slice2:      []Struct2{{Name: "RRR"}, {}},
	}

	err := defaulter.SetDefaults(&s1)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", s1)

	s2 := MyStruct{
		Name:        "AAA",
		Age:         1,
		Float:       2,
		Ignore:      "",
		IgnoreEmpty: "FF",
		Slice:       []*Struct2{{}, {}},
	}

	err = defaulter.SetDefaults(&s2)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", s2)
	log.Println("SLice2 is nil:", s2.Slice2 == nil)
}
