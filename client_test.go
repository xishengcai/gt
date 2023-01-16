package gt

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	go mockServer(t)

	result := personExample{}

	url := "http://localhost/json-example"
	err := NewDefaultClient().
		SetURL(url).
		Get().
		Do().
		InTo(&result, JSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", result)
}

func TestHttpLog(t *testing.T) {
	go mockServer(t)

	for i := 0; i < 100; i++ {
		go func() {
			result := personExample{}
			url := "http://localhost/json-example"
			err := NewDefaultClient().
				SetURL(url).
				SetLog("hello").
				Get().
				Do().
				InTo(&result, JSON)
			if err != nil {
				t.Fatal(err)
			}
			//t.Logf("%+v", result)
		}()
	}

	time.Sleep(10 * time.Second)

}
