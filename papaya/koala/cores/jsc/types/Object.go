package types

// property descriptor

type JSPropDesc struct {
  Value        any
  Writable     bool
  Get          func(any)
  Set          func() any
  Configurable bool
  Enumerable   bool
}

// [properties] Object

type JSObjectImpl interface {
  Length() int
  Name() string
  Prototype() any
  Assign(target any, sources ...any) any          // #Assign<>
  GetOwnPropertyDescriptor(obj any, prop any) any // JSPropDesc
  GetOwnPropertyDescriptors(obj any) any          // KMap<JSPropDesc>
  GetOwnPropertyNames(obj any) []string
  GetOwnPropertySymbols(obj any) any // Array<T>
  Is(value1 any, value2 any) bool    // Compare<T>
  PreventExtensions(obj any) any     // #Prevent<Extension>, #Prevent<Assign>
  Seal(obj any) any
  Create(proto any, propertiesObject any) any              // create new instance of `object`
  DefineProperties(obj any, props any) any                 // props like `JSPropDesc` with `KMap`
  DefineProperty(obj any, prop string, descriptor any) any // descriptor like `JSPropDesc` with name as `string`
  Freeze(obj any) any                                      // make it prevent extensible, non writable, configurable
  GetPrototypeOf(obj any) any                              // like `JSObjectPrototype`
  SetPrototypeOf(obj any, prototype any) any               // set prototype as `object`
  IsExtensible(obj any) bool
  IsFrozen(obj any) bool
  IsSealed(obj any) any
  Keys(obj any) []string
  Entries(obj any) []any          // like `[]KMap`
  FromEntries(iterable []any) any // like `[]KMap` into `JSObject`
  Values(obj any) []any
  HasOwn(obj any, prop string) bool // has own `properties`
}

// [properties] Object.prototype

type JSObjectPrototypeImpl interface {
  Constructor(args ...any) any // instance of `this`
  HasOwnProperty(prop string) bool
  IsPrototypeOf(obj any) any
  PropertyIsEnumerable(prop string) bool
  ToString() string
  ValueOf() any                                         // set `ValueOf` can make it conversion into `any`
  ToLocaleString(locales string, options string) string // set language to make it conversion
}
