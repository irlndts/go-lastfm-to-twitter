# go-lastfm
simple API wrapper for last.fm

```go
client, _ := lastfm.New("<API_KEY>")
top, := client.User.TopArtists("<USERNAME>", lastfm.PeriodWeek, 1, 0)
```
