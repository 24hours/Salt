package salt

import(
	"testing"
	"fmt"
	"encoding/json"
	"net/url"
	"strings"
)

func TestMain(t *testing.T){
	// make sure the session is sane 
	v := Session()
	if v == nil{
		t.Error(" Session() return NIL")
	}else if v.Request == nil{
		t.Error("Quest.Request is NIL")
	}else if v.Response != nil{
		t.Error("Quest.Response is NOT NIL")
	}


	s, _ := Get("http://httpbin.org")
	if s.Response.StatusCode != 200 {
		t.Error("httpbin.org error expect code 200 , get ", s.Response.StatusCode)
	}
	s, _ = Get("http://httpbin.org/status/202")
	if s.Response.StatusCode != 202 {
		t.Error("httpbin.org error expect code 202 , get ", s.Response.StatusCode)
	}

	s, _ = Custom("CuStoM", "http://httpbin.org")
	// actually this function meant to submit some malformed method
	// a properly implemented server should return 405 
	if s.Response.StatusCode != 405{
		t.Error("Expect code 405, got " , s.Response.StatusCode)
	}

	s, _ = Get("http://httpbin.org/delay/3")
	// this test simple to check that response time work properly
	// probbably meaningless by itself
	if s.ResponseTime < 3{
		t.Error("Expect response time 3s, got " , s.ResponseTime)
	}

	fmt.Println()
}

func TestGet(t *testing.T){
	type items struct{
		Item1 string
		Item2 string
	}
	type resp1 struct{
		Args items
	} 

	var a resp1

	s, _ := Get("http://httpbin.org/get?item1=1&item2=2")

	err := json.Unmarshal(s.Raw, &a)
	if err != nil{
		t.Error("Error occur on decoding JSON String")
	}

	if !strings.EqualFold("1", a.Args.Item1) || !strings.EqualFold("2",a.Args.Item2) {
		t.Error("GET param submitted Properly  (Item1, Item2) expect (1,2) get (", a.Args.Item1, ", " , a.Args.Item2,")")
		t.Error("Requests url : ", s.Request.URL.String())
	}

	b := Session()
	b.Query = map[string][]string{"item1":{"2"}, "item2":{"2"}}
	resp, _ := b.Get("http://httpbin.org/get")

	err = json.Unmarshal(resp.Raw, &a)
	if err != nil{
		t.Error("Error occur on decoding JSON String")
	}

	if !strings.EqualFold("2", a.Args.Item1) || !strings.EqualFold("2",a.Args.Item2) {
		t.Error("GET param submitted Properly  (Item1, Item2) expect (2,2) get (", a.Args.Item1, ", " , a.Args.Item2,")")
	}

	b = Session()
	b.Query = map[string][]string{"item1":{"3"}}
	resp, _ = b.Get("http://httpbin.org/get?item2=3")

	err = json.Unmarshal(resp.Raw, &a)
	if err != nil{
		t.Error("Error occur on decoding JSON String")
	}
	if !strings.EqualFold("3", a.Args.Item1) || !strings.EqualFold("3",a.Args.Item2) {
		t.Error("GET param submitted Properly  (Item1, Item2) expect (3,3) get (", a.Args.Item1, ", " , a.Args.Item2,")")
	}
}

func TestPost(t *testing.T){
	var resp interface{}
	tester := Session()
	s, _ := tester.Post("http://httpbin.org/post", url.Values{"item1" :{"1"}, "item2" :{"2"}})
	_ = json.Unmarshal(s.Raw, &resp)
	dat := resp.(map[string]interface{})
	form := dat["form"].(map[string]interface{})
	// t.Error(s.Text)
	if !strings.EqualFold("1",form["item1"].(string)) || !strings.EqualFold("2", form["item2"].(string)){
		t.Error("ttpbin.org Post Param expect {1 2} , get", form)
	}

	// Simply check if Data are cleared
	// if we send GET request to /post, the server will return 403 error 
	// if this error did not occur , that mean previous request contaminated currect request
	s, _ = tester.Get("http://httpbin.org/post")
	// t.Error(s.Text)
	if s.Response.StatusCode != 405 {
		t.Error("httpbin.org error expect code 405 , get ", s.Response.StatusCode)
	}


	// s, _ = Get("http://httpbin.org/post")
	// if s.Response.StatusCode != 202 {
	// 	t.Error("httpbin.org Post Param expect {1 2}, get ", s.Response.StatusCode)
	// }

}

// This test are quite unnecessary, since it is testing Golang library http.Header type 
// but in this context, I want to make sure that Header value persist across requests
func TestHeader(t *testing.T){
	var h interface{}
	
	sess := Session()
	sess.Request.Header.Set("User-Agent", "Salt Testing Suit/Golang")
	resp, _ := sess.Get("http://httpbin.org/headers")

	err := json.Unmarshal(resp.Raw, &h)
	if err != nil{
		t.Error("Error occur on decoding JSON String")
	}

	dat := h.(map[string]interface{})
	header := dat["headers"].(map[string]interface{})
	if !strings.EqualFold("Salt Testing Suit/Golang",header["User-Agent"].(string)){
		t.Error("Header not set correctly expecting \"Salt Testing Suit/Golang\" got", header["User-Agent"].(string))
	}

	resp, _ = sess.Get("http://httpbin.org/headers")

	h = nil
	err = json.Unmarshal(resp.Raw, &h)
	if err != nil{
		t.Error("Error occur on decoding JSON String")
	}

	dat = h.(map[string]interface{})
	header2 := dat["headers"].(map[string]interface{})
	if !strings.EqualFold("Salt Testing Suit/Golang",header2["User-Agent"].(string)){
		t.Error("Header should persist correctly expecting \"Salt Testing Suit/Golang\" got", header2["User-Agent"].(string))
	}

	// This part are untestable, The idea is that some server (shared hosting)
	// use Host for virtual host, so Host field should be overwrittable regardless or requested url 
	// sess.Request.Header.Set("Host", "www.google.com")
	// resp, _ = sess.Get("http://httpbin.org/")

	// err = json.Unmarshal(resp.Raw, &h)
	// if err != nil{
	// 	t.Error("Error occur on decoding JSON String")
	// }
	// t.Error(h)
	// dat = h.(map[string]interface{})
	// header2 = dat["headers"].(map[string]interface{})
	// if !strings.EqualFold("www.google.com",header2["Host"].(string)){
	// 	t.Error("Host should be overwritable expecting \"www.google.com\" got", header2["Host"].(string))
	// }

}







