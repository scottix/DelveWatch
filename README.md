# DelveWatch
Allows you to watch files and restart the delve debugger on change. Used for development purposes only. Do not run in production.

# Install
```go get github.com/scottix/DelveWatch```

# Usage
```
Usage of ./DelveWatch:
  -api string
    	Specify the delve api version --api-version=<input> (default "2")
  -args string
    	Additional args to program
  -delve string
    	Specify the delve command to run
  -listen string
    	Specify the delve listen variable --listen=<input> (default ":2345")
  -trace
    	Trace debug
  -verbose
    	Verbose output when running
```
