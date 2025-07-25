// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: comment.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LikeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	VId           int64                  `protobuf:"varint,2,opt,name=VId,proto3" json:"VId,omitempty"`
	Uid           int64                  `protobuf:"varint,3,opt,name=uid,proto3" json:"uid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LikeRequest) Reset() {
	*x = LikeRequest{}
	mi := &file_comment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeRequest) ProtoMessage() {}

func (x *LikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeRequest.ProtoReflect.Descriptor instead.
func (*LikeRequest) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{0}
}

func (x *LikeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LikeRequest) GetVId() int64 {
	if x != nil {
		return x.VId
	}
	return 0
}

func (x *LikeRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type LikeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LikeResponse) Reset() {
	*x = LikeResponse{}
	mi := &file_comment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LikeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeResponse) ProtoMessage() {}

func (x *LikeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeResponse.ProtoReflect.Descriptor instead.
func (*LikeResponse) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{1}
}

type GetCommentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Vid           int64                  `protobuf:"varint,1,opt,name=Vid,proto3" json:"Vid,omitempty"`
	Page          int64                  `protobuf:"varint,2,opt,name=Page,proto3" json:"Page,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCommentRequest) Reset() {
	*x = GetCommentRequest{}
	mi := &file_comment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentRequest) ProtoMessage() {}

func (x *GetCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentRequest.ProtoReflect.Descriptor instead.
func (*GetCommentRequest) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{2}
}

func (x *GetCommentRequest) GetVid() int64 {
	if x != nil {
		return x.Vid
	}
	return 0
}

func (x *GetCommentRequest) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

type GetCommentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comm          []*Comments            `protobuf:"bytes,1,rep,name=comm,proto3" json:"comm,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCommentResponse) Reset() {
	*x = GetCommentResponse{}
	mi := &file_comment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentResponse) ProtoMessage() {}

func (x *GetCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentResponse.ProtoReflect.Descriptor instead.
func (*GetCommentResponse) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{3}
}

func (x *GetCommentResponse) GetComm() []*Comments {
	if x != nil {
		return x.Comm
	}
	return nil
}

type Comments struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             int64                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	VId            int64                  `protobuf:"varint,2,opt,name=VId,proto3" json:"VId,omitempty"`
	Content        string                 `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
	Ctime          int64                  `protobuf:"varint,4,opt,name=Ctime,proto3" json:"Ctime,omitempty"`
	Like           int64                  `protobuf:"varint,5,opt,name=Like,proto3" json:"Like,omitempty"`
	Dislike        int64                  `protobuf:"varint,6,opt,name=Dislike,proto3" json:"Dislike,omitempty"`
	ReplyCount     int64                  `protobuf:"varint,7,opt,name=ReplyCount,proto3" json:"ReplyCount,omitempty"`
	Picture        string           `protobuf:"bytes,8,opt,name=Picture,proto3" json:"Picture,omitempty"`
	User           *User            `protobuf:"bytes,9,opt,name=User,proto3" json:"User,omitempty"`
	MentionedUsers []*MentionedUser `protobuf:"bytes,10,rep,name=MentionedUsers,proto3" json:"MentionedUsers,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Comments) Reset() {
	*x = Comments{}
	mi := &file_comment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Comments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comments) ProtoMessage() {}

func (x *Comments) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comments.ProtoReflect.Descriptor instead.
func (*Comments) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{4}
}

func (x *Comments) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Comments) GetVId() int64 {
	if x != nil {
		return x.VId
	}
	return 0
}

func (x *Comments) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Comments) GetCtime() int64 {
	if x != nil {
		return x.Ctime
	}
	return 0
}

func (x *Comments) GetLike() int64 {
	if x != nil {
		return x.Like
	}
	return 0
}

func (x *Comments) GetDislike() int64 {
	if x != nil {
		return x.Dislike
	}
	return 0
}

func (x *Comments) GetReplyCount() int64 {
	if x != nil {
		return x.ReplyCount
	}
	return 0
}

func (x *Comments) GetPicture() string {
	if x != nil {
		return x.Picture
	}
	return ""
}

func (x *Comments) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *Comments) GetMentionedUsers() []*MentionedUser {
	if x != nil {
		return x.MentionedUsers
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	AvatarURL     string                 `protobuf:"bytes,3,opt,name=AvatarURL,proto3" json:"AvatarURL,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_comment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{5}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetAvatarURL() string {
	if x != nil {
		return x.AvatarURL
	}
	return ""
}

