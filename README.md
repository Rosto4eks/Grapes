# grapes
grapes is a minimalistic router written in golang

provides the ability to handle parameterized requests, send files, serve static requests, and more

- [Installation](#installation)
- [Example](#example)
- [Tree structure](#tree-structure)
- [Context](#context)
- Parameters:
  - [Catch-all parameters](#catch-all-parameters)
  - [Named parameters](#named-parameters)
  - [Query parameters](#query-parameters)
- Serving files:
  - [Send file](#send-file)
  - [Serving static files](#serving-static-files)
- Forms:
  - [Parse Form values](#parse-form-values)
  - [Parse Form files](#parse-form-files)
- Additional:
  - [HttpNotFound](#http-not-found)
  
## Installation
To use grapes just create .mod file and add package
```
go get -u github.com/Rosto4eks/grapes
```
## Example
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  // create new router
  r := grapes.NewRouter()

  r.Get("/", func(ctx grapes.Context) {
    ctx.String("Hello World!")
  })

  r.Get("/Home", func(ctx grapes.Context) {
    ctx.Json(grapes.Obj{
      "title": "Home",
      "id": 1,
    })
  })

  r.Post("/Home", func(ctx grapes.Context) {
    id := ctx.Query("Id")
    ctx.Json(grapes.Obj{
      "id": id,
    })
  })

  // start listening at port 80
  r.Run(80)
}
```
```
GET  /          -> "Hello World!"
GET  /Home      -> {"id": 1,"title": "Home"}
POST /Home?id=1 -> {"id": 1}
```
## Tree structure
The router applies tree structure which allows to use parameterized requests like /*, /:id

tree consists of nodes, each node has Handlers and Childrens (next nodes)

when the request comes, router splits the url and searches into tree suitable node 
```
- / 
  ├ /*
  ├ /:id
  |    ├ /name
  |    └ /*
  └ /Home
        ├ /info
        └ /:other
```
## Context
Context is the superstructure over http.ResponseWriter and http.Request

it also provides some additional functionality for processing requests and sending responses, like sending JSON, files as response, extarcting parameters from request url
```go
from
func (w http.ResponseWriter, r *http.Request) {
  ...
}
to
func(ctx grapes.Context) {
  ...
})
```

## Catch-all parameters
symbol * means, that router will handle any request with any part instead of symbol *
```go
package main

import "github.com/Rosto4eks/grapes"

func catchAll(ctx grapes.Context) {
  ctx.Json(grapes.Obj{
    "title": "Home",
  })
}

func catchInfo(ctx grapes.Context) {
  ctx.Json(grapes.Obj{
    "title": "Info",
  })
}

func main() {
  r := grapes.NewRouter()

  r.Get("/*", catchAll)

  r.Get("/*/Info/*", catchInfo)

  r.Run(80)
}
```
```
/                       -> 404 page not found
/credits                -> {"title": "Home"}
/Home                   -> {"title": "Home"}
/Home/golang            -> 404 page not found
/Home/Info              -> 404 page not found
/Home/Info/github       -> {"title": "Info"}
/Home/Info/chess        -> {"title": "Info"}
/Home/Info/chess/queen  -> 404 page not found
```

## Named parameters
They are like catch-all params, but you can extract param value from url 

function ctx.NamedParam(param) will return it
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()

  r.Get("/:id", func (ctx grapes.Context) {
    id := ctx.NamedParam("id")
    ctx.Json(grapes.Obj{
      "id": id,
    })
  })

  r.Get("/:id/:name", func (ctx grapes.Context) {
    id := ctx.NamedParam("id")
    name := ctx.NamedParam("name")
    ctx.Json(grapes.Obj{
      "id": id,
      "name": name,
    })
  })

  r.Run(80)
}
```
```
/                       -> 404 page not found
/0x3DE3E                -> {"id": "0x3DE3E"}
/0x3DE3E/Rosto4eks      -> {"id": "0x3DE3E", "name": "Rosto4eks"}
/0x3DE3E/Rosto4eks/info -> 404 page not found
```
## Query parameters
quoery params can be obtained using the function ctx.Query(param)
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()

  r.Get("/Home", func (ctx grapes.Context) {
    id := ctx.Query("id")
    ctx.Json(grapes.Obj{
      "id": id,
    })
  })

  r.Run(80)
}
```
```
/Home?id=1 -> {"id": "1"}
```
## Send file
To send file, use function ctx.File(path)
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()

  r.Get("/", func (ctx grapes.Context) {
    ctx.File("public/home.html")
  })

  r.Run(80)
}
```

## Serving static files
To serve static files, use function router.Static(path, pattern)

pattern must start with '/' symbol

The example below shows how to use html with css (path in link to css will look like "/public/css/home.css")

Folder structure:
```
main.go
public
     ├ home.html
     └ css
         └ home.css
```
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()
  r.Static("public", "/public")
  r.Get("/", func (ctx grapes.Context) {
    ctx.File("public/home.html")
  })

  r.Run(80)
}
```

## Parse form values
html form wich will be used in the next 2 examples
```html
<form action="/" method="post" enctype="multipart/form-data">
  <input name="name" type="text">
  <input name="password" type="text">
  <input name="file" type="file">
  <input type="submit">
</form>
```
To parse values from form, use ctx.GetFormValue(key)
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()
  r.Static("public", "/public")

  r.Get("/", func (ctx grapes.Context) {
    ctx.SendFile("public/home.html")
  })
  r.Post("/", func(ctx grapes.Context) {
    name, _ := ctx.GetFormValue("name")
    password, _ := ctx.GetFormValue("password")

    ctx.SendJson(grapes.Obj{
      "name":name,
      "password":password,
      })
  })

  r.Run(80)
}
```

## Parse form files
To parse form file, use ctx.GetFormFile(key)
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()
  r.Static("public", "/public")

  r.Get("/", func (ctx grapes.Context) {
    ctx.SendFile("public/home.html")
  })
  r.Post("/", func(ctx grapes.Context) {
    file, header, _ := ctx.GetFormFile("file")
    ctx.SendJson(grapes.Obj{
      "fileName": header.Filename,
      "fileSize": header.Size,
    })
  })

  r.Run(80)
}
```

## Http-not-found
router.HttpNotFound is the build-in handler wich provides processing unhandled routes, also you can change this fuction:
```go
package main

import "github.com/Rosto4eks/grapes"

func main() {
  r := grapes.NewRouter()
  
  r.HttpNotFound = func(ctx grapes.Context) {
    ctx.SendString("Oops! Page not found.")
  }

  r.Get("/", func (ctx grapes.Context) {
    ctx.SendString("Hello")
  })

  r.Run(80)
}
```
```
/          -> "Hello"
/Home      -> "Oops! Page not found."
/Home/Info -> "Oops! Page not found."
```
