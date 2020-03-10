# macro ![wip](https://img.shields.io/badge/-work%20in%20progress-red) ![draft](https://img.shields.io/badge/-draft-red)
A *go generate* macro processor, using AST processing and go itself - inspired by rust macros.
Technically, a temporary go module is created for the actual go module which is solely
created based on your #[...] macro directives. It is invoked with the according AST
input and output context parameters.

## example

Use `go get github.com/ee4g/macro` in your project and create a file like *gen.go*:
```go
package main

import "github.com/ee4g/macro"

//go:generate go run gen.go

func main(){
    macro.MustApply()
}
```

Usage is as follows
```go
package main

import "fmt"

// Import macros should be the first macros in a file and provide the scope
// for all following macros. *import* is a reserved keyword. You need
// to also declare a
//  #[require] github.com/mycompany/mypackage 1.2.3
// within your go.mod file.
//
// #[import] github.com/mycompany/mypackage // just a normal go import


// A macro is a bunch of method calls within a comment group.
// They are executed with the current context which is the 
// directly following AST type, ignoring any comment).
// So, the following macros are applied to Entity.
//
// #[mypackage.Stringer()]
// #[mypackage.SQL("my_entity_table")

// Entity has a normal comment which is not polluted from the 
// macro declaration above. Macros are implementation specific and therefore
// unimportant to document here.
type Entity struct{
    ID string
    Name string
}

/*
 This is also a multiline macro, applied to the method hello below.
    #[
        mypackage.Get("/hello/{name}/say")
        mypackage.OpenAPI()
        // a macro comment
        mypackage.SecurityRole(&mypackage.Role{"admin}") // another comment
    ]
*/

// hello is a Rest endpoint.
func hello(name string)string{
    return "hello "+name
}
```

To define a custom macro, simply declare it like this:
```go
package mypackage

// Macros is a hardcoded name and must always be present
// and is created once per macro.
type Macros struct{
    Method tbd.Method // Method is injected automatically
    Type tbd.Type // Type is injected automatically
    AST tbd.Type // AST is injected and contains the annotated tree
    Out tbd.Out // Out is injected and is used to generate code
}

// the Macro signature must be reflected here
func (m *Macros) OpenAPI(){
    // whatever
}

func (m *Macros) Get(route string){
    // whatever
}

func (m *Macros) SecurityRole(r *Role){
    // whatever
}

// you can define custom types
type Role struct{
    Name string
}
```