## Fan In, Fan Out

-> sometimes in pipeline a stage become computationally extensive and take time to compute

can we break computationally intensive stage into multiple goroutines and run them in parallel to speed it up ?

What is fan-out ?

-> Multiple goroutines are started to read data from the single channel

-> Distribute work amongst a group of workers goroutines to parallelize the CPU usage and the I/O usage

-> Helps computationally extensive stage to run faster

What is Fan-in ?

-> Process of combining multiple results into one channel

-> we create merge goroutines, to read data from multiple input channels and send the data to a single output channel