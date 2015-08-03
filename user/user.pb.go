// Code generated by protoc-gen-go.
// source: user.proto
// DO NOT EDIT!

/*
Package user is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	User
	UserList
	Delta
*/
package user

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type User struct {
	ID   string `protobuf:"bytes,1,opt" json:"ID,omitempty"`
	Name string `protobuf:"bytes,2,opt" json:"Name,omitempty"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}

type UserList struct {
	Users []*User `protobuf:"bytes,1,rep" json:"Users,omitempty"`
}

func (m *UserList) Reset()         { *m = UserList{} }
func (m *UserList) String() string { return proto.CompactTextString(m) }
func (*UserList) ProtoMessage()    {}

func (m *UserList) GetUsers() []*User {
	if m != nil {
		return m.Users
	}
	return nil
}

type Delta struct {
	User string `protobuf:"bytes,1,opt" json:"User,omitempty"`
	Name string `protobuf:"bytes,2,opt" json:"Name,omitempty"`
}

func (m *Delta) Reset()         { *m = Delta{} }
func (m *Delta) String() string { return proto.CompactTextString(m) }
func (*Delta) ProtoMessage()    {}

// Client API for UserQuery service

type UserQueryClient interface {
	Search(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserList, error)
	Update(ctx context.Context, in *Delta, opts ...grpc.CallOption) (*User, error)
	Delete(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	Create(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
}

type userQueryClient struct {
	cc *grpc.ClientConn
}

func NewUserQueryClient(cc *grpc.ClientConn) UserQueryClient {
	return &userQueryClient{cc}
}

func (c *userQueryClient) Search(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserList, error) {
	out := new(UserList)
	err := grpc.Invoke(ctx, "/user.UserQuery/Search", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userQueryClient) Update(ctx context.Context, in *Delta, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserQuery/Update", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userQueryClient) Delete(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserQuery/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userQueryClient) Create(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserQuery/Create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserQuery service

type UserQueryServer interface {
	Search(context.Context, *User) (*UserList, error)
	Update(context.Context, *Delta) (*User, error)
	Delete(context.Context, *User) (*User, error)
	Create(context.Context, *User) (*User, error)
}

func RegisterUserQueryServer(s *grpc.Server, srv UserQueryServer) {
	s.RegisterService(&_UserQuery_serviceDesc, srv)
}

func _UserQuery_Search_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(User)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(UserQueryServer).Search(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _UserQuery_Update_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(Delta)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(UserQueryServer).Update(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _UserQuery_Delete_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(User)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(UserQueryServer).Delete(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _UserQuery_Create_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(User)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(UserQueryServer).Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _UserQuery_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserQuery",
	HandlerType: (*UserQueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _UserQuery_Search_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _UserQuery_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserQuery_Delete_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _UserQuery_Create_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
