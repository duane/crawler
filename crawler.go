package main

import (
  "code.google.com/p/go-html-transform/css/selector"
  "code.google.com/p/go-html-transform/h5"
  "io"
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
  resp, err := http.Get(root_url)
  if err != nil {
    log.Fatal(err.Error())
  }

  links, err := get_links(resp.Body)
  if err != nil {
    log.Fatal(err.Error())
  }
}
