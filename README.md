# go-lastfm-to-twitter

simple cli tool to publish your lastfm top to twitter

![last.fm top](http://s3-eu-central-1.amazonaws.com/irlndts.moscow/wp-content/uploads/2018/12/19171411/2018-12-19_17-13-30-e1545228962769.png)

### Installation

```
go install -u github.com/irlndts/go-lastfm-to-twitter
```

### Usage
```
go-lastfm-to-twitter list \
  --user=kleto4kin \
  --period=week \
  --publish=true \
  --lastfm-key=<KEY> \
  --twitter-consumer=<KEY> \
  --twitter-token=<KEY> \
  --twitter-secret=<KEY> \
  --twitter-token-secret=<KEY>
```
  
Find lastfm key [here](https://www.last.fm/api/account/create) and twitter keys [here](https://developer.twitter.com/en/apps)
