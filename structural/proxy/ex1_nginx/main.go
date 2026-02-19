package main

type Server interface {
	HandleRequest(url, method string) (int, string)
}

type NginxProxy struct {
	server       Server
	rateLimit    int
	requestCount map[string]int
	cache        map[string]string
}

func (n *NginxProxy) HandleRequest(url, method string) (int, string) {
	if n.requestCount[url] >= n.rateLimit{
		return 429,"Not Allowed"
	}
	n.requestCount[url]++
	if val,ok := n.cache[url] ; ok{
		return 200,val
	} 
	
	val,msg := n.server.HandleRequest(url,method)
	n.cache[url] = msg

	return val,n.cache[url]
}
