package rpc

type rpcError struct {
	s string
}

//GetExecPluginListError rpc request error
func GetExecPluginListError(text string) error {
	return &rpcError{text}
}

func (e *rpcError) Error() string {
	return e.s
}
