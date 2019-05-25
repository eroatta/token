# token-splitex [WIP]

Collection of token splitting and expanding algorithms.
The following lists show the supported algorithms.

## Splitting algorithms

* **Conserv**: This is the reference algorithm. Each token is split using separation markers as underscores, numbers and regular camel case.
* **Greedy**: This algorithm is based on a greedy approach and uses several lists to find the best split, analyzing the token looking for preffixes and suffixes. **[WIP]**
* **Samurai**: This algorithm splits identifiers into sequences of words by mining word frequencies in source code.
* **GenTest**: **[WIP]**

## Expansion algorithms

* **Basic**: *TODO*
* **AMAP**: *TODO*
* **Normalize**: *TODO*

## Usage

Every splitting algorithm implements a `Splitter` interface, which defines the required contract for an algorithm.

### Conserv

A new Conserv splitter can be created using: NewConserv().

A token can be splitted calling the splitting function: conserv.Split(token)
