// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: proto/engine_service.proto

package engine_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A program by a user to be executed.
type Program struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id of the user associated with the program and request.
	UserId uint32 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	// The source code of the program to be executed.
	SourceCode string `protobuf:"bytes,2,opt,name=sourceCode,proto3" json:"sourceCode,omitempty"`
	// If given, the user's input to be written to stdin during program execution.
	Input *string `protobuf:"bytes,3,opt,name=input,proto3,oneof" json:"input,omitempty"`
}

func (x *Program) Reset() {
	*x = Program{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Program) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Program) ProtoMessage() {}

func (x *Program) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Program.ProtoReflect.Descriptor instead.
func (*Program) Descriptor() ([]byte, []int) {
	return file_proto_engine_service_proto_rawDescGZIP(), []int{0}
}

func (x *Program) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Program) GetSourceCode() string {
	if x != nil {
		return x.SourceCode
	}
	return ""
}

func (x *Program) GetInput() string {
	if x != nil && x.Input != nil {
		return *x.Input
	}
	return ""
}

// A test case to judge whether a program is valid.
type TestCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Data to test a program against.
	TestData string `protobuf:"bytes,1,opt,name=testData,proto3" json:"testData,omitempty"`
	// The expected output from a valid program for the test data.
	ExpectedOutput string `protobuf:"bytes,2,opt,name=expectedOutput,proto3" json:"expectedOutput,omitempty"`
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_proto_engine_service_proto_rawDescGZIP(), []int{1}
}

func (x *TestCase) GetTestData() string {
	if x != nil {
		return x.TestData
	}
	return ""
}

func (x *TestCase) GetExpectedOutput() string {
	if x != nil {
		return x.ExpectedOutput
	}
	return ""
}

// A program with a list of test cases to be executed against.
type ProgramWithTests struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user's program
	Program *Program `protobuf:"bytes,1,opt,name=program,proto3" json:"program,omitempty"`
	// The list of test cases
	Tests []*TestCase `protobuf:"bytes,2,rep,name=tests,proto3" json:"tests,omitempty"`
}

func (x *ProgramWithTests) Reset() {
	*x = ProgramWithTests{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgramWithTests) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgramWithTests) ProtoMessage() {}

func (x *ProgramWithTests) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgramWithTests.ProtoReflect.Descriptor instead.
func (*ProgramWithTests) Descriptor() ([]byte, []int) {
	return file_proto_engine_service_proto_rawDescGZIP(), []int{2}
}

func (x *ProgramWithTests) GetProgram() *Program {
	if x != nil {
		return x.Program
	}
	return nil
}

func (x *ProgramWithTests) GetTests() []*TestCase {
	if x != nil {
		return x.Tests
	}
	return nil
}

