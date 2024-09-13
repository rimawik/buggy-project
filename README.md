# buggy-project
Fixed issues in a buggy project

## Issue 1:
Wrong use of WaitGroup (line 16) 
WaitGroup should be used in situations where different goroutines have to synchronize their completion for a common task. for HTTP requests, each request is independent, so we don't need to group the APIs together using a global WaitGroup. each API has to have its own WaitGroup.
Using the same Waitgroup will have a negative impact. For instance, if one request increments the WaitGroup while another request is finishing, the second request may not return until the first request completes, potentially leading to deadlocks.


## Issue 2:
Lack of error handling(line 21, and line 40)
all possible errors have to be handled correctly, so the developer knows what are the issues that are causing the program to fail.


## Issue 3:
DB connection not closed (line 21)
we need to call `defer db.Close` to make sure the connection is closed properly when the app terminates to prevent resource exhaustion on the database side.

## Issue 4:
Added db.Ping() to ensure that the connection is valid (line  26). If the connection to db couldn't be established, it is better we terminate the app.

## Issue 5:
missing status code (line 79)
Without specifying the error code, the request will return an error, but the status code will be 200, indicating success.
