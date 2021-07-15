package commands

import "github.com/valyala/fastjson"

func fromJsonValue(value *fastjson.Value) Req {
	content := Req{}
	content.id = value.GetInt64("id")
	props := value.GetObject("properties")
	content.preferredNames = make(map[string]string)
	content.variantNames = make(map[string]string)
	content.metadata = make(map[string]string)
	props.Visit(func(key []byte, v *fastjson.Value) {

	})
	content.latitude = value.GetFloat64("geom:latitude")
	return content
}
