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
    "github.com/eroatta/token-splitex/lists"
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

### Samurai

Samurai algoritm, proposed by Hill et all, receives a token and splits it based on frequency information (local and global) and two lists of common prefixes and suffixes. For each token analysed Samurai starts by executing a _mixedCaseSplit_ algorithm, which outputs a delimited token and then applies a _sameCaseSplit_ algorithm to each part of the newly delimited token.
The source code must be mined to extract and create two string frequency tables, which are passed to Samurai as `TokenContext`.
Once we have our frequency tables and the lists of common prefixes and suffixes, we can call the splitting function on Samurai, providing the token, the context and the lists of words: `samurai.Split(token, context, prefixes, suffixes)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token-splitex/samurai"
    "github.com/eroatta/token-splitex/lists"
)

func main() {
    localFreqTable := samurai.NewFrequencyTable()
    localFreqTable.SetOccurrences("http", 100)
    localFreqTable.SetOccurrences("response", 100)

    globalFreqTable := samurai.NewFrequencyTable()
    globalFreqTable.SetOccurrences("http", 120)
    globalFreqTable.SetOccurrences("response", 120)

    tokenContext := samurai.NewTokenContext(localFreqTable, globalFreqTable)

    splitted := samurai.Split("httpresponse", tokenContext, lists.Prefixes, lists.Suffixes)

    fmt.Println(splitted) // [http response]
}
```
