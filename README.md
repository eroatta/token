# token-splitex [WIP]

Collection of token splitting and expanding algorithms.
The following lists show the supported algorithms.

## Splitting algorithms

* **Conserv**: This is the reference algorithm. Each token is split using separation markers as underscores, numbers and regular camel case.
* **Greedy**: This algorithm is based on a greedy approach and uses several lists to find the best split, analyzing the token looking for preffixes and suffixes. *Related paper:* [Identifier splitting: A study of two techniques (Feild, Binkley and Lawrie)]([https://link](https://www.academia.edu/2852176/Identifier_splitting_A_study_of_two_techniques))
* **Samurai**: This algorithm splits identifiers into sequences of words by mining word frequencies in source code.
* **GenTest**: **[WIP]**

## Expansion algorithms

* **Basic**: *TODO*
* **AMAP**: *TODO*
* **Normalize**: *TODO*

## Usage

### Conserv

A token can be splitted just by calling the splitting function: `conserv.Split(token)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token-splitex/conserv"
)

func main() {
    splitted := conserv.Split("httpResponse")

    fmt.Println(splitted) // [http response]
}
```

### Greedy

Greedy looks for the longest prefix and the longest suffix that are "on a list" (i.e. in the dicctionary, on the list of abbreviations, or on the stop list), so it requires the list to be passed as a parameter. We can build the list using the provider `ListBuilder`.
Once we have our list, we can call the splitting function on Greedy, providing the token and the list of words: `greedy.Split(token, list)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token-splitex/greedy"
)

func main() {
    listBuilder := greedy.NewListBuilder()
    list := listBuilder.Dicctionary(lists.Dicctionary).
        KnownAbbreviations(lists.KnownAbbreviations).
        StopList(lists.Stop).build()

    splitted := greedy.Split("httpResponse", list)

    fmt.Println(splitted) // [http response]
}
```
