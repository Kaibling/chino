package utils

import (
	"bytes"
	"chino/models"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

func Request(method, url string, body []byte) ([]byte, error) {
	data := bytes.NewBuffer(body)
	r, err := http.NewRequest(method, url, data)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	rbody, _ := ioutil.ReadAll(response.Body)
	return rbody, nil
}

func GetContext(key string, r *http.Request) interface{} {
	parameter := r.Context().Value(models.String(key))
	if parameter == nil {
		panic("context parameter '" + key + "' missing")
	}
	return parameter
}

func NewULID() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}
