package etcd

import (
	"encoding/json"
	"github.com/coreos/etcd/store"
	"io/ioutil"
	"path"
)

func Get(key string) (*[]store.Response, error) {
	logger.Debugf("get %s [%s]", key, client.cluster.Leader)
	resp, err := sendRequest("GET", path.Join("keys", key), "")

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return convertResponse(b)

}

// Get the value of the key from a given machine address
// If the given machine is not available, return error
// Mainly use for testing purpose
func GetFrom(key string, addr string) (*[]store.Response, error) {
	httpPath := createHttpPath(addr, path.Join(version, "keys", key))

	resp, err := client.httpClient.Get(httpPath)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	return convertResponse(b)
}

// convert byte stream to response
func convertResponse(b []byte) (*[]store.Response, error) {

	var results []store.Response
	var result store.Response

	err := json.Unmarshal(b, &result)

	if err != nil {
		err = json.Unmarshal(b, &results)

		if err != nil {
			return nil, err
		}

	} else {
		results = make([]store.Response, 1)
		results[0] = result
	}
	return &results, nil
}
