# gdav

webdav server by golang !

base on golang.org/x/net/webdav 

# introduce

multi user

user power limit
  
  hide some file or subfolder

# hide syntax

https://golang.org/pkg/path/filepath/#Glob

```sh

'*'         matches any sequence of non-Separator characters
	'?'         matches any single non-Separator character
	'[' [ '^' ] { character-range } ']'
	            character class (must be non-empty)
	c           matches character c (c != '*', '?', '\\', '[')
	'\\' c      matches character c

character-range:
	c           matches character c (c != '\\', '-', ']')
	'\\' c      matches character c
	lo '-' hi   matches character c for lo <= c <= hi

```