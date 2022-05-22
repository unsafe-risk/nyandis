package main

import (
	"fmt"

	"github.com/unsafe-risk/nyandis/transfer/message"
	"github.com/unsafe-risk/nyandis/transfer/optional"
	"google.golang.org/protobuf/proto"
)

func main() {
	a := message.Object{
		Type:       message.TYPE_BOOLEAN,
		FloatValue: optional.Float(255.),
	}
	out, err := proto.Marshal(&a)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
