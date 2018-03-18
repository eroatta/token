# dispersal

Tool for splitting and expanding identifiers, supporting the following algorithms:

* **Conserv**: This is the reference algorithm. Each token is split using separation markers as underscores, numbers and regular camel case.
* **Greedy**: This algorithm is based on a greedy approach and uses several lists to find the best splits, analyzing the token looking for preffixes and suffixes.
* **Samurai**: This algorithm splits identifiers into sequences of words by mining word frequencies in source code.
* **GenTest**

## Usage

### GET /splits?identifier=$token

```
    [
        "sub"'
        "sub2",
        "sub3"
    ]
```
