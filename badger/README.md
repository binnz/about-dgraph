## badger特性
separates keys from values to minimize I/O amplification
slower for range key-value iteration Now
The keys are stored in LSM tree, while the values are stored in a write-ahead log called the value log
When required, the values are directly read from the log stored on SSD, utilizing its vastly superior random read performance.（SSD-centric design.）
compaction process: multiple files are read into memory, sorted, and written back. Sorting is essential for efficient retrieval, for both key lookups and range iterations



CAP:
difficult to come up with a solution that could be both “Consistent and Available”



Mutex
互斥锁两种模式
正常模式：排队最前面的goroutine和新来的一起竞争。新来的已经在CPU当中，所以更可能抢到锁。
饥饿模式：1ms后唤醒等待的goroutine没有抢到，锁的所有权会直接从释放锁(unlock)的goroutine转交给队列头的goroutine，其他新来的不再自选抢锁，自觉去排队。
什么时候再切回正常模式：If a waiter receives ownership of the mutex and sees that either (1) it is the last waiter in the queue, or (2) it waited for less than 1 ms（拿到锁所花的时间小于1ms）, it switches mutex back to normal operation mode.

type Mutex struct {
	state int32
	sema	uint32
}

state 
最低位1表示锁是否加锁
第2位表示锁是否被唤醒
第3位标记这把锁是否为饥饿状态

sema
sema is used to provide the function of sleeping and waking goroutine, which is equivalent to a waiting queue.


## DB

type lockedKeys struct {
	sync.RWMutex
	keys map[uint64]struct{}
}