module example

go 1.20

replace signature => ../signature

replace extension => ../extension

require signature v0.1.0

require extension v0.1.0

require github.com/loopholelabs/polyglot v1.1.3 // indirect
