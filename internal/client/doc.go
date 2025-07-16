// Package client applies options to make http client calls
// to server.
// API contains default http.Client with 20000 connects at one time
// and url to the server
// NewAPI creates new instance of API.
// It fills out default client with 20000 connects at one time to the server
// and input url
// Response combines response body and status code from server
// API.GetAll sends request to get all data from storage
// API:
// resp, err := API.GetAll()
// if err != nil {
//    return err
// }
// fmt.Println(string(resp.Body), resp.StatusCode)
// API.Delete takes param and sends request to delete instance from
// server by this param. Param has to be json format string
// API:
// resp, err := API.Delete(`{"key":"key_value"}`)
// if err != nil {
//    return err
// }
// fmt.Println(string(resp.Body), resp.StatusCode)
// API.GetByID takes param and sends request to get instance from
// server by this param. Param has to be json format string
// API:
// resp, err := API.GetByID(`{"key":"key_value"}`)
// if err != nil {
//    return err
// }
// fmt.Println(string(resp.Body), resp.StatusCode)
// API.AddOrUpdate takes param and sends request to insert or update new instance
// to the server. Param has to be json format string.
// Response.Body is always nil
// API:
// resp, err := API.AddOrUpdate(`{
//    "key":"dasha",
//    "entity": {
//		"misha": 20
//    }
// }`)
// if err != nil {
//    return err
// }
// fmt.Println(string(resp.Body), resp.StatusCode)
package client
