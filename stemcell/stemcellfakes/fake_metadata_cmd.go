// Code generated by counterfeiter. DO NOT EDIT.
package stemcellfakes

import (
	io "io"
	sync "sync"

	stemcell "github.com/cf-platform-eng/tileinspect/stemcell"
)

type FakeMetadataCmd struct {
	WriteMetadataStub        func(io.Writer) error
	writeMetadataMutex       sync.RWMutex
	writeMetadataArgsForCall []struct {
		arg1 io.Writer
	}
	writeMetadataReturns struct {
		result1 error
	}
	writeMetadataReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMetadataCmd) WriteMetadata(arg1 io.Writer) error {
	fake.writeMetadataMutex.Lock()
	ret, specificReturn := fake.writeMetadataReturnsOnCall[len(fake.writeMetadataArgsForCall)]
	fake.writeMetadataArgsForCall = append(fake.writeMetadataArgsForCall, struct {
		arg1 io.Writer
	}{arg1})
	fake.recordInvocation("WriteMetadata", []interface{}{arg1})
	fake.writeMetadataMutex.Unlock()
	if fake.WriteMetadataStub != nil {
		return fake.WriteMetadataStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.writeMetadataReturns
	return fakeReturns.result1
}

func (fake *FakeMetadataCmd) WriteMetadataCallCount() int {
	fake.writeMetadataMutex.RLock()
	defer fake.writeMetadataMutex.RUnlock()
	return len(fake.writeMetadataArgsForCall)
}

func (fake *FakeMetadataCmd) WriteMetadataCalls(stub func(io.Writer) error) {
	fake.writeMetadataMutex.Lock()
	defer fake.writeMetadataMutex.Unlock()
	fake.WriteMetadataStub = stub
}

func (fake *FakeMetadataCmd) WriteMetadataArgsForCall(i int) io.Writer {
	fake.writeMetadataMutex.RLock()
	defer fake.writeMetadataMutex.RUnlock()
	argsForCall := fake.writeMetadataArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeMetadataCmd) WriteMetadataReturns(result1 error) {
	fake.writeMetadataMutex.Lock()
	defer fake.writeMetadataMutex.Unlock()
	fake.WriteMetadataStub = nil
	fake.writeMetadataReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeMetadataCmd) WriteMetadataReturnsOnCall(i int, result1 error) {
	fake.writeMetadataMutex.Lock()
	defer fake.writeMetadataMutex.Unlock()
	fake.WriteMetadataStub = nil
	if fake.writeMetadataReturnsOnCall == nil {
		fake.writeMetadataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.writeMetadataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeMetadataCmd) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.writeMetadataMutex.RLock()
	defer fake.writeMetadataMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMetadataCmd) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ stemcell.MetadataCmd = new(FakeMetadataCmd)
