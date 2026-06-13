package main

import "sync"

var store = map[string]string{}
var storeMu sync.RWMutex
