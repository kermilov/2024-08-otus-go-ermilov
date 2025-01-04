package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// go test -bench=BenchmarkGetDomainStat -gcflags="-N" .
func BenchmarkGetDomainStat(b *testing.B) {
	b.Skip("для истории, с помощью этого теста получен файл bench_before.log / bench_after.log")
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		GetDomainStat(data, "biz")
		b.StopTimer()
	}
}
