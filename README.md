# concurrent-quiz

I'd been learning concurrency and developed a simple quiz game with GO channels.

```
$ go build && ./concurrent-quiz --file="problems.csv" --time=2
```

```

Usage of ./concurrent-quiz:
  --file string
        You can specify a path to CSV with questions so that the program can load new questions. (default "./problems.csv")
  --time int
        You can specify a time for answering the questions. (default 3)
```