type SubmitCommentRequest struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	UserId         int64                  `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Vid            int64                  `protobuf:"varint,2,opt,name=Vid,proto3" json:"Vid,omitempty"`
	Content        string                 `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
	UserName       string                 `protobuf:"bytes,4,opt,name=UserName,proto3" json:"UserName,omitempty"`
	Picture        string           `protobuf:"bytes,5,opt,name=Picture,proto3" json:"Picture,omitempty"`
	MentionedUsers []*MentionedUser `protobuf:"bytes,6,rep,name=MentionedUsers,proto3" json:"MentionedUsers,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *SubmitCommentRequest) Reset() {
	*x = SubmitCommentRequest{}
	mi := &file_comment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitCommentRequest) ProtoMessage() {}

func (x *SubmitCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitCommentRequest.ProtoReflect.Descriptor instead.
func (*SubmitCommentRequest) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{6}
}

func (x *SubmitCommentRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SubmitCommentRequest) GetVid() int64 {
	if x != nil {
		return x.Vid
	}
	return 0
}

func (x *SubmitCommentRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SubmitCommentRequest) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *SubmitCommentRequest) GetPicture() string {
	if x != nil {
		return x.Picture
	}
	return ""
}

func (x *SubmitCommentRequest) GetMentionedUsers() []*MentionedUser {
	if x != nil {
		return x.MentionedUsers
	}
	return nil
}

type SubmitCommentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SubmitCommentResponse) Reset() {
	*x = SubmitCommentResponse{}
	mi := &file_comment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmitCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitCommentResponse) ProtoMessage() {}

func (x *SubmitCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitCommentResponse.ProtoReflect.Descriptor instead.
func (*SubmitCommentResponse) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{7}
}

type MentionedUser struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	UserName      string                 `protobuf:"bytes,2,opt,name=UserName,proto3" json:"UserName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MentionedUser) Reset() {
	*x = MentionedUser{}
	mi := &file_comment_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MentionedUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MentionedUser) ProtoMessage() {}

func (x *MentionedUser) ProtoReflect() protoreflect.Message {
	mi := &file_comment_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MentionedUser.ProtoReflect.Descriptor instead.
func (*MentionedUser) Descriptor() ([]byte, []int) {
	return file_comment_proto_rawDescGZIP(), []int{8}
}

func (x *MentionedUser) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *MentionedUser) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

var File_comment_proto protoreflect.FileDescriptor

const file_comment_proto_rawDesc = "" +
	"\n" +
	"\rcomment.proto\x12\acomm.v1\"A\n" +
	"\vLikeRequest\x12\x0e\n" +
	"\x02Id\x18\x01 \x01(\x03R\x02Id\x12\x10\n" +
	"\x03VId\x18\x02 \x01(\x03R\x03VId\x12\x10\n" +
	"\x03uid\x18\x03 \x01(\x03R\x03uid\"\x0e\n" +
	"\fLikeResponse\"9\n" +
	"\x11GetCommentRequest\x12\x10\n" +
	"\x03Vid\x18\x01 \x01(\x03R\x03Vid\x12\x12\n" +
	"\x04Page\x18\x02 \x01(\x03R\x04Page\";\n" +
	"\x12GetCommentResponse\x12%\n" +
	"\x04comm\x18\x01 \x03(\v2\x11.comm.v1.CommentsR\x04comm\"\xa7\x02\n" +
	"\bComments\x12\x0e\n" +
	"\x02Id\x18\x01 \x01(\x03R\x02Id\x12\x10\n" +
	"\x03VId\x18\x02 \x01(\x03R\x03VId\x12\x18\n" +
	"\aContent\x18\x03 \x01(\tR\aContent\x12\x14\n" +
	"\x05Ctime\x18\x04 \x01(\x03R\x05Ctime\x12\x12\n" +
	"\x04Like\x18\x05 \x01(\x03R\x04Like\x12\x18\n" +
	"\aDislike\x18\x06 \x01(\x03R\aDislike\x12\x1e\n" +
	"\n" +
	"ReplyCount\x18\a \x01(\x03R\n" +
	"ReplyCount\x12\x18\n" +
	"\aPicture\x18\b \x01(\tR\aPicture\x12!\n" +
	"\x04User\x18\t \x01(\v2\r.comm.v1.UserR\x04User\x12>\n" +
	"\x0eMentionedUsers\x18\n" +
	" \x03(\v2\x16.comm.v1.MentionedUserR\x0eMentionedUsers\"H\n" +
	"\x04User\x12\x0e\n" +
	"\x02Id\x18\x01 \x01(\x03R\x02Id\x12\x12\n" +
	"\x04Name\x18\x02 \x01(\tR\x04Name\x12\x1c\n" +
	"\tAvatarURL\x18\x03 \x01(\tR\tAvatarURL\"\xd0\x01\n" +
	"\x14SubmitCommentRequest\x12\x16\n" +
	"\x06UserId\x18\x01 \x01(\x03R\x06UserId\x12\x10\n" +
	"\x03Vid\x18\x02 \x01(\x03R\x03Vid\x12\x18\n" +
	"\aContent\x18\x03 \x01(\tR\aContent\x12\x1a\n" +
	"\bUserName\x18\x04 \x01(\tR\bUserName\x12\x18\n" +
	"\aPicture\x18\x05 \x01(\tR\aPicture\x12>\n" +
	"\x0eMentionedUsers\x18\x06 \x03(\v2\x16.comm.v1.MentionedUserR\x0eMentionedUsers\"\x17\n" +
	"\x15SubmitCommentResponse\"C\n" +
	"\rMentionedUser\x12\x16\n" +
	"\x06UserId\x18\x01 \x01(\x03R\x06UserId\x12\x1a\n" +
	"\bUserName\x18\x02 \x01(\tR\bUserName2\xdc\x01\n" +
	"\x0eCommentService\x12N\n" +
	"\rSubmitComment\x12\x1d.comm.v1.SubmitCommentRequest\x1a\x1e.comm.v1.SubmitCommentResponse\x12E\n" +
	"\n" +
	"GetComment\x12\x1a.comm.v1.GetCommentRequest\x1a\x1b.comm.v1.GetCommentResponse\x123\n" +
	"\x04Like\x12\x14.comm.v1.LikeRequest\x1a\x15.comm.v1.LikeResponseB\x10Z\x0ecomm/v1;commv1b\x06proto3"

