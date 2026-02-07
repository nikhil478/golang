## Deep Dive Into channels

``` ch := make(chan int, 3) ```

``` 

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32                              
	timer    *timer // timer feeding this chan. 
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters
	bubble   *synctestBubble

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}

type waitq struct {
	first *sudog
	last  *sudog
}

```

-> any goroutine to do any operation must first aquire mutex lock on the channel

-> buf is circular ring buffer which actually stores data and is for only buffered channel
-> dataqsize is size of the buffer
-> qcount is total data elements in the queue
-> sendx and recvx indicates current index of buffer where you can send data and recv data
-> recvq and sendq are waiting queue which are used to store blocked go routines the goroutines which who are blocked while send data and recv data through the channel
-> timer & synctestBubble is something new i can see

-> waitq is the linked list of goroutines, the elements in the linked list is represented by the sudog struct

```

// sudog (pseudo-g) represents a g in a wait list, such as for sending/receiving
// on a channel.
//
// sudog is necessary because the g ↔ synchronization object relation
// is many-to-many. A g can be on many wait lists, so there may be
// many sudogs for one g; and many gs may be waiting on the same
// synchronization object, so there may be many sudogs for one object.
//
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.

type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog

	elem maybeTraceablePtr // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	// waiters is a count of semaRoot waiting list other than head of list,
	// clamped to a uint16 to fit in unused space.
	// Only meaningful at the head of the list.
	// (If we wanted to be overly clever, we could store a high 16 bits
	// in the second entry in the list.)
	waiters uint16

	parent   *sudog             // semaRoot binary tree
	waitlink *sudog             // g.waiting list or semaRoot
	waittail *sudog             // semaRoot
	c        maybeTraceableChan // channel
}

```

In the sudog struct, we have the field g, which is reference to the goroutine, and elem field is pointer to the memory which contains the value to be sent, or to which the received value will be written


when we create a channel with built in function make, 
	-> hchan struct is allocated in the heap 
	-> returns a reference to the allocated memory
	-> since ch is pointer, it can be sent between the functions which can perform, send or receive operation on the channel