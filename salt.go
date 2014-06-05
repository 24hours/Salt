package salt

import(
	"net/http"
	"net/url"
	"io/ioutil"
	"time"
	"strings"
	"io"
	"fmt"
)

type Salt struct{
	client *http.Client
	Request *http.Request
	Response *http.Response
	ResponseTime time.Duration
	Query url.Values
	Data url.Values
	Text string
	Raw []byte
}

func NewQuest() *Salt{
	var err error
	ret := new(Salt)
	ret.client = new(http.Client)
	// NewRequest does some internal initialization 
	// Do not attempt to use new(http.Request)
	ret.Request, err = http.NewRequest("", "", nil)
	if err != nil{
		panic("Fail to create NewRequest, any action are impossible to continue")
	}
	ret.Query = nil
	ret.Response = nil
	ret.Data = nil 
	return ret
}

func (self *Salt) Get(url string) (*Salt, error){
	self.Request.Method = "Get"
	_, err := self.do(url)
	return self, err
}

func (self *Salt) Head(url string) (*Salt, error){
	self.Request.Method = "Head"
	_, err := self.do(url)
	return self, err
}

func (self *Salt) Post(url string, data url.Values) (*Salt, error){
	self.Request.Method = "Post"
	self.Data = data
	_, err := self.do(url)
	return self, err
}

func (self *Salt) Custom(method string, url string) (*Salt, error){
	self.Request.Method = method
	_, err := self.do(url)
	return self, err
}



func (self *Salt) do(urls string) (*http.Response, error){
	u, err := url.Parse(urls)
	self.Request.URL = u
	if self.Query == nil{
		self.Query = url.Values{}
	}
	
	for k, v := range self.Request.URL.Query(){
		self.Query[k] = v
	}
	self.Request.URL.RawQuery = self.Query.Encode()

	//post data 
	if self.Data != nil{
		self.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// [fix] I have a hunch that this part will cause some weird bug 
		body := io.Reader(strings.NewReader(self.Data.Encode()))
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
   			rc = ioutil.NopCloser(body)
   		}
		self.Request.Body = rc
		// [fix] if data exist, then POST method must be used 
		self.Request.Method = "Post"
	}

	// user requested to override Host field in request header 
	// this action should be honored, user are expected to handle the error
	host := self.Request.Header.Get("Host")
	if !strings.EqualFold("", host){
		self.Request.Host = host
		fmt.Println(self.Request.Host)
	}

	//send the request now 
	tic := time.Now()
	ret , err := self.client.Do(self.Request)
	if err != nil{
		self.Response = nil
		return nil, err
	}
	toc := time.Now()
	self.ResponseTime = toc.Sub(tic)

	self.Raw, err = ioutil.ReadAll(ret.Body)
	defer ret.Body.Close()
	if err != nil{
		self.Response = nil
		return nil, err
	}

	self.Text = string(self.Raw)
	self.Request = ret.Request
	self.Data = nil
	self.Query = nil
 	self.Response = ret
	return ret, nil
}








func (req *Salt) Hello(){
	fmt.Println("req.val")
}