# DWIM

## About

DWIM is a suite of tools, currently targeted at developers, to facilitate code generation from a semantic model.

Essentially, this is a low-code framework for defining data structures, API's, and interfaces from a declarative model, acting as the source of truth. This alleviates several pain points in the development cycle of services and data structures which evolve over time...

 1. Backwards-compatible changes
 2. Schema consistency between API's and data structures

## Code Generation

General args for generating code from input schemas

```
$ go run cmd/dwim/main.go <input-dir> <output-dir>
```

`<input-dir>` is the input directory where schemas to be translated into code exist

`<output-dir>` is where the generated code will live


### Examples

#### Generate GoLang structs from the built-in DWIM schemas for ERC-721 tokens, and Ethereum Core concepts
```
$ go run cmd/dwim/main.go schemas/dwim gen/schemas/DWIM

Parsing RDF schema for file(erc_721.ttl)
Parsing RDF schema for file(eth_core.ttl)
```

When done, the files in `gen/schemas/dwim` will be updated (if needed), to reflect the current state of the input path provided `schemas/dwim`

