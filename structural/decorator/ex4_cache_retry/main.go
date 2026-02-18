package main

type DataFetcher func(key string) (string, error)

func WithCache(d DataFetcher, c map[string]string) DataFetcher {
	return func(key string) (string, error) {
		if _, ok := c[key]; ok {
			return c[key], nil
		} else {
			return d(key)
		}
	}
}

func WithRetry(d DataFetcher, maxRetries int) DataFetcher {
	return func(key string) (string, error) {
		var lastErr error
		for i := 0; i < maxRetries; i++ {
			val, err := d(key)
			if err == nil {
				return val, nil
			}
			lastErr = err
		}
		return "", lastErr
	}
}
