package circuitbreaker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
)

var Client http.Client

var RETRIES = 3

func CallUsingCircuitBreaker(breakerName string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1)

	errors := hystrix.Go(breakerName, func() error {
		req, _ := http.NewRequest(method, url, nil)
		return callWithRetries(req, output)
	}, func(err error) error {
		logrus.Errorf("In fallback function for breaker %v, error: %v", breakerName, err.Error())
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		logrus.Errorf("Circuit state is: %v", circuit.IsOpen())
		return err
	})

	select {
	case out := <-output:
		logrus.Debugf("Call in breaker %v successful", breakerName)
		return out, nil

	case err := <-errors:
		return nil, err
	}
}

func callWithRetries(req *http.Request, output chan []byte) error {
	r := retrier.New(retrier.ConstantBackoff(RETRIES, 100*time.Millisecond), nil)

	attempt := 0

	err := r.Run(func() error {
		attempt++
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			responseBody, err := ioutil.ReadAll(resp.Body)

			if err == nil {
				output <- responseBody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("Status was %v", resp.StatusCode)
		}

		logrus.Errorf("Retrier failed, attempt %v", attempt)
		return err
	})

	return err
}
