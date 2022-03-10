# Channel

## Inherently interesting

- goroutine-safe
  - hchan **mutex**
- store and pass values between goroutine
  - copying into and out of hchan **buffer**
- provide FIFO semantics
  - hchan **sudog queues**
  - calls into the runtime scheduler(gopark, goready)
- can cause goroutines to block and unblock
