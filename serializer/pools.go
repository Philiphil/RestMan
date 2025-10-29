package serializer

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"sync"
)

var jsonEncoderPool = sync.Pool{
	New: func() interface{} {
		return json.NewEncoder(&bytes.Buffer{})
	},
}

var xmlEncoderPool = sync.Pool{
	New: func() interface{} {
		return xml.NewEncoder(&bytes.Buffer{})
	},
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
