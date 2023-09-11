package hw10programoptimization

import (
	"archive/zip"
	"log"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Close()

	data, err := r.File[0].Open()
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}
