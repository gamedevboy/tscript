# TinyScript

- TinyScript is a register based fast script language written by golang.
- It's VM design inspired by LUA and can only used on 64 bit machine.
- The runtime provide `FAST FIELD ACCESS` and `FUNCTION INLINE CACHE` like V8.

## GET STARTED

### Basic Runtime Enveriment

```golang

scriptAssembly := &struct {
    *assembly.Component
}{}

scriptAssembly.Component = assembly.NewScriptAssembly(scriptAssembly)
loader.LoadAssembly(scriptAssembly, "./scripts/player/player")

scriptContext := &struct {
    *context.Component
}{}
scriptContext.Component = context.NewScriptContext(scriptContext, scriptAssembly, 1024)

scriptContext.Run()

```

## FUNCTION CALL

Functions can be started by 'function', 'func' or '#', and can be used like JavaScript

Example

```
func Hello() {
    function nest() {
        println("Hello World")
    }

    nest()
}

Hello()
```

## PRIMARY TYPES

Just like JavaScript, TinyScript is a dynamic language, and it only provide very few built-in types;
`Number`, `String`, `Bool`, `Array`, `Map`, `Object`

You can define your own type by using ES6 syntax

```
class Example : Foo {
    constructor(name) {
        this.name = name
    }
}

var e = new Example('Big')
println(e.name)
```

## USING MUTI-FILES

In the script files, you can use `#import "..."` to include other files.

> ALL FILES IMPORTED SHARE THE SAME SCRIPT CONTEXT

## CONTRIBUTE GUID

Any one with any kind of contributions are welcomed, so please help to impove this project.

- Compiler
- Runtime
- Interpreter
- Documents
- Test Cases
- Examples

## TO BE CONTINUED
