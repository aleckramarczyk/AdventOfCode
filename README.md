## Running Locally

***

```shell
AOC_SESSION_COOKIE=cookie go run cmd/main.go 
```
This will default to running both parts of todays solution.

You can change this behavior with these arguments:
- `-y 2023`
  - specifies the year of the solution you want to run
- `-d 1`
  - specifies the day of solution 
- `-p 1`
  - specifies which part of the solution
  - use 0 or leave unspecified to run both parts
- `-s cookie`
  - pass session cookie as an argument
  - will override environment variable