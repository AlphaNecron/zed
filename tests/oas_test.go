package tests

import (
	"encoding/json"
	"fmt"
	"github.com/ogen-go/ogen"
	"github.com/stretchr/testify/assert"
	"necron.dev/zed"
	"testing"
	"time"
)

type OASTest struct {
	Username        string    `zed:"username,err:'username is required',pattern:'^[a-zA-Z0-9_]{6,20}$',pattern_err:'invalid username',min_len:6,max_len:20,min_len_err:'username is too short',max_len_err:'username is too long',required"`
	Password        string    `zed:"password,err:'password is required',min_len:6,max_len:20,min_len_err:'password is too short',max_len_err:'password is too long',required"`
	Email           string    `zed:"email,err:'email is required',pattern:'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',pattern_err:'invalid email',required"`
	RememberSession bool      `zed:"rememberSession,err:'invalid rememberSession'"`
	Age             int       `zed:"age,err:'age is required',min:0,exclusive_min,max:100,min_err:'invalid age',max_err:'invalid age'"`
	Birthdate       time.Time `zed:"birthdate,err:'birthdate is required',required"`
}

func TestOASStruct(t *testing.T) {
	f, e := zed.StructForE[OASTest]("expected struct")
	assert.NoError(t, e)
	x, e := json.MarshalIndent(ogen.NewSpec().AddSchema("uwu", f.ToSchema()), "", "\t")
	fmt.Println(string(x), e)
}
