#  Test-ddos

## lib simple attack fetch ddos

```go
	var origins = []string{}
			for i := 0; i < workers; i++ {
				origins = append(origins, fmt.Sprintf("http://example%v.com", i))
			}
			d, err := ddos.New(url, workers, method, "", origins)
			if err != nil {
				panic(err)
			}
			d.Run()
			d.Result()
			fmt.Printf("\nDDoS attack server:%s", url)
```