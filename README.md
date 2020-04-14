# token

[![GoDoc](https://godoc.org/github.com/eroatta/token?status.svg)](https://godoc.org/github.com/eroatta/token)
[![Go Report Card](https://goreportcard.com/badge/github.com/eroatta/token)](https://goreportcard.com/report/github.com/eroatta/token)

Collection of token splitting and expanding algorithms.
The following lists show the supported algorithms.

## Splitting algorithms

* **Conserv**: This is the reference algorithm.
Each token is split using separation markers as underscores, numbers and regular camel case.
* **Greedy**: This algorithm is based on a greedy approach and uses several lists to find the best split, analyzing the token looking for preffixes and suffixes.
*Related paper:* [Identifier splitting: A study of two techniques (Feild, Binkley and Lawrie)]([https://link](https://www.academia.edu/2852176/Identifier_splitting_A_study_of_two_techniques))
* **Samurai**: This algorithm splits identifiers into sequences of words by mining word frequencies in source code.
This is a technique to split identifiers into their component terms by mining frequencies in large source code bases, and relies on two assumptions:
  1. A substring composed an identifier is also likely to be used in other parts of the program or in other programs alone or as part of other identifiers.
  2. Given two possible splits of a given identifier, the split that most likely represents the developer's intent partitions the identifier into terms occurring more often in the program.
* **GenTest**: This is a splitting algorithm that consists of two parts: generation and test. The generation part of GenTest generates all possible splittings; the test part, however, evaluates a scoring function against each proposed splitting.
GenTest uses a set of metrics to characterize the quality of the split.

## Expansion algorithms

* **Basic**: It's the basic abbreviation and acronym expansion algorithm, which was proposed by Lawrie. This algorithm uses lists of words from the source code, a dictionary and also a phrase list to match a token to the possible expansions.
* **AMAP**: This algorithm applies an automated approach to mining abbreviation expansions from source code.
It's based on a scoped approach which uses contextual information at the method, program, and general software level to automatically select the most appropriate expansion for a given abbreviation.
* **Normalize**: This algorithm is based on GenTest, and selects the best expansion for a given token on a context using the one that produces the highest score.

## Usage

### Conserv

A token can be splitted just by calling the splitting function: `conserv.Split(token)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/conserv"
)

func main() {
    splitted := conserv.Split("httpResponse")

    fmt.Println(splitted) // "http response"
}
```

### Greedy

Greedy looks for the longest prefix and the longest suffix that are "on a list" (i.e. in the dictionary, on the list of abbreviations, or on the stop list), so it requires the list to be passed as a parameter.
We can build the list using the provider `ListBuilder`.

Once we have our list, we can call the splitting function on Greedy, providing the token and the list of words: `greedy.Split(token, list)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/greedy"
    "github.com/eroatta/token/lists"
)

func main() {
    listBuilder := greedy.NewListBuilder()
    list := listBuilder.Dictionary(lists.Dictionary).
        KnownAbbreviations(lists.KnownAbbreviations).
        StopList(lists.Stop).Build()

    splitted := greedy.Split("httpResponse", list)

    fmt.Println(splitted) // "http response"
}
```

### Samurai

Samurai algoritm, proposed by Hill et all, receives a token and splits it based on frequency information (local and global) and two lists of common prefixes and suffixes.
For each token analysed Samurai starts by executing a _mixedCaseSplit_ algorithm, which outputs a delimited token and then applies a _sameCaseSplit_ algorithm to each part of the newly delimited token.
The source code must be mined to extract and create two string frequency tables, which are passed to Samurai as `TokenContext`.

Once we have our frequency tables and the lists of common prefixes and suffixes, we can call the splitting function on Samurai, providing the token, the context and the lists of words: `samurai.Split(token, context, prefixes, suffixes)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/samurai"
    "github.com/eroatta/token/lists"
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

    fmt.Println(splitted) // "http response"
}
```

### GenTest

GenTest requires a similarity calculator, because it relies on the fact that words (expanded words) should be found co-located in the documentation or in general text.
Also, it requires a list of words that form the context where the token is located, and an expansion set.

Once we have our similarity calculator, context words and list of possible expansions, we can call the splitting function on GenTest, providing each required parameter: `gentest.Split(token, similarityCalculator, context, expansions)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/expansion"
    "github.com/eroatta/token/gentest"
    "github.com/eroatta/token/lists"
)

func main() {
    var simCalculator gentest.SimilarityCalculator
    context := lists.NewBuilder().Add("http").Add("response").Build()
    possibleExpansions := expansion.NewSetBuilder().AddList(lists.Dictionary).Build()

    splitted := gentest.Split("httpResponse", simCalculator, context, possibleExpansions)

    fmt.Println(splitted) // [http response]
}
```

### Basic

The Basic expansion algorithm works independently on soft words in the context of the source code for a particular function.
It uses four lists of potential expansions: a list of natural-language words extracted from the code, a list of phrases extracted from the code, a list of programming language specific words referred to as a _stoplist_, and finally a natural-language dictionary.
On our implementation, the stoplist and the dictionary are merged.
A word from one of these lists matches an abbreviation if it begins with the same letter and the individual letters of the abbreviation occur, in order, in the word.

Once we have our sets of possible expansions, we can call the expansion function on Basic, providing each required parameter: `basic.Expand(token, srcWords, phraseList, regularExpansions)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/basic"
    "github.com/eroatta/token/expansion"
)

func main() {
    srcWords := expansion.NewSetBuilder().AddStrings("connection", "client").Build()
    phraseList := map[string]string{
        "json": "java-script-object-notation",
    }

    expanded := basic.Expand("json", srcWords, phraseList, basic.DefaultExpansions)

    fmt.Println(expanded) // ["java script object notation"]
}
```

### AMAP

Automatically expanding abbreviations requires the following steps: (1) identifying whether a token is a non-dictionary word, and therefore a short form candidate; (2) searching for potential long forms for the given short form; and (3) selecting the most appropriate long form from among the set of mined potential long form candidates.
The automatic long form mining technique is inspired by the static scoping of variables represented by a compiler’s symbol table.
When looking for potential long forms, the algorithm starts at the closest scope to the short form, such as type names and statements, and gradually broadens the scope to include the method, its omments, and the package comments.
If the technique is still unsuccessful in finding a long form, it attempts to find the most likely long form found within the program and Go.

Once the token scope and the reference text are defined, we can call the expansion function on AMAP, providing each required parameter: `amap.Expand(token, scope, referece)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/amap"
)

func main() {
    // func Marshal(v interface{}) ([]byte, error)
    variableDeclarations := []string{"v interface"}
    methodName := "marshal"
    methodBodyText := ""
    methodComments := []string{
        "marshal returns the java script object notation encoding of v",
    }
    packageComments := []string{
        "package json implements encoding and decoding of java script object notation as defined in rfc",
        "the mapping between java script object notation and go",
    }

    scope := amap.NewTokenScope(variableDeclarations, methodName, methodBodyText,
        methodComments, packageComments)

    reference := []string{}

    expanded := amap.Expand("json", scope, reference)

    fmt.Println(expanded) // [java script object notation]
}

```

### Normalize

Normalize is based on GenTest, and requires a similarity calculator, because it relies on the fact that words (expanded words) should be found co-located in the documentation or in general text.
Also, it requires a list of words that form the context where the token is located, and an expansion set.

Once we have our similarity calculator, context words and list of possible expansions, we can call the expanding function on GenTest, providing each required parameter: `gentest.Expand(token, similarityCalculator, context, expansions)`.

```go
package main

import (
    "fmt"

    "github.com/eroatta/token/expansion"
    "github.com/eroatta/token/gentest"
    "github.com/eroatta/token/lists"
)

func main() {
    simCalculator := similarityCalculatorMock{
        "http-response": 1.0,
    }

    context := lists.NewBuilder().Add("http").Add("response").Build()
    possibleExpansions := expansion.NewSetBuilder().AddStrings("response").Build()

    expanded := gentest.Expand("httpresp", simCalculator, context, possibleExpansions)

    fmt.Println(expanded) // [http response]
}

type similarityCalculatorMock map[string]float64

func (s similarityCalculatorMock) Similarity(word string, another string) float64 {
    var key string
    if word < another {
        key = word + "-" + another
    } else {
        key = another + "-" + word
    }

    return s[key]
}
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
