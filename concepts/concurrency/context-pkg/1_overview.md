### Context Package 


-> In go servers, each incoming request is handled in its own goroutines

-> we need a way to propagate request-scoped data down the call graph 


Context Package servers two primary purpose:

-> provides API for cancelling branches of call-graph

-> provides a data-bag for transporting request-scoped data through call-graph


A context carries a deadline, cancelation signal and request-scoped values across API boundaries, Its methods are safe for simultaneous use by multiple goroutines

```
type Context interface {

    // done returns a ch that is closed when this context is ccanceled or times out
    Done() <-chan struct{}


    // Err indicates why this context was canceled, after the done channel is closed
    Err() error

    // Deadline returns the time when this context will be canceled, if any
    Deadline()( deadline time.Time, ok bool)

    // Value returns the val;ue associated with key or nil if none
    Value(Key interface{}) interface{}
}

```

Context package provides functions to create new context

-> context.Background()

    - background returns a empty context
    - root of any context tree
    - never canceled, has no value, has no deadline
    - typically used by main function
    - acts as a top level context for incoming request

-> context.TODO()

    - TODO() return an  empty context
    - TODO intended purpose is to server as placeholder


What is context package used for ?

context package can be used to send

    -> Request-scoped values
    -> Cancellation signals



## Context package for cancellation

-> Context is immutable

-> Context package provides function to add new behavior

-> To add cancellation behavior we have functions like

    - context.WithCancel()

        - withcancel returns a copy of parent with a new done channel
        - channel can be used to close context done channel
        - closing the done channel indicated to an operation to abandon its work and return 
        - cancelling the context return releases the resources associated with it 
        - cancel() does not wait for the work to stop
        - cancel() may be called by multiple goroutines simultaneously
        - after the first cancel , other calls will do nothing


    - context.WithTimeout()

        - WithTimeout() takes parent context and time duration as input

        -Withtimeout returns a new context that closes its done channel after the given timeout duration

        - WithTimeout is useful for setting a deadline on the requests to backend servers

        - withTimout is a wrapper over with WithDeadline()

    - context.WithDeadline()

        - withdeadline takes parent context and clock time as input

        - withdeadline returns a new context that closes its done channel when the machine clock advances past the given deadline

-> The derived context is passed to child goroutines to facilitate their cancellation



Difference is using WithTimeOut and WithDeadline 

-> WithTimeout - timer countdown begins from the moment the context is created

-> WithDeadline() -  Set explicit time when timer will expire

Summary:

How context package can be used for cancellation of an operation?

->Context pkg provides functions to derive new context values from existing ones to add cancellation behavior

-> context.WithCancel() - is used to create a cancellation context

-> cancel() is used to close the done channel

-> On receiving the close signal, goroutine is suppose to abandon it operation and return


Summary:

    - context.WithDeadline(): is used to set deadline to an operation

    - context.WithDeadline(): creates a new context, whose done channel get closed when machine clock advances past the given deadline

    - ctx.Deadline() can used to know if a deadline is associated with the context

## Context Package as Data bag

- Context package can used to transport request-scoped data down the call graph

- context.WithValue() provides a way to associate request-scoped values with a context


## Go Idioms to follow


# Incoming requests to a server should create a Context

-> create context early in processing task or request

-> create a top level context

-> http.Request value already contains a context


# Outgoing calls to servers should accpet a Context

-> Higher level calls need to tell lower level calls how long they are willing to wait


-> http.DefaultClient.DO() method to respect cancellation signal on timer expiry and return with error message


# Pass a Context to function performing I/O

-> Any function that is performing I/O should accpet a context value as its first parameter and respect any timeout or deadline configured by the caller

-> Any API that takes a Context, the idiom is to have the first parameter accept the Context value

# Any change to a Context value creates a new Context value that is then propagated forward

When a Context is canceled, all Contexts derived from it are also canceled

    - If a parent Context is cancelled, all children derived by that parent Context are cancelled as well

# Use TODO context if you are unsure about which Context to use

-> if a function is not responsible for creating top level context

-> we need a temporary top-level Context until we figured out where the actual Context will come from

# Use Context values only for request-scoped data

-> Do not use the Context  value to pass data into a function which becomes essential for its successful execution

-> a function should be able to execute its logic with an empty context Value



Summary:

    -> Incoming res to a server should create a Context

    -> Outgoing calls to servers should accept a Context

    -> Any function that is performing I/O should accept a Context value

    -> Any change to a Context value creates a new Context value that is then propagated forward

    -> if a parent Context is cancelled, all children derived from it are also cancelled

    -> use TODO context if you are unsure about which Context to use

    -> Use context values only for request-scoped data, not for passing optional params to functions