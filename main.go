package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// GetDirs get all directories under the `dir`
func GetDirs(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	paths := []string{dir}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() && strings.HasPrefix(name, ".") == false {
			paths = append(paths, GetDirs(filepath.Join(dir, name))...)
			continue
		}
	}

	return paths
}

func onChange(op fsnotify.Op, dir string, name string) {
	log.Println(op, dir, name)
}

func main() {
	// 終了検知用のチャネルを作成。SIGINTを受け取る
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)

	// Working Directoryを取得
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start watch ...", wd)
	log.Println("[Ctrl+c]: Finish watch")

	// fsnotifyのウォッチャーを作成
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 監視ルーチン
	go func() {
		for {
			select {
			case evt := <-watcher.Events:
				// evt.Nameに絶対パスが入ってるので、Working Directoryからの相対パスを取得
				dir, err := filepath.Rel(wd, filepath.Dir(evt.Name))
				if err != nil {
					log.Println(err)
					continue
				}
				// 絶対パスからファイル名（ディレクトリ名）を抽出
				name := filepath.Base(evt.Name)

				if evt.Op == fsnotify.Remove || evt.Op == fsnotify.Rename {
					onChange(evt.Op, dir, name)
					continue
				}

				// .で始まるディレクトリは無視する。.gitとか.vscodeとか
				info, err := os.Stat(evt.Name)
				if err != nil {
					log.Println(err)
					continue
				}
				if info.IsDir() && strings.HasPrefix(name, ".") {
					continue
				}

				// onCnage2か所に分かれちゃったの気持ち悪いけど、
				// 削除した後とかってもうFileInfo取得できないからしょうがないじゃん？
				onChange(evt.Op, dir, name)

				// 新しくディレクトリが追加されたら、監視対象にする！
				if info.IsDir() && evt.Op == fsnotify.Create {
					watcher.Add(evt.Name)
				}

			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	// Working Directory配下のすべてのディレクトリを取得
	// fsnotifyは子フォルダまで見てくれないので、全部監視対象にする
	targets := GetDirs(wd)
	for _, dir := range targets {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// SIGINTが来るまで待機
	<-done
	log.Println("Bye!")
}
