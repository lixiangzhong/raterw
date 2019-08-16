# raterw
基于golang.org/x/time/rate 实现的限速io.Writer,io.Reader


# Writer
```go
func main() {
	b := []byte("qweasdzxc123")
	buf := bytes.NewBuffer(nil)
	ratew := raterw.NewRateWriter(buf, 4, 4) //4byte/s 的速度写入
	n, err := ratew.Write(b)
	log.Println(n, err)
}
```

# Reader

```go
	b := []byte("qweasdzxc123")
	r := bytes.NewReader(b)
	reader := raterw.NewRateReader(r, 2, 2) //2byte/s 的速度读取
	var p = make([]byte, 9) 
	n, err := reader.Read(p)
	log.Println(string(p), n, err)
```

# AddRate
```go
import ( 
		"github.com/lixiangzhong/raterw"
		"golang.org/x/time/rate"
	)
```
```go
	l := rate.NewLimiter(10, 20)
	rr := raterw.AddRateWriter(w, l)
	rw := raterw.AddRateReader(r, l)
```