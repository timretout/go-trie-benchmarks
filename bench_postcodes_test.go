package main

import (
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"

	badgerodon "github.com/badgerodon/collections/trie"
	claudiu "github.com/claudiu/trie"
	derekparker "github.com/derekparker/trie"
	dghubble "github.com/dghubble/trie"
	timretout "github.com/timretout/trie"
	viant "github.com/viant/ptrie"
)

var postcodes []string

func init() {
	onspd := "../ONSPD/Data/ONSPD_MAY_2020_UK.csv"

	file, err := os.Open(onspd)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(file)
	r.ReuseRecord = true

	// Skip header
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		postcode := strings.ReplaceAll(record[0], " ", "")
		postcodes = append(postcodes, postcode)
	}

}

var tries = []struct {
	Name   string
	New    func() interface{}
	Insert func(interface{}, string)
	Exists func(interface{}, string) bool
}{
	{"badgerodon",
		func() interface{} { return badgerodon.New() },
		func(tr interface{}, s string) { tr.(*badgerodon.Trie).Insert(s, struct{}{}) },
		func(tr interface{}, s string) bool { return tr.(*badgerodon.Trie).Has(s) },
	},
	{"claudiu",
		func() interface{} { return claudiu.NewTrie() },
		func(tr interface{}, s string) { tr.(*claudiu.Trie).Add(s) },
		func(tr interface{}, s string) bool { return tr.(*claudiu.Trie).Find(s) != nil },
	},
	{"derekparker",
		func() interface{} { return derekparker.New() },
		func(tr interface{}, s string) { tr.(*derekparker.Trie).Add(s, struct{}{}) },
		func(tr interface{}, s string) bool { _, ok := tr.(*derekparker.Trie).Find(s); return ok },
	},
	{"dghubble",
		func() interface{} { return dghubble.NewRuneTrie() },
		func(tr interface{}, s string) { tr.(*dghubble.RuneTrie).Put(s, struct{}{}) },
		func(tr interface{}, s string) bool { return tr.(*dghubble.RuneTrie).Get(s) != nil },
	},
	{"timretout",
		func() interface{} { return timretout.New() },
		func(tr interface{}, s string) { tr.(*timretout.Trie).Insert(s) },
		func(tr interface{}, s string) bool { return tr.(*timretout.Trie).Exists(s) },
	},
	{"viant",
		func() interface{} { return viant.New() },
		func(tr interface{}, s string) { tr.(viant.Trie).Put([]byte(s), "") },
		func(tr interface{}, s string) bool { return tr.(viant.Trie).Has([]byte(s)) },
	},
}

func BenchmarkImportONSPD(b *testing.B) {
	for _, bm := range tries {
		b.Run(bm.Name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				tr := bm.New()
				for _, v := range postcodes {
					bm.Insert(tr, v)
				}
			}
		})
	}
}

func BenchmarkONSPDSequentialExists(b *testing.B) {
	for _, bm := range tries {
		b.Run(bm.Name, func(b *testing.B) {
			tr := bm.New()
			for _, v := range postcodes {
				bm.Insert(tr, v)
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				if !bm.Exists(tr, postcodes[n%len(postcodes)]) {
					b.Error("something went wrong")
				}
			}
		})
	}
}

func BenchmarkONSPDRandomExists(b *testing.B) {
	for _, bm := range tries {
		b.Run(bm.Name, func(b *testing.B) {
			tr := bm.New()
			for _, v := range postcodes {
				bm.Insert(tr, v)
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				v := postcodes[rand.Intn(len(postcodes))]
				if !bm.Exists(tr, v) {
					b.Error("something went wrong")
				}
			}
		})
	}
}
