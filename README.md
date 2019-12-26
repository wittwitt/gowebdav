# gdav

webdav server by golang !

base on golang.org/x/net/webdav 

# introduce

multi user

user power limit
  
  hide some file or subfolder

# hide syntax


only support path prefix

```conf
hides = [ 
	"oneproject/.git/",
	"down/",
	"cc/a.txt",
]
```
> only math ${root}/cc/a.txt
> no math ${root}/cc/a.txt/../..
> math ${root}/down/../..
> math ${root}/oneproject/.git/../..