// Contains the result of executing a program.
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The stdout of the program.
	StandardOutput string `protobuf:"bytes,1,opt,name=standardOutput,proto3" json:"standardOutput,omitempty"`
	// The stderr produced by the program.
	// nil if there was no problem with program execution.
	StandardError string `protobuf:"bytes,2,opt,name=standardError,proto3" json:"standardError,omitempty"`
	// The time, in ms, that the program took to terminate execution.
	ElapsedTime string `protobuf:"bytes,3,opt,name=elapsedTime,proto3" json:"elapsedTime,omitempty"`
	// The amount of memory used by the program.
	MemoryUsage string `protobuf:"bytes,4,opt,name=memoryUsage,proto3" json:"memoryUsage,omitempty"`
	// Indicates any problem the engine itself had with executing the program.
	// For example, the program may have timed out or consumed too much memory.
	EngineError string `protobuf:"bytes,5,opt,name=engineError,proto3" json:"engineError,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_proto_engine_service_proto_rawDescGZIP(), []int{3}
}

func (x *Result) GetStandardOutput() string {
	if x != nil {
		return x.StandardOutput
	}
	return ""
}

func (x *Result) GetStandardError() string {
	if x != nil {
		return x.StandardError
	}
	return ""
}

func (x *Result) GetElapsedTime() string {
	if x != nil {
		return x.ElapsedTime
	}
	return ""
}

func (x *Result) GetMemoryUsage() string {
	if x != nil {
		return x.MemoryUsage
	}
	return ""
}

func (x *Result) GetEngineError() string {
	if x != nil {
		return x.EngineError
	}
	return ""
}

// Contains the result of executing a program against a set of test cases.
type TestResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Indicates whether the program successfully passed against all test cases.
	Accepted bool `protobuf:"varint,1,opt,name=accepted,proto3" json:"accepted,omitempty"`
	// The first test case that the program failed on.
	FailedTestIndex *TestCase `protobuf:"bytes,2,opt,name=failedTestIndex,proto3,oneof" json:"failedTestIndex,omitempty"`
	// Result of the program when tested against the first failed test case.
	FailedResult *Result `protobuf:"bytes,3,opt,name=failedResult,proto3,oneof" json:"failedResult,omitempty"`
}

func (x *TestResult) Reset() {
	*x = TestResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_engine_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResult) ProtoMessage() {}

func (x *TestResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_engine_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResult.ProtoReflect.Descriptor instead.
func (*TestResult) Descriptor() ([]byte, []int) {
	return file_proto_engine_service_proto_rawDescGZIP(), []int{4}
}

func (x *TestResult) GetAccepted() bool {
	if x != nil {
		return x.Accepted
	}
	return false
}

func (x *TestResult) GetFailedTestIndex() *TestCase {
	if x != nil {
		return x.FailedTestIndex
	}
	return nil
}

func (x *TestResult) GetFailedResult() *Result {
	if x != nil {
		return x.FailedResult
	}
	return nil
}

var File_proto_engine_service_proto protoreflect.FileDescriptor

var file_proto_engine_service_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x07,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x19, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x22, 0x4e, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x26, 0x0a, 0x0e,
	0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x22, 0x57, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x57,
	0x69, 0x74, 0x68, 0x54, 0x65, 0x73, 0x74, 0x73, 0x12, 0x22, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x67,
	0x72, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x50, 0x72, 0x6f, 0x67,
	0x72, 0x61, 0x6d, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x05,
	0x74, 0x65, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x52, 0x05, 0x74, 0x65, 0x73, 0x74, 0x73, 0x22, 0xbc, 0x01,
	0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x6e,
	0x64, 0x61, 0x72, 0x64, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x72, 0x64, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x72, 0x64, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x72,
	0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x6c, 0x61,
	0x70, 0x73, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x65, 0x6d, 0x6f,
	0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0xb9, 0x01, 0x0a,
	0x0a, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x12, 0x38, 0x0a, 0x0f, 0x66, 0x61, 0x69, 0x6c, 0x65,
	0x64, 0x54, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x09, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x66,
	0x61, 0x69, 0x6c, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x88, 0x01,
	0x01, 0x12, 0x30, 0x0a, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x48, 0x01, 0x52, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x88, 0x01, 0x01, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x54, 0x65,
	0x73, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x66, 0x61, 0x69, 0x6c,
	0x65, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x69, 0x0a, 0x0d, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x08, 0x2e,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x31, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x54, 0x65, 0x73, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x11, 0x2e, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x57, 0x69, 0x74,
	0x68, 0x54, 0x65, 0x73, 0x74, 0x73, 0x1a, 0x0b, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x42, 0x25, 0x5a, 0x23, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_engine_service_proto_rawDescOnce sync.Once
	file_proto_engine_service_proto_rawDescData = file_proto_engine_service_proto_rawDesc
)

func file_proto_engine_service_proto_rawDescGZIP() []byte {
	file_proto_engine_service_proto_rawDescOnce.Do(func() {
		file_proto_engine_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_engine_service_proto_rawDescData)
	})
	return file_proto_engine_service_proto_rawDescData
}

var file_proto_engine_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_engine_service_proto_goTypes = []interface{}{
	(*Program)(nil),          // 0: Program
	(*TestCase)(nil),         // 1: TestCase
	(*ProgramWithTests)(nil), // 2: ProgramWithTests
	(*Result)(nil),           // 3: Result
	(*TestResult)(nil),       // 4: TestResult
}
var file_proto_engine_service_proto_depIdxs = []int32{
	0, // 0: ProgramWithTests.program:type_name -> Program
	1, // 1: ProgramWithTests.tests:type_name -> TestCase
	1, // 2: TestResult.failedTestIndex:type_name -> TestCase
	3, // 3: TestResult.failedResult:type_name -> Result
	0, // 4: EngineService.GetProgramResult:input_type -> Program
	2, // 5: EngineService.GetTestedResult:input_type -> ProgramWithTests
	3, // 6: EngineService.GetProgramResult:output_type -> Result
	4, // 7: EngineService.GetTestedResult:output_type -> TestResult
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_engine_service_proto_init() }
func file_proto_engine_service_proto_init() {
	if File_proto_engine_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_engine_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Program); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_engine_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCase); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_engine_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgramWithTests); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_engine_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_engine_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_engine_service_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_proto_engine_service_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_engine_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_engine_service_proto_goTypes,
		DependencyIndexes: file_proto_engine_service_proto_depIdxs,
		MessageInfos:      file_proto_engine_service_proto_msgTypes,
	}.Build()
	File_proto_engine_service_proto = out.File
	file_proto_engine_service_proto_rawDesc = nil
	file_proto_engine_service_proto_goTypes = nil
	file_proto_engine_service_proto_depIdxs = nil
}
