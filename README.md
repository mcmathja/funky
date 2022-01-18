# funky

funky is a zero-dependency library providing generic operations over Go collections in order to support functional programming paradigms.

In other words, it makes Go a bit more funky.

## What does it do?

funky supports your standard Map, Filter, and Reduce, but also a bunch of other operations. Have a look at the documentation for each package to see what's available.

## Isn't Go imperative / procedural / etc.?

Go is a wonderfully concise language with a mercifully small feature set. Included among those carefully curated functionalities are first-class functions and closures. If the language designers did not intend for it to be used in a functional way, then that was a big mistake!

## I don't really understand generics...

That's okay! We did all of the hard work so that you don't have to. Well, us and the Go compiler's steadily improving type inference. The only time you'll likely see type arguments is when you're declaring an empty collection.

## Why can't I use method chaining?

Many other languages define collection operations as methods on a collection type, allowing you to combine several sequential operations into a single statement. Consider this example from Scala:

```scala
List(1, 2, 3).map(_ * 2).filter(_ % 3 != 0).reduce(_ + _) // 6
```

Go's generics implementation has a unique restriction with regard to methods: while a function can be generic with respect to as many type parameters as you choose, a method can only be generic with respect to the type parameters of the struct that method is defined on. This means that functions like Filter can be implemented as methods, since it only operates on the type of the collection, while functions like Map cannot be implemented as methods, since it requires an additional type that is not (necessarily) the type of the collection:

```go
// This is okay, because it only uses T, which is one of slice's type parameters.
func (Slice[T]) Filter([]T, func(T) bool) []T

// This is NOT okay, because it uses U, which is not one of slice's type parameters.
func (Slice[T]) Map([]T, func(T) U) []U
```

The details around why this is difficult (and perhaps impossible) for Go to accomplish are related to how generic methods satisfy interfaces - a detailed discussion of the issue can be found in the generics proposal [here](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods). But note that even if fully generic method definitions were to become available in later versions of Go, they would still need to be defined against type aliases (as alluded to in the code sample above: `Slice[T]` instead of `[]T`), which would presumably introduce a great deal of obnoxious type conversions (`Slice[T](mySlice)`, `[]T(myFunkySlice)`) for clients of a collections library.

So simply put, not all common collection operations can be implemented as methods. Rather than split the difference and implement some operations as methods and others as package level functions, funky implements them all as package level functions for consistency's sake.

