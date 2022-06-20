package unpacker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

// UnpackFunc must return a pointer to a struct and the resource id
// this needs to be moved to shared as well
type UnpackFunc func(s *schema.ResourceData) (interface{}, string, error)

// Universal TODO universalUnpack - implement me
func Universal(payload reflect.Type, s *schema.ResourceData) (interface{}, string, error) {
	//d := &util.ResourceData{ResourceData: s}
	//var t = reflect.TypeOf(payload)
	//var v = reflect.ValueOf(payload)
	//if t.Kind() == reflect.Ptr {
	//	t = t.Elem()
	//	v = v.Elem()
	//}
	//
	//for i := 0; i < t.NumField(); i++ {
	//	thing := v.Field(i)
	//
	//	switch thing.Kind() {
	//	case reflect.String:
	//		v.SetString(thing.String())
	//	case reflect.Int:
	//		v.SetInt(thing.Int())
	//	case reflect.Bool:
	//		v.SetBool(thing.Bool())
	//	}
	//}
	//result := KeyPairPayLoad{
	//	PairName:    d.GetString("pair_name", false),
	//	PairType:    d.GetString("pair_type", false),
	//	Alias:       d.GetString("alias", false),
	//	PrivateKey:  strings.ReplaceAll(d.GetString("private_key", false), "\t", ""),
	//	PublicKey:   strings.ReplaceAll(d.GetString("public_key", false), "\t", ""),
	//	Unavailable: d.GetBool("unavailable", false),
	//}
	//return &result, result.PairName, nil
	return nil,"", nil
}
