/*
Package jsonrpc provides types and functions for transport independent JSONRPC handling

Main entry point is Process(3) function that accepts reader of incoming data, message dispatcher interface and response writer
In case of unexpected panic Process(3) will always return Internal Error with panic message as error message.
Process(3) will call ProcessBatchMessage(3) or ProcessSingleMessage(3) based on first JSON token.

ProcessBatchMessage blocks until all incoming messages will be processed.
*/
package jsonrpc
