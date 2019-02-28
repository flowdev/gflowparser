# gflowparser
Flow DSL parser for the Go programming language build with flowdev/gparselib.

## Flow DSL
The flow DSL is used to show the flow of data between components. So it consists of two main objects:
- components that perform computations, I/O, etc. and
- arrows that connect components and transport the defined data into the second component.

Components can have explicit ports that they use to receive or provide data.
The default ports are `in` for input and `out` for output.

So a very simple flow looks like this:
```flowdev
in (data)-> [component1] (data)-> [component2] (data)-> out
```
and is rendered to:
![simple flow](img/simple.svg)
As you can see, the ports at the outer level (usually at the very start and end
of a flow line) have to be stated explicitly even if they are the standard
ports.

Generally new lines and comments are fine when seperating flow lines and
within parentheses (`(` and `)`) and square brackets (`[` and `]`).

### Data and data types
Multiple data for arrows are supported and can either be seperated by a comma (`,`)
to keep them on the same line or by a pipe (`|`) to have multiple lines:
```flowdev
in (data1, data2, data3)-> [component1] (data4 | data5 | data6)-> [
component2] (data1, data3 | data5, data6)-> out
```
![multiple data](img/multiData.svg)

Data types can be upper case (exported) or lower case (local). A package name
separated by dot (`.`) can be prepended to a data type.
Simple data types like `string`, `int` and `bool` don't provide much
information so instead the more descriptive name of the parameter should be
used instead. 
```flowdev
in (localDataType, ExportedLocalDataType | otherpackage.ExportedDataType)-> [
component1] (ExportedLocalDataType | descriptiveNameForString)-> out
```
![data types](img/dataTypes.svg)

### Ports
Ports have lower case names and can have an optional index (array ports).
The maximum index is fix at design time (compile time) as anything else would
be too hard to debug.
In the Go code port names are appended to the function or method name seperated
by an underscore (`_`).
```flowdev
in (d)-> myInPort [component1] out (d)-> arrayIn:1 [component2] arrayOut:1 (d)-> out
[component1] specialOut (data)-> arrayIn:2 [component2] arrayOut:2 (data)-> extraOut
```
![ports](img/ports.svg)

### Continuations

### Components

### Plugins

### Splits, merges and circles
