# Implementing locking with Redlock / Redsync
Redlock is a locking algorithm for implementing distributed locks with Redis.  
The Golang implementation of Redlock resides at:
```html
https://github.com/go-redsync/redsync
```
in a so called Redsync library.

## Prerequisites
Excerpt from Redis website:
```yaml
First, it’s important to understand that Redlock is designed to be operated over a minimum of 3 machines with independent Redis instances. This avoids any single-point of failure in your locking mechanism (which would be a deadlock on all resources!). The other point to understand is that, while the clocks do not need to be 100% synchronized, the clocks do need to function in the same way – e.g. time moves at precisely the same pace – 1 second on machine A is the same as 1 second on machine B.
```
### Infrastructure
Excerpts from Redis website:
```yaml
If Redis is configured, as by default, to fsync on disk every second, it is possible that after a restart our key is missing. In theory, if we want to guarantee the lock safety in the face of any kind of instance restart, we need to enable fsync=always in the persistence settings. This will affect performance due to the additional sync overhead.
```
```yaml
Using delayed restarts it is basically possible to achieve safety even without any kind of Redis persistence available, however note that this may translate into an availability penalty. For example if a majority of instances crash, the system will become globally unavailable for TTL (here globally means that no resource at all will be lockable during this time).
```
## Lock Manager - Embedded
## Lock Manager - Standalone
### Topology
Several processing instances and one lock manager.
### Flow - High Level
    1. Processing instance I[n] wishes to lock the processing of payload / event E[m]. This event resides in a key value cache. In order to retrieve the payload, key KeyE[m] should be used for retrieval from the cache.
    2. To understand if processing can be started, I[n] asks the lock manager (LockM) if a lock on KeyE[m] could be obtained. Failure to get a lock on KeyE[m] would mean that another instance is processing the E[m] payload.
### Flow - Low Level
    1. Create lock manager. Lock manager would receive requests from many processing instances.  
    Requests for which the lock process was not yet initiated would be placed in memory cache with a key as arrival timestamp.  
    Once the lock process is initiated the request is deleted from memory cache.
    2. Several requests are received for the same KeyE[m], only the first one - I[n] would be used to initiate the locking process.  
    On successful lock of the first request, the remaining requests would not be honored, but placed on hold.  
    3. If an unlock is received from I[n], the remaining requests would be deleted as this unlock would represent a successful process of payload.
    4. If no unlock request has been received, after TTL the second arrived request is used to initiate the locking process. The process repeats with next request and up to sucessful processing or no more requests.
### Assumptions
- TTL value for request is global and mandatory.
- Requests are ordered as they are received.
- The lock manager can also work with one persistence / Redis instance. Of course, several instances would provide better failure prevention.