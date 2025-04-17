It would be better to do something like this


Serialize(any serialized, any medium)

1    What's serialized ?
        A pointer
            -> unpoint it and go to 1
        A slice
            -> itt trought it and go to 1
        A struct
            -> itt trough it and go to 1 IF field is not excluded
        A map
            -> itt trough it and go to 1
        A primitive
            convert (2)

I didnt did at first cause the only target was struct other cases came as edge case
now everything's there and implemeted and working
At first nothing was mastered, but now it is and now is time to create a better worflow

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