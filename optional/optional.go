package optional

import (
    "reflect"
    "fmt"
)

type Optional[T any] struct {
    hasValue bool
    value    T
}

func NewOptional[T any]() Optional[T] {
    opt := Optional[T]{}
    opt.hasValue = false;
    return opt;
}

func NewOptionalWithVal[T any](value T) Optional[T] {
    opt := Optional[T]{}
    opt.SetValue(value)
    return opt;
}

func (this *Optional[T]) String() string {
    if this.hasValue {
        return fmt.Sprint(this.value)
    } else {
        return "<"+reflect.TypeOf(this.value).String()+"?>"
    }
}

func (this *Optional[T]) HasValue() bool {
    return this.hasValue;
}

func (this *Optional[T]) GetValue() T {
    if !this.hasValue {
        panic("Called `GetValue()` on an `Optional` with no value")
    }
    return this.value;
}

func (this *Optional[T]) GetValueOr(def T) T {
    if this.hasValue {
        return this.value;
    } else {
        return def;
    }
}

func (this* Optional[T]) SetValue(val T) {
    this.value = val
    this.hasValue = true
}
