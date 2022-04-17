#### package dieselfluid/math

**package math**

>Mathematical tools for working with 32-bit floating point numbers for rendering
>purposes typically. If you require scientific numerical computing consider using
>the golang standard package gonum

>All package numerical operations are assumed to immutable with comments displaying
>a @mutate or @nomutate clause indicating whether the calling object state is changed.

@See package compute for numerical intergation

### Usage
* Vectors should be declared explicilty with their array structure lenghts identified
it is considered an error to declar a new vector with `var x = Vec{}` this places
  a 0 length vector array *

> Create new vector array objects
```
var x = Vec{1, 2, 3}
var y = Vec{1, 1, 1}
var eq = float32(6.0)
  if Dot(x, y) == eq && Dot(x, y) == eq {
  } else {
    t.Errorf("Vector failed %f", x[0])
  }
}
```
