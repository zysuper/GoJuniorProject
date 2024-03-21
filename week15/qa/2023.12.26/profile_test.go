package qa

import (
	"github.com/stretchr/testify/require"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestProfile(t *testing.T) {
	go func() {
		for i := 0; i < 10000; i++ {
			t.Log(i)
			time.Sleep(time.Second)
		}
	}()
	err := http.ListenAndServe(":8081", nil)
	require.NoError(t, err)
}
