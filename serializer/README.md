About serialization groups.

They allows DTO on demand

right now we use the same tag serializationGroup for input and output,
which is consufing.

create and update routes are using read's config as an output
which is okay in theory but in practice sometimes we do weird thing
so we need to be able to give those a proper output serializationGoup for the sake of it

OR ...
I would argue post/put output should be get compliant
So ...
allows POST proper output and use get as backup


___________
Now, It is what's already happenning ... almost
it is not for pointer
and validation are a bit messy right now

goos: darwin
goarch: amd64
pkg: github.com/philiphil/restman/test/serializer
cpu: VirtualApple @ 2.50GHz
BenchmarkSerializer_Deserialize-14    	   21493	     55857 ns/op	    5167 B/op	       8 allocs/op
<.>
goos: darwin
goarch: amd64
pkg: github.com/philiphil/restman/test/serializer
cpu: VirtualApple @ 2.50GHz
BenchmarkSerializer_Serialize-14    	    2247	    489309 ns/op	  348702 B/op	    6268 allocs/op
PASS
ok  	github.com/phil