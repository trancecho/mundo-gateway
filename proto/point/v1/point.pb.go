// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: point/v1/point.proto

package point

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// 错误码枚举
type ErrorCode int32

const (
	ErrorCode_UNKNOWN_ERROR       ErrorCode = 0
	ErrorCode_POINTS_INSUFFICIENT ErrorCode = 1 // 积分不足
	ErrorCode_OPERATION_FAILED    ErrorCode = 2
	ErrorCode_INVALID_REQUEST     ErrorCode = 3
	ErrorCode_NONE_ERROR          ErrorCode = 4 // 无错误
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "UNKNOWN_ERROR",
		1: "POINTS_INSUFFICIENT",
		2: "OPERATION_FAILED",
		3: "INVALID_REQUEST",
		4: "NONE_ERROR",
	}
	ErrorCode_value = map[string]int32{
		"UNKNOWN_ERROR":       0,
		"POINTS_INSUFFICIENT": 1,
		"OPERATION_FAILED":    2,
		"INVALID_REQUEST":     3,
		"NONE_ERROR":          4,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_point_v1_point_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_point_v1_point_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{0}
}

// 用户信息
type UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId             string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username           string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Experience         int64  `protobuf:"varint,3,opt,name=experience,proto3" json:"experience,omitempty"`                                             // 当前经验值
	Points             int64  `protobuf:"varint,4,opt,name=points,proto3" json:"points,omitempty"`                                                     // 当前积分
	Level              int32  `protobuf:"varint,5,opt,name=level,proto3" json:"level,omitempty"`                                                       // 当前等级
	IsSigned           bool   `protobuf:"varint,6,opt,name=is_signed,json=isSigned,proto3" json:"is_signed,omitempty"`                                 // 是否已签到
	ContinuousSignDays int32  `protobuf:"varint,7,opt,name=continuous_sign_days,json=continuousSignDays,proto3" json:"continuous_sign_days,omitempty"` // 连续签到天数
	TotalSignDays      int32  `protobuf:"varint,8,opt,name=total_sign_days,json=totalSignDays,proto3" json:"total_sign_days,omitempty"`                // 总签到天数
}

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	mi := &file_point_v1_point_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{0}
}

func (x *UserInfo) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserInfo) GetExperience() int64 {
	if x != nil {
		return x.Experience
	}
	return 0
}

func (x *UserInfo) GetPoints() int64 {
	if x != nil {
		return x.Points
	}
	return 0
}

func (x *UserInfo) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *UserInfo) GetIsSigned() bool {
	if x != nil {
		return x.IsSigned
	}
	return false
}

func (x *UserInfo) GetContinuousSignDays() int32 {
	if x != nil {
		return x.ContinuousSignDays
	}
	return 0
}

func (x *UserInfo) GetTotalSignDays() int32 {
	if x != nil {
		return x.TotalSignDays
	}
	return 0
}

// 积分/经验变更请求
type UpdatePointsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId          string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	DeltaPoints     int64  `protobuf:"varint,2,opt,name=delta_points,json=deltaPoints,proto3" json:"delta_points,omitempty"`             // 积分变化量（正加负扣）
	DeltaExperience int64  `protobuf:"varint,3,opt,name=delta_experience,json=deltaExperience,proto3" json:"delta_experience,omitempty"` // 经验变化量
	Reason          string `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`                                           // 变更原因（如"签到"、"发帖"）
}

func (x *UpdatePointsRequest) Reset() {
	*x = UpdatePointsRequest{}
	mi := &file_point_v1_point_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePointsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePointsRequest) ProtoMessage() {}

func (x *UpdatePointsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePointsRequest.ProtoReflect.Descriptor instead.
func (*UpdatePointsRequest) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{1}
}

func (x *UpdatePointsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdatePointsRequest) GetDeltaPoints() int64 {
	if x != nil {
		return x.DeltaPoints
	}
	return 0
}

func (x *UpdatePointsRequest) GetDeltaExperience() int64 {
	if x != nil {
		return x.DeltaExperience
	}
	return 0
}

