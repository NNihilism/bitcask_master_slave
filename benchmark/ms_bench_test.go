package benchmark

import (
	"fmt"
	"testing"
	"time"

	msClient "bitcask_master_slave/client/ms_client"

	"github.com/stretchr/testify/assert"
)

var cli *msClient.Client

func init() {
	cli = msClient.NewClient(&msClient.MSClientConfig{
		MasterHost: "127.0.0.1",
		MasterPort: "8992",
	})
	time.Sleep(time.Second * 10)
	initDataForGet()
}

func initDataForGet() {
	writeCount := 30000
	for i := 0; i < writeCount; i++ {
		_, err := cli.Set(getKeyAndValue(i))
		if err != nil {
			panic(err)
		}
		fmt.Println(i)
	}
	fmt.Println("After init.")
}

func BenchmarkMS_Get(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := cli.Get([][]byte{key2[i]})
		assert.Nil(b, err)
	}
}
