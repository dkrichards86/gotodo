# gotodo
gotodo is a CLI client to manage your todos. It uses the ubiquitous [todo.txt](http://todotxt.org/)
format for specifying todo attributes, and saves them to a file on your local system.

![Tests](https://github.com/dkrichards86/gotodo/workflows/Tests/badge.svg)
![Lint](https://github.com/dkrichards86/gotodo/workflows/Lint/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dkrichards86/gotodo)](https://goreportcard.com/report/github.com/dkrichards86/gotodo)

## Usage
```
NAME:
   gotodo - CLI client to manage your todos

USAGE:
   gotodo [global options] command [command options] [arguments...]

COMMANDS:
   list, ls      Shows a lists of your todos
   add           Creates a new todo
   edit          Edits an existing todo
   pri           Updates the priority of a todo
   depri         Removes the priority from a todo
   addproject    Adds a new project to a todo
   addcontext    Adds a new context to a todo
   addattribute  Adds a new attribute to a todo
   complete, do  Marks a todo as complete
   resume        Marks a todo as incomplete
   remove, rm    Removes a todo
   projects      Shows a list of projects
   contexts      Shows a list of contexts
   attributes    Shows a list of custom attributes
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Contributing

Ff you spot bugs or have features that you'd really like to see in gotodo, please check out the 
[contributing page](.github/CONTRIBUTING.md)