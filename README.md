# ltime
time.now消耗性能，所以写这个包

时间精度是20毫秒左右，对时间要求不高的，可以直接调用，原本写的是用锁的方式，现在改成了用原子锁的方式
