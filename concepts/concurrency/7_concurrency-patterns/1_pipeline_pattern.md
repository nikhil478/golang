### Pipeline

-> Process streams, or batches of data


G1 - ch1 -> G2 - ch2 -> G3 - ch3 -> G4

Stage - take data in, perform an operation on it and send the data out

## Stages

-> Seperate the concerns of each stage

-> process individual stage concurrently

-> could consume and return the same type

func square(in <-chan int) <- chan int {}

-> This enables composability of pipeline

square(sqaure(generator(2,3)))

## Image processing pipeline

Input: List of images
Output: Thumbnail images

G1 - ch1 -> G2 - ch2 -> G3 

G1: generate a list of images to process
G2: generate thumbnail images
G3: store thumbnail image to disk or transfer to storage bucket in cloud


Summary:

-> Pipelines are used to process Streams or batches of data

-> Pipelines enables us to make an efficient use of I/O and multiple cpus cores

-> pipelines is a series of stages, connected by channels