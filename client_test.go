package gt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	go mockServer(t)

	result := personExample{}

	url := "http://localhost/json-example"
	err := NewDefaultClient().
		GET(url).
		Do().
		InTo(&result, JSON)
	assert.Nil(t, err)
	t.Logf("%+v", result)
}

func TestHttpLog(t *testing.T) {
	go mockServer(t)

	for i := 0; i < 100; i++ {
		go func() {
			result := personExample{}
			url := "http://localhost/json-example"
			err := NewDefaultClient().
				EnableLog().
				GET(url).
				Do().
				InTo(&result, JSON)
			assert.Nil(t, err)
		}()
	}

	time.Sleep(10 * time.Second)

}
