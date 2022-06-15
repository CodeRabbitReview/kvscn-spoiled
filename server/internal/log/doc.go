// Package zlog inits logr.Logger instance
// logger will save logs into file storage.json and send them in console
// API calls:
// Log.WithName("test").Error(err, "error type") - send error log with test name
// Log.WithName("test").Info("action is done") - send info log with test name
package zlog
