### sync.Pool

-> create and make available pool of things for use

b := bufPool.Get().(*bytes.Buffer)

bufPool.Put(b)
