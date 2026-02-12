Sync Once

-> Rune one-time initialization functions

once.Do(funcValue)

-> sync.Once ensure that only one call to Do ever calls the function passed in - even on different goroutines

(pretty useful for singleton pattern)