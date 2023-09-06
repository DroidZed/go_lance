package services

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/stretchr/testify/assert"
)

const dummyPwd = "mysecrea56sd0a5s415asd415as5qaas]as'asda.da/*-as=-?/\\'sd;asda00..6as5d45-tpassword"

func TestCompareSecureToPlain(t *testing.T) {
	hashedPassword, err := cryptor.HashPassword(dummyPwd)
	if err != nil {
		t.Fatal(err)
	}

	result := cryptor.CompareSecureToPlain(hashedPassword, dummyPwd)

	assert.Equal(t, true, result)
}

func BenchmarkHashPassword(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cryptor.HashPassword(dummyPwd)
	}
}