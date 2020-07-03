## Install

```bash
go get -u github.com/IrvinYoung/gutil
```

## Basic usage

```go
import "github.com/IrvinYoung/gutil/aut"

// init default inscance
if err := aut.InitDefaultJwtEngine("<custom sign>", 6*time.Second, nil, nil); err != nil {
	log.Fatalln(err)
}

// create a JWT with custom data
m := map[string]interface{}{
  "id":"111",
  "name":"alice",
}
JwtStr, err := aut.NewJwt(m)
if err != nil {
	t.Fatal(err)
}
log.Println("token=",JwtStr)

// verify JWT
r, err := aut.VerifyJwt(JwtStr)
if err != nil {
	log.Println("verify failed =", err)
} else {
	log.Println("PASS: token content = %+v\n", r)
}

```

## Using custom instance 

```go
import "github.com/IrvinYoung/gutil/aut"

//define Invalidate function 
func Invalidate(fields ...interface{}) {
	id := fields[0].(string)
	expire := int64(fields[1].(float64))

	tm := expire - time.Now().Unix()
	if tm <= 0 {
		return
	}
	cli, err := initRedisCache()
	if err != nil {
		log.Println(err)
	}
	err = cli.Set(fmt.Sprintf("%s_%d", id, expire), "1", time.Duration(tm)*time.Second).Err()
	fmt.Println("invalidate result =",err)
}

//define IsScrap function
func IsScrap(m map[string]interface{}) (err error) {
	id := m["id"].(string)
	expire := int64(m["expire"].(float64))
	key := fmt.Sprintf("%s_%d", id, expire)

	cli, err := initRedisCache()
	if err != nil {
		log.Println(err)
	}
	val := cli.Exists(key).Val()
	if val >= 1 {
		err = errors.New("scrap token")
		return
	}
	return
}

//we record scrap token by redis
func initRedisCache() (cli *redis.Client, err error) {
	cli = redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	})
	pong, err := cli.Ping().Result()
	if err != nil {
		return
	}
	if pong != "PONG" {
		err = errors.New("redis don't have a right pong result =" + pong)
		return
	}
	return
}

//init custom inscance
if err := aut.InitDefaultJwtEngine("your Sign", 6*time.Second, Invalidate, IsScrap); err != nil {
		log.Fatalln(err)
	}

//create a token
m := map[string]interface{}{
  "id":"111",
  "name":"alice",
}
JwtStr, err := aut.NewJwt(m)
if err != nil {
	log.Println(err)
}
log.Println("JWT =", JwtStr)

//verify
r, err := aut.VerifyJwt(JwtStr)
if err != nil {
	log.Println("verify failed =", err)
} else {
	log.Println("PASS: token content = %+v\n", r)
}

//invalidate
aut.InvalidateJwt(m["id"], r["expire"])

//expire
r, err = aut.VerifyJwt(JwtStr)
if err != nil {
	log.Println("EXPIRE: verify failed =", err)
}

```

