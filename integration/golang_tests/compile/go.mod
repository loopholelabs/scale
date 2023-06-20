module compile

go 1.20

replace signature v0.1.0 => ../signature

replace example v0.1.0 => ../function

require signature v0.1.0

require example v0.1.0

require github.com/loopholelabs/polyglot v1.1.1 // indirect
