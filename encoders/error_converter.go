package encoders

import "reflect"

func ConvertErrorsToString(i *[]interface{}) {
	val := reflect.ValueOf(i).Elem()
	count := val.Len()
	for i := 0; i < count; i++ {
		elemVal := val.Index(i)
		elem := elemVal.Interface()
		switch elem.(type) {
		case error:
			elemVal.Set(reflect.ValueOf(elem.(error).Error()))
		}
	}
}
