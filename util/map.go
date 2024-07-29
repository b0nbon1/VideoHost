package util

func MapValues[T any](values []T, mapFunc func(T) interface{}) []interface{}{
  var newValues []interface{}
  for  _, v := range values {
    newValue := mapFunc(v)
    newValues = append(newValues, newValue)
  }
  return newValues
}