func (x *UpdatePointsRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// 通用响应
type CommonResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success   bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message   string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	ErrorCode ErrorCode `protobuf:"varint,3,opt,name=error_code,json=errorCode,proto3,enum=mundo.system.point.ErrorCode" json:"error_code,omitempty"`
}

func (x *CommonResponse) Reset() {
	*x = CommonResponse{}
	mi := &file_point_v1_point_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommonResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonResponse) ProtoMessage() {}

func (x *CommonResponse) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonResponse.ProtoReflect.Descriptor instead.
func (*CommonResponse) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{2}
}

func (x *CommonResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CommonResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CommonResponse) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_UNKNOWN_ERROR
}

// 点赞请求
type LikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PostId       string `protobuf:"bytes,2,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	TargetUserId string `protobuf:"bytes,3,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"` // 被点赞的用户ID
}

func (x *LikeRequest) Reset() {
	*x = LikeRequest{}
	mi := &file_point_v1_point_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeRequest) ProtoMessage() {}

func (x *LikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[3]
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
	return file_point_v1_point_proto_rawDescGZIP(), []int{3}
}

func (x *LikeRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *LikeRequest) GetPostId() string {
	if x != nil {
		return x.PostId
	}
	return ""
}

func (x *LikeRequest) GetTargetUserId() string {
	if x != nil {
		return x.TargetUserId
	}
	return ""
}

// 获取用户信息请求
type GetUserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	mi := &file_point_v1_point_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserInfoRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// 签到请求
type SignRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *SignRequest) Reset() {
	*x = SignRequest{}
	mi := &file_point_v1_point_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignRequest) ProtoMessage() {}

func (x *SignRequest) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignRequest.ProtoReflect.Descriptor instead.
func (*SignRequest) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{5}
}

func (x *SignRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// 签到响应
type SignResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success            bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message            string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	ErrorCode          ErrorCode `protobuf:"varint,3,opt,name=error_code,json=errorCode,proto3,enum=mundo.system.point.ErrorCode" json:"error_code,omitempty"`
	Points             int64     `protobuf:"varint,4,opt,name=points,proto3" json:"points,omitempty"`                                                     // 签到获得的积分
	Experience         int64     `protobuf:"varint,5,opt,name=experience,proto3" json:"experience,omitempty"`                                             // 签到获得的经验
	ContinuousSignDays int32     `protobuf:"varint,6,opt,name=continuous_sign_days,json=continuousSignDays,proto3" json:"continuous_sign_days,omitempty"` // 连续签到天数
}

func (x *SignResponse) Reset() {
	*x = SignResponse{}
	mi := &file_point_v1_point_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignResponse) ProtoMessage() {}

func (x *SignResponse) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignResponse.ProtoReflect.Descriptor instead.
func (*SignResponse) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{6}
}

func (x *SignResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SignResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *SignResponse) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_UNKNOWN_ERROR
}

func (x *SignResponse) GetPoints() int64 {
	if x != nil {
		return x.Points
	}
	return 0
}

func (x *SignResponse) GetExperience() int64 {
	if x != nil {
		return x.Experience
	}
	return 0
}

func (x *SignResponse) GetContinuousSignDays() int32 {
	if x != nil {
		return x.ContinuousSignDays
	}
	return 0
}

// 后台统计数据
type AdminStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LevelDistribution []*LevelDistribution `protobuf:"bytes,1,rep,name=level_distribution,json=levelDistribution,proto3" json:"level_distribution,omitempty"`
	AvgPoints         float32              `protobuf:"fixed32,2,opt,name=avg_points,json=avgPoints,proto3" json:"avg_points,omitempty"`
	MonthlyPointsUsed int64                `protobuf:"varint,3,opt,name=monthly_points_used,json=monthlyPointsUsed,proto3" json:"monthly_points_used,omitempty"`
}

func (x *AdminStats) Reset() {
	*x = AdminStats{}
	mi := &file_point_v1_point_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminStats) ProtoMessage() {}

func (x *AdminStats) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminStats.ProtoReflect.Descriptor instead.
func (*AdminStats) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{7}
}

func (x *AdminStats) GetLevelDistribution() []*LevelDistribution {
	if x != nil {
		return x.LevelDistribution
	}
	return nil
}

func (x *AdminStats) GetAvgPoints() float32 {
	if x != nil {
		return x.AvgPoints
	}
	return 0
}

func (x *AdminStats) GetMonthlyPointsUsed() int64 {
	if x != nil {
		return x.MonthlyPointsUsed
	}
	return 0
}

// 等级分布
type LevelDistribution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Level     int32 `protobuf:"varint,1,opt,name=level,proto3" json:"level,omitempty"`
	UserCount int64 `protobuf:"varint,2,opt,name=user_count,json=userCount,proto3" json:"user_count,omitempty"`
}

func (x *LevelDistribution) Reset() {
	*x = LevelDistribution{}
	mi := &file_point_v1_point_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LevelDistribution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LevelDistribution) ProtoMessage() {}

func (x *LevelDistribution) ProtoReflect() protoreflect.Message {
	mi := &file_point_v1_point_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LevelDistribution.ProtoReflect.Descriptor instead.
func (*LevelDistribution) Descriptor() ([]byte, []int) {
	return file_point_v1_point_proto_rawDescGZIP(), []int{8}
}

func (x *LevelDistribution) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *LevelDistribution) GetUserCount() int64 {
	if x != nil {
		return x.UserCount
	}
	return 0
}

var File_point_v1_point_proto protoreflect.FileDescriptor

var file_point_v1_point_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x02, 0x0a, 0x08, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78,
	0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x73,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75,
	0x6f, 0x75, 0x73, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x64, 0x61, 0x79, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x53,
	0x69, 0x67, 0x6e, 0x44, 0x61, 0x79, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x64, 0x61, 0x79, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x6e, 0x44, 0x61, 0x79, 0x73, 0x22,
	0x94, 0x01, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x5f, 0x65, 0x78, 0x70,
	0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x64,
	0x65, 0x6c, 0x74, 0x61, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x82, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3c, 0x0a,
	0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1d, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x65, 0x0a, 0x0b, 0x4c,
	0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x2d, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x26, 0x0a, 0x0b, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0xea, 0x01, 0x0a, 0x0c, 0x53, 0x69,
	0x67, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3c,
	0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f,
	0x75, 0x73, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x64, 0x61, 0x79, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x6f, 0x75, 0x73, 0x53, 0x69,
	0x67, 0x6e, 0x44, 0x61, 0x79, 0x73, 0x22, 0xb1, 0x01, 0x0a, 0x0a, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x54, 0x0a, 0x12, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x5f, 0x64,
	0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x44, 0x69, 0x73, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x44,
	0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x61,
	0x76, 0x67, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x09, 0x61, 0x76, 0x67, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x6c, 0x79, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x75, 0x73, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x6c, 0x79,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x55, 0x73, 0x65, 0x64, 0x22, 0x48, 0x0a, 0x11, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x2a, 0x72, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x4f, 0x49, 0x4e, 0x54, 0x53, 0x5f, 0x49,
	0x4e, 0x53, 0x55, 0x46, 0x46, 0x49, 0x43, 0x49, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x14, 0x0a,
	0x10, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x52,
	0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x4f, 0x4e, 0x45,
	0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x04, 0x32, 0xe8, 0x04, 0x0a, 0x0b, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x64, 0x0a, 0x04, 0x53, 0x69, 0x67, 0x6e,
	0x12, 0x1f, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x22, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a,
	0x22, 0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x12, 0x95,
	0x01, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x41,
	0x6e, 0x64, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x27, 0x2e, 0x6d,
	0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2b, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x25, 0x3a, 0x01, 0x2a, 0x22, 0x20, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x5f, 0x65, 0x78, 0x70, 0x65,
	0x72, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x78, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x26, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x23, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d,
	0x12, 0x6b, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x6b, 0x65, 0x12,
	0x1f, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x22, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a, 0x22,
	0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x69, 0x6b, 0x65, 0x12, 0x74, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x26,
	0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6d, 0x75, 0x6e, 0x64, 0x6f, 0x2e, 0x73,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x2e, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x73, 0x74,
	0x61, 0x74, 0x73, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x63, 0x68, 0x6f, 0x2f, 0x6d, 0x75, 0x6e, 0x64,
	0x6f, 0x2d, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2d, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_point_v1_point_proto_rawDescOnce sync.Once
	file_point_v1_point_proto_rawDescData = file_point_v1_point_proto_rawDesc
)

func file_point_v1_point_proto_rawDescGZIP() []byte {
	file_point_v1_point_proto_rawDescOnce.Do(func() {
		file_point_v1_point_proto_rawDescData = protoimpl.X.CompressGZIP(file_point_v1_point_proto_rawDescData)
	})
	return file_point_v1_point_proto_rawDescData
}

var file_point_v1_point_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_point_v1_point_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_point_v1_point_proto_goTypes = []any{
	(ErrorCode)(0),              // 0: mundo.system.point.ErrorCode
	(*UserInfo)(nil),            // 1: mundo.system.point.UserInfo
	(*UpdatePointsRequest)(nil), // 2: mundo.system.point.UpdatePointsRequest
	(*CommonResponse)(nil),      // 3: mundo.system.point.CommonResponse
	(*LikeRequest)(nil),         // 4: mundo.system.point.LikeRequest
	(*GetUserInfoRequest)(nil),  // 5: mundo.system.point.GetUserInfoRequest
	(*SignRequest)(nil),         // 6: mundo.system.point.SignRequest
	(*SignResponse)(nil),        // 7: mundo.system.point.SignResponse
	(*AdminStats)(nil),          // 8: mundo.system.point.AdminStats
	(*LevelDistribution)(nil),   // 9: mundo.system.point.LevelDistribution
}
var file_point_v1_point_proto_depIdxs = []int32{
	0, // 0: mundo.system.point.CommonResponse.error_code:type_name -> mundo.system.point.ErrorCode
	0, // 1: mundo.system.point.SignResponse.error_code:type_name -> mundo.system.point.ErrorCode
	9, // 2: mundo.system.point.AdminStats.level_distribution:type_name -> mundo.system.point.LevelDistribution
	6, // 3: mundo.system.point.UserService.Sign:input_type -> mundo.system.point.SignRequest
	2, // 4: mundo.system.point.UserService.UpdatePointsAndExperience:input_type -> mundo.system.point.UpdatePointsRequest
	5, // 5: mundo.system.point.UserService.GetUserInfo:input_type -> mundo.system.point.GetUserInfoRequest
	4, // 6: mundo.system.point.UserService.ProcessLike:input_type -> mundo.system.point.LikeRequest
	5, // 7: mundo.system.point.UserService.GetAdminStats:input_type -> mundo.system.point.GetUserInfoRequest
	3, // 8: mundo.system.point.UserService.Sign:output_type -> mundo.system.point.CommonResponse
	3, // 9: mundo.system.point.UserService.UpdatePointsAndExperience:output_type -> mundo.system.point.CommonResponse
	1, // 10: mundo.system.point.UserService.GetUserInfo:output_type -> mundo.system.point.UserInfo
	3, // 11: mundo.system.point.UserService.ProcessLike:output_type -> mundo.system.point.CommonResponse
	8, // 12: mundo.system.point.UserService.GetAdminStats:output_type -> mundo.system.point.AdminStats
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_point_v1_point_proto_init() }
func file_point_v1_point_proto_init() {
	if File_point_v1_point_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_point_v1_point_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_point_v1_point_proto_goTypes,
		DependencyIndexes: file_point_v1_point_proto_depIdxs,
		EnumInfos:         file_point_v1_point_proto_enumTypes,
		MessageInfos:      file_point_v1_point_proto_msgTypes,
	}.Build()
	File_point_v1_point_proto = out.File
	file_point_v1_point_proto_rawDesc = nil
	file_point_v1_point_proto_goTypes = nil
	file_point_v1_point_proto_depIdxs = nil
}
