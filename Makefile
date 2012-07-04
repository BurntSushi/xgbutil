all: callback.go types_auto.go gofmt

install:
	go install -p 6 . ./ewmh ./gopher ./icccm ./keybind ./motif ./mousebind \
		./xcursor ./xevent ./xgraphics ./xinerama ./xprop ./xrect ./xwindow

gofmt:
	gofmt -w *.go */*.go examples/*/*.go
	colcheck *.go */*.go examples/*/*.go

callback.go:
	scripts/write-events callbacks > xevent/callback.go

types_auto.go:
	scripts/write-events evtypes > xevent/types_auto.go

tags:
	find ./ \( -name '*.go' -and -not -wholename './tests/*' -and -not -wholename './examples/*' \) -print0 | xargs -0 gotags > TAGS

loc:
	find ./ -name '*.go' -and -not -wholename './tests*' -and -not -name '*keysymdef.go' -and -not -name '*gopher.go' -print | sort | xargs wc -l

ex-%:
	go run examples/$*/main.go

gopherimg:
	go-bindata -f GopherPng -p gopher -i gopher/gophercolor-small.png -o gopher/gopher.go