var (
	file_comment_proto_rawDescOnce sync.Once
	file_comment_proto_rawDescData []byte
)

func file_comment_proto_rawDescGZIP() []byte {
	file_comment_proto_rawDescOnce.Do(func() {
		file_comment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_comment_proto_rawDesc), len(file_comment_proto_rawDesc)))
	})
	return file_comment_proto_rawDescData
}

var file_comment_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_comment_proto_goTypes = []any{
	(*LikeRequest)(nil),           // 0: comm.v1.LikeRequest
	(*LikeResponse)(nil),          // 1: comm.v1.LikeResponse
	(*GetCommentRequest)(nil),     // 2: comm.v1.GetCommentRequest
	(*GetCommentResponse)(nil),    // 3: comm.v1.GetCommentResponse
	(*Comments)(nil),              // 4: comm.v1.Comments
	(*User)(nil),                  // 5: comm.v1.User
	(*SubmitCommentRequest)(nil),  // 6: comm.v1.SubmitCommentRequest
	(*SubmitCommentResponse)(nil), // 7: comm.v1.SubmitCommentResponse
	(*MentionedUser)(nil),         // 8: comm.v1.MentionedUser
}
var file_comment_proto_depIdxs = []int32{
	4, // 0: comm.v1.GetCommentResponse.comm:type_name -> comm.v1.Comments
	5, // 1: comm.v1.Comments.User:type_name -> comm.v1.User
	8, // 2: comm.v1.Comments.MentionedUsers:type_name -> comm.v1.MentionedUser
	8, // 3: comm.v1.SubmitCommentRequest.MentionedUsers:type_name -> comm.v1.MentionedUser
	6, // 4: comm.v1.CommentService.SubmitComment:input_type -> comm.v1.SubmitCommentRequest
	2, // 5: comm.v1.CommentService.GetComment:input_type -> comm.v1.GetCommentRequest
	0, // 6: comm.v1.CommentService.Like:input_type -> comm.v1.LikeRequest
	7, // 7: comm.v1.CommentService.SubmitComment:output_type -> comm.v1.SubmitCommentResponse
	3, // 8: comm.v1.CommentService.GetComment:output_type -> comm.v1.GetCommentResponse
	1, // 9: comm.v1.CommentService.Like:output_type -> comm.v1.LikeResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_comment_proto_init() }
func file_comment_proto_init() {
	if File_comment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_comment_proto_rawDesc), len(file_comment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comment_proto_goTypes,
		DependencyIndexes: file_comment_proto_depIdxs,
		MessageInfos:      file_comment_proto_msgTypes,
	}.Build()
	File_comment_proto = out.File
	file_comment_proto_goTypes = nil
	file_comment_proto_depIdxs = nil
}
