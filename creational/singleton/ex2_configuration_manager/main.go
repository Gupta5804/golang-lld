package main

import (
	"sync"
)

type ConfigManager struct {
	config map[string]string
	mu sync.RWMutex
}

func (c *ConfigManager) Set(key, value string){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config[key] = value
}

func (c *ConfigManager) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.config[key]
}

var (
	instance *ConfigManager
	once sync.Once
)

func GetInstance() *ConfigManager {
	once.Do(func(){
		instance = &ConfigManager{
			config: make(map[string]string),
		}
	})
	return instance
}