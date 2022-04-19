// Package handlers implements taking http request and sending responses to client
// in JSON format.
// Handlers package provides Storager interface - that is 4 functions of
// simple and common operation in storage GetAll, Get, Put, Delete.
//
//
// There is a Storage struct to connect 2 values logger and storage.
// Logger is a log.Logger. Storage implements Storager interface
// in this struct there are a lot of methods:
// Storage.ServeHTTP that handlers all http request and call specific
// method by input URL and http method. Storage.ServeHTTP checks input URL for correctness.
// If URL does not have /api prefix - Storage.ServeHTTP sends http.StatusNotFound.
// If URL has more than one nested / - Storage.ServeHTTP sends http.StatusNotAcceptable.
// That it tries to match specific URL to method
// If request method is http.MethodGet and there is only one / - Storage.ServeHTTP
// calls Storage.GetAll method
// If request method is http.MethodGet and there is one / with id after it - Storage.ServeHTTP
// calls Storage.Get method
// If request method is http.MethodPut or http.MethodPost and there is only one / - Storage.ServeHTTP
// calls Storage.Put method
// If request method is http.MethodDelete and there is one / with id after it - Storage.ServeHTTP
// calls Storage.Delete method
//
//
// Storage.GetAll sends http.StatusNotFound and error with "no data in storage"
//in the body if there is not any data in Storage.storage.
// If everything is OK sends array of JSON objects
//
// Storage.Get takes data from storage by key if no data
//in Storage.storage or no data by input key sends http.StatusNotFound
// If everything is OK sends JSON object with response key
//
// Storage.Put takes data from request body if some error appears in
//this action - Storage.Put sends http.StatusInternalServerError and error
// If everything is OK update or create new instance in storage.
//
// Storage.Delete removes data by input key. If no data in
//storage will be sent no data in storage and http.StatusNotFound.
// If no value by this key models.ErrNoSuchKey and http.StatusNotFound will be sent.
// If everything is OK data by input key will be removed.
// Storage.OutHTML takes all data from storage and make html file with all data
// If some errors appears in making template or sending error
// Server log error and sends http.StatusInternalServerError to client
package handlers
