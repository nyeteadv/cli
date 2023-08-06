// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package ssh

import (
	context "context"
	sync "sync"
)

// Ensure, that SshServerHostServerMock does implement SshServerHostServer.
// If this is not the case, regenerate this file with moq.
var _ SshServerHostServer = &SshServerHostServerMock{}

// SshServerHostServerMock is a mock implementation of SshServerHostServer.
//
//	func TestSomethingThatUsesSshServerHostServer(t *testing.T) {
//
//		// make and configure a mocked SshServerHostServer
//		mockedSshServerHostServer := &SshServerHostServerMock{
//			StartRemoteServerAsyncFunc: func(contextMoqParam context.Context, startRemoteServerRequest *StartRemoteServerRequest) (*StartRemoteServerResponse, error) {
//				panic("mock out the StartRemoteServerAsync method")
//			},
//			mustEmbedUnimplementedSshServerHostServerFunc: func()  {
//				panic("mock out the mustEmbedUnimplementedSshServerHostServer method")
//			},
//		}
//
//		// use mockedSshServerHostServer in code that requires SshServerHostServer
//		// and then make assertions.
//
//	}
type SshServerHostServerMock struct {
	// StartRemoteServerAsyncFunc mocks the StartRemoteServerAsync method.
	StartRemoteServerAsyncFunc func(contextMoqParam context.Context, startRemoteServerRequest *StartRemoteServerRequest) (*StartRemoteServerResponse, error)

	// mustEmbedUnimplementedSshServerHostServerFunc mocks the mustEmbedUnimplementedSshServerHostServer method.
	mustEmbedUnimplementedSshServerHostServerFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// StartRemoteServerAsync holds details about calls to the StartRemoteServerAsync method.
		StartRemoteServerAsync []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// StartRemoteServerRequest is the startRemoteServerRequest argument value.
			StartRemoteServerRequest *StartRemoteServerRequest
		}
		// mustEmbedUnimplementedSshServerHostServer holds details about calls to the mustEmbedUnimplementedSshServerHostServer method.
		mustEmbedUnimplementedSshServerHostServer []struct {
		}
	}
	lockStartRemoteServerAsync                    sync.RWMutex
	lockmustEmbedUnimplementedSshServerHostServer sync.RWMutex
}

// StartRemoteServerAsync calls StartRemoteServerAsyncFunc.
func (mock *SshServerHostServerMock) StartRemoteServerAsync(contextMoqParam context.Context, startRemoteServerRequest *StartRemoteServerRequest) (*StartRemoteServerResponse, error) {
	if mock.StartRemoteServerAsyncFunc == nil {
		panic("SshServerHostServerMock.StartRemoteServerAsyncFunc: method is nil but SshServerHostServer.StartRemoteServerAsync was just called")
	}
	callInfo := struct {
		ContextMoqParam          context.Context
		StartRemoteServerRequest *StartRemoteServerRequest
	}{
		ContextMoqParam:          contextMoqParam,
		StartRemoteServerRequest: startRemoteServerRequest,
	}
	mock.lockStartRemoteServerAsync.Lock()
	mock.calls.StartRemoteServerAsync = append(mock.calls.StartRemoteServerAsync, callInfo)
	mock.lockStartRemoteServerAsync.Unlock()
	return mock.StartRemoteServerAsyncFunc(contextMoqParam, startRemoteServerRequest)
}

// StartRemoteServerAsyncCalls gets all the calls that were made to StartRemoteServerAsync.
// Check the length with:
//
//	len(mockedSshServerHostServer.StartRemoteServerAsyncCalls())
func (mock *SshServerHostServerMock) StartRemoteServerAsyncCalls() []struct {
	ContextMoqParam          context.Context
	StartRemoteServerRequest *StartRemoteServerRequest
} {
	var calls []struct {
		ContextMoqParam          context.Context
		StartRemoteServerRequest *StartRemoteServerRequest
	}
	mock.lockStartRemoteServerAsync.RLock()
	calls = mock.calls.StartRemoteServerAsync
	mock.lockStartRemoteServerAsync.RUnlock()
	return calls
}

// mustEmbedUnimplementedSshServerHostServer calls mustEmbedUnimplementedSshServerHostServerFunc.
func (mock *SshServerHostServerMock) mustEmbedUnimplementedSshServerHostServer() {
	if mock.mustEmbedUnimplementedSshServerHostServerFunc == nil {
		panic("SshServerHostServerMock.mustEmbedUnimplementedSshServerHostServerFunc: method is nil but SshServerHostServer.mustEmbedUnimplementedSshServerHostServer was just called")
	}
	callInfo := struct {
	}{}
	mock.lockmustEmbedUnimplementedSshServerHostServer.Lock()
	mock.calls.mustEmbedUnimplementedSshServerHostServer = append(mock.calls.mustEmbedUnimplementedSshServerHostServer, callInfo)
	mock.lockmustEmbedUnimplementedSshServerHostServer.Unlock()
	mock.mustEmbedUnimplementedSshServerHostServerFunc()
}

// mustEmbedUnimplementedSshServerHostServerCalls gets all the calls that were made to mustEmbedUnimplementedSshServerHostServer.
// Check the length with:
//
//	len(mockedSshServerHostServer.mustEmbedUnimplementedSshServerHostServerCalls())
func (mock *SshServerHostServerMock) mustEmbedUnimplementedSshServerHostServerCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockmustEmbedUnimplementedSshServerHostServer.RLock()
	calls = mock.calls.mustEmbedUnimplementedSshServerHostServer
	mock.lockmustEmbedUnimplementedSshServerHostServer.RUnlock()
	return calls
}
