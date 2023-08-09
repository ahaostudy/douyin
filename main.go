package main

import "main/router"

func main() {
  r := router.InitRouter()

  if err := r.Run(); err != nil {
    panic(err)
  }
}
