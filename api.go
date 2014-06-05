// Wrapper for http.Request, http.Response and http.Client
// attemp to provide python-requests alike library on golang 
package salt
import(
	"net/url"
)

// Perform Get request to "url"
func Get(url string) (*Salt, error){
	sess := Session()
	_, err := sess.Get(url)
	if err != nil{
		return nil, err
	}
	return sess,nil
}

// Perform Head request to "url"
func Head(url string) (*Salt, error){	
	sess := Session()
	_, err := sess.Head(url)
	if err != nil{
		return nil, err
	}
	return sess,nil
}

// Perform Post request to "url"
func Post(url string, data url.Values) (*Salt, error){	
	sess := Session()
	_, err := sess.Post(url, data)
	if err != nil{
		return nil, err
	}
	return sess,nil
}

// Send Request with custom method 
// User are expected to handle go any error throw by server
func Custom(method string, url string) (*Salt, error){	
	sess := Session()
	_, err := sess.Custom(method, url)
	if err != nil{
		return nil, err
	}
	return sess,nil
}


// Create a new session 
// User can customize header value , cookies 
// the value created will be persistance across requests 
func Session() *Salt{
	return NewQuest()
}