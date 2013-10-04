package main

import (
  "code.google.com/p/go-html-transform/css/selector"
  "code.google.com/p/go-html-transform/h5"
  "github.com/jmhodges/levigo"
  "io"
  "io/ioutil"
  "log"
  "net/http"
)

const root_url string = "http://www.reddit.com"

func get_links(reader io.Reader) (links []string, err error) {
  t, err := h5.New(reader)
  if err != nil {
    return
  }

  top := t.Top()
  chain, err := selector.Selector("a")
  if err != nil {
    return
  }

  if chain == nil {
    panic("chain is nil")
  }

  link_nodes := chain.Find(top)
  for _, link_node := range link_nodes {
    for _, attr := range link_node.Attr {
      if attr.Key == "href" {
        links = append(links, attr.Val)
      }
    }
  }
  return
}

func main() {
  opts := levigo.NewOptions()
  opts.SetCache(levigo.NewLRUCache(3 << 30))
  opts.SetCreateIfMissing(true)
  db, err := levigo.Open("./level.db", opts)
  if err != nil {
    log.Fatal(err.Error())
  }

  resp, err := http.Get(root_url)
  if err != nil {
    log.Fatal(err.Error())
  }

  body_str, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  wo := levigo.NewWriteOptions()

  err = db.Put(wo, []byte(root_url), body_str)
  if err != nil {
    log.Fatal(err.Error())
  }

  _, err = get_links(resp.Body)
  if err != nil {
    log.Fatal(err.Error())
  }
}
