Goでカレントフォルダのファイルのイベントを拾うだけのコード

## How to Build

```
go get -u golang.org/x/sys/...
go get -u github.com/fsnotify/fsnotify
go build
```

## How to Play

```
go install
```

お好きなフォルダで

```
$ watch-sample
```

## 所感

`fsnotify`がディレクトリツリー全体まで見てくれないので、自分で取ってきて追加してく必要がある。
もちろん監視してるディレクトリで新しく子ディレクトリが作られたときはもちろん監視に追加してやらんといかん。
これが少しだけ面倒だけど、インターフェイスはすごく簡素で使いやすかった。
