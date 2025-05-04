# gorm_struct_warpper
gorm 解析 struct 的 warpper

## 用法
```go
// map解析成struct
type demomap struct {
	GObj `json:"-"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (o *demomap) Scan(value interface{}) error {
	return o.GObj.Scan(value, &o)
}

func (o demomap) Value() (driver.Value, error) {
	return o.GObj.Value(o)
}

// 列表解析成struct
type demolist struct {
	GObj `json:"-"`
	list []string
}

func (o *demolist) At(index int) string {
	return o.list[index]
}

func (o *demolist) Scan(value interface{}) error {
	return o.GObj.Scan(value, &o.list)
}

func (o demolist) Value() (driver.Value, error) {
	return o.GObj.Value(o.list)
}

func (o *demolist) MarshalJSON() ([]byte, error) {
	a, e := json.Marshal(&o.list)
	return a, e
}

func (o *demolist) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &o.list)
}

```