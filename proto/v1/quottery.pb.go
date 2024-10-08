// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.23.4
// source: quottery.proto

package qubicpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BetInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                      uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatorId               string                 `protobuf:"bytes,2,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	Description             string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Options                 []*BetInfo_Option      `protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty"`
	Oracles                 []*BetInfo_Oracle      `protobuf:"bytes,5,rep,name=oracles,proto3" json:"oracles,omitempty"`
	Votes                   []*BetInfo_Vote        `protobuf:"bytes,6,rep,name=votes,proto3" json:"votes,omitempty"`
	MinimumBetAmount        uint64                 `protobuf:"varint,7,opt,name=minimum_bet_amount,json=minimumBetAmount,proto3" json:"minimum_bet_amount,omitempty"`
	MaximumBetSlotPerOption uint32                 `protobuf:"varint,8,opt,name=maximum_bet_slot_per_option,json=maximumBetSlotPerOption,proto3" json:"maximum_bet_slot_per_option,omitempty"`
	OpenTime                *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=open_time,json=openTime,proto3" json:"open_time,omitempty"`
	CloseTime               *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=close_time,json=closeTime,proto3" json:"close_time,omitempty"`
	EndTime                 *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *BetInfo) Reset() {
	*x = BetInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetInfo) ProtoMessage() {}

func (x *BetInfo) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetInfo.ProtoReflect.Descriptor instead.
func (*BetInfo) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{0}
}

func (x *BetInfo) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *BetInfo) GetCreatorId() string {
	if x != nil {
		return x.CreatorId
	}
	return ""
}

func (x *BetInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *BetInfo) GetOptions() []*BetInfo_Option {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *BetInfo) GetOracles() []*BetInfo_Oracle {
	if x != nil {
		return x.Oracles
	}
	return nil
}

func (x *BetInfo) GetVotes() []*BetInfo_Vote {
	if x != nil {
		return x.Votes
	}
	return nil
}

func (x *BetInfo) GetMinimumBetAmount() uint64 {
	if x != nil {
		return x.MinimumBetAmount
	}
	return 0
}

func (x *BetInfo) GetMaximumBetSlotPerOption() uint32 {
	if x != nil {
		return x.MaximumBetSlotPerOption
	}
	return 0
}

func (x *BetInfo) GetOpenTime() *timestamppb.Timestamp {
	if x != nil {
		return x.OpenTime
	}
	return nil
}

func (x *BetInfo) GetCloseTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CloseTime
	}
	return nil
}

func (x *BetInfo) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

type ActiveBets struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActiveBetIds []uint32 `protobuf:"varint,1,rep,packed,name=active_bet_ids,json=activeBetIds,proto3" json:"active_bet_ids,omitempty"`
}

func (x *ActiveBets) Reset() {
	*x = ActiveBets{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActiveBets) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveBets) ProtoMessage() {}

func (x *ActiveBets) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveBets.ProtoReflect.Descriptor instead.
func (*ActiveBets) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{1}
}

func (x *ActiveBets) GetActiveBetIds() []uint32 {
	if x != nil {
		return x.ActiveBetIds
	}
	return nil
}

type BasicInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fees                 *BasicInfo_Fees          `protobuf:"bytes,1,opt,name=fees,proto3" json:"fees,omitempty"`
	MinimumBetSlotAmount uint64                   `protobuf:"varint,2,opt,name=minimum_bet_slot_amount,json=minimumBetSlotAmount,proto3" json:"minimum_bet_slot_amount,omitempty"`
	IssuedBets           uint64                   `protobuf:"varint,3,opt,name=issued_bets,json=issuedBets,proto3" json:"issued_bets,omitempty"`
	MoneyFlowData        *BasicInfo_MoneyFlowData `protobuf:"bytes,4,opt,name=money_flow_data,json=moneyFlowData,proto3" json:"money_flow_data,omitempty"`
	EconomicsData        *BasicInfo_EconomicsData `protobuf:"bytes,5,opt,name=economics_data,json=economicsData,proto3" json:"economics_data,omitempty"`
	GameOperatorId       string                   `protobuf:"bytes,6,opt,name=game_operator_id,json=gameOperatorId,proto3" json:"game_operator_id,omitempty"`
}

func (x *BasicInfo) Reset() {
	*x = BasicInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicInfo) ProtoMessage() {}

func (x *BasicInfo) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicInfo.ProtoReflect.Descriptor instead.
func (*BasicInfo) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{2}
}

func (x *BasicInfo) GetFees() *BasicInfo_Fees {
	if x != nil {
		return x.Fees
	}
	return nil
}

func (x *BasicInfo) GetMinimumBetSlotAmount() uint64 {
	if x != nil {
		return x.MinimumBetSlotAmount
	}
	return 0
}

func (x *BasicInfo) GetIssuedBets() uint64 {
	if x != nil {
		return x.IssuedBets
	}
	return 0
}

func (x *BasicInfo) GetMoneyFlowData() *BasicInfo_MoneyFlowData {
	if x != nil {
		return x.MoneyFlowData
	}
	return nil
}

func (x *BasicInfo) GetEconomicsData() *BasicInfo_EconomicsData {
	if x != nil {
		return x.EconomicsData
	}
	return nil
}

func (x *BasicInfo) GetGameOperatorId() string {
	if x != nil {
		return x.GameOperatorId
	}
	return ""
}

type BetOptionBettors struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BettorIds []string `protobuf:"bytes,1,rep,name=bettor_ids,json=bettorIds,proto3" json:"bettor_ids,omitempty"`
}

func (x *BetOptionBettors) Reset() {
	*x = BetOptionBettors{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetOptionBettors) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetOptionBettors) ProtoMessage() {}

func (x *BetOptionBettors) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetOptionBettors.ProtoReflect.Descriptor instead.
func (*BetOptionBettors) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{3}
}

func (x *BetOptionBettors) GetBettorIds() []string {
	if x != nil {
		return x.BettorIds
	}
	return nil
}

type BetInfo_Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	State       uint32 `protobuf:"varint,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *BetInfo_Option) Reset() {
	*x = BetInfo_Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetInfo_Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetInfo_Option) ProtoMessage() {}

func (x *BetInfo_Option) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetInfo_Option.ProtoReflect.Descriptor instead.
func (*BetInfo_Option) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{0, 0}
}

func (x *BetInfo_Option) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *BetInfo_Option) GetState() uint32 {
	if x != nil {
		return x.State
	}
	return 0
}

type BetInfo_Oracle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FeePercentage float32 `protobuf:"fixed32,2,opt,name=fee_percentage,json=feePercentage,proto3" json:"fee_percentage,omitempty"`
}

func (x *BetInfo_Oracle) Reset() {
	*x = BetInfo_Oracle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetInfo_Oracle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetInfo_Oracle) ProtoMessage() {}

func (x *BetInfo_Oracle) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetInfo_Oracle.ProtoReflect.Descriptor instead.
func (*BetInfo_Oracle) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{0, 1}
}

func (x *BetInfo_Oracle) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BetInfo_Oracle) GetFeePercentage() float32 {
	if x != nil {
		return x.FeePercentage
	}
	return 0
}

type BetInfo_Vote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OracleId  uint32 `protobuf:"varint,1,opt,name=oracle_id,json=oracleId,proto3" json:"oracle_id,omitempty"`
	WonOption uint32 `protobuf:"varint,2,opt,name=won_option,json=wonOption,proto3" json:"won_option,omitempty"`
}

func (x *BetInfo_Vote) Reset() {
	*x = BetInfo_Vote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetInfo_Vote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetInfo_Vote) ProtoMessage() {}

func (x *BetInfo_Vote) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetInfo_Vote.ProtoReflect.Descriptor instead.
func (*BetInfo_Vote) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{0, 2}
}

func (x *BetInfo_Vote) GetOracleId() uint32 {
	if x != nil {
		return x.OracleId
	}
	return 0
}

func (x *BetInfo_Vote) GetWonOption() uint32 {
	if x != nil {
		return x.WonOption
	}
	return 0
}

type BasicInfo_Fees struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlotPerDay   uint64 `protobuf:"varint,1,opt,name=slot_per_day,json=slotPerDay,proto3" json:"slot_per_day,omitempty"`
	GameOperator uint64 `protobuf:"varint,2,opt,name=game_operator,json=gameOperator,proto3" json:"game_operator,omitempty"`
	Shareholder  uint64 `protobuf:"varint,3,opt,name=shareholder,proto3" json:"shareholder,omitempty"`
	Burn         uint64 `protobuf:"varint,4,opt,name=burn,proto3" json:"burn,omitempty"`
}

func (x *BasicInfo_Fees) Reset() {
	*x = BasicInfo_Fees{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicInfo_Fees) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicInfo_Fees) ProtoMessage() {}

func (x *BasicInfo_Fees) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicInfo_Fees.ProtoReflect.Descriptor instead.
func (*BasicInfo_Fees) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{2, 0}
}

func (x *BasicInfo_Fees) GetSlotPerDay() uint64 {
	if x != nil {
		return x.SlotPerDay
	}
	return 0
}

func (x *BasicInfo_Fees) GetGameOperator() uint64 {
	if x != nil {
		return x.GameOperator
	}
	return 0
}

func (x *BasicInfo_Fees) GetShareholder() uint64 {
	if x != nil {
		return x.Shareholder
	}
	return 0
}

func (x *BasicInfo_Fees) GetBurn() uint64 {
	if x != nil {
		return x.Burn
	}
	return 0
}

type BasicInfo_MoneyFlowData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total       uint64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	IssueBet    uint64 `protobuf:"varint,2,opt,name=issue_bet,json=issueBet,proto3" json:"issue_bet,omitempty"`
	JoinBet     uint64 `protobuf:"varint,3,opt,name=join_bet,json=joinBet,proto3" json:"join_bet,omitempty"`
	FinalizeBet uint64 `protobuf:"varint,4,opt,name=finalize_bet,json=finalizeBet,proto3" json:"finalize_bet,omitempty"`
}

func (x *BasicInfo_MoneyFlowData) Reset() {
	*x = BasicInfo_MoneyFlowData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicInfo_MoneyFlowData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicInfo_MoneyFlowData) ProtoMessage() {}

func (x *BasicInfo_MoneyFlowData) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicInfo_MoneyFlowData.ProtoReflect.Descriptor instead.
func (*BasicInfo_MoneyFlowData) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{2, 1}
}

func (x *BasicInfo_MoneyFlowData) GetTotal() uint64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *BasicInfo_MoneyFlowData) GetIssueBet() uint64 {
	if x != nil {
		return x.IssueBet
	}
	return 0
}

func (x *BasicInfo_MoneyFlowData) GetJoinBet() uint64 {
	if x != nil {
		return x.JoinBet
	}
	return 0
}

func (x *BasicInfo_MoneyFlowData) GetFinalizeBet() uint64 {
	if x != nil {
		return x.FinalizeBet
	}
	return 0
}

type BasicInfo_EconomicsData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EarnedAmountShareholder uint64 `protobuf:"varint,1,opt,name=earned_amount_shareholder,json=earnedAmountShareholder,proto3" json:"earned_amount_shareholder,omitempty"`
	PaidAmountShareholder   uint64 `protobuf:"varint,2,opt,name=paid_amount_shareholder,json=paidAmountShareholder,proto3" json:"paid_amount_shareholder,omitempty"`
	EarnedAmountBetWinner   uint64 `protobuf:"varint,3,opt,name=earned_amount_bet_winner,json=earnedAmountBetWinner,proto3" json:"earned_amount_bet_winner,omitempty"`
	DistributedAmount       uint64 `protobuf:"varint,4,opt,name=distributed_amount,json=distributedAmount,proto3" json:"distributed_amount,omitempty"`
	BurnedAmount            uint64 `protobuf:"varint,5,opt,name=burned_amount,json=burnedAmount,proto3" json:"burned_amount,omitempty"`
}

func (x *BasicInfo_EconomicsData) Reset() {
	*x = BasicInfo_EconomicsData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quottery_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicInfo_EconomicsData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicInfo_EconomicsData) ProtoMessage() {}

func (x *BasicInfo_EconomicsData) ProtoReflect() protoreflect.Message {
	mi := &file_quottery_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicInfo_EconomicsData.ProtoReflect.Descriptor instead.
func (*BasicInfo_EconomicsData) Descriptor() ([]byte, []int) {
	return file_quottery_proto_rawDescGZIP(), []int{2, 2}
}

func (x *BasicInfo_EconomicsData) GetEarnedAmountShareholder() uint64 {
	if x != nil {
		return x.EarnedAmountShareholder
	}
	return 0
}

func (x *BasicInfo_EconomicsData) GetPaidAmountShareholder() uint64 {
	if x != nil {
		return x.PaidAmountShareholder
	}
	return 0
}

func (x *BasicInfo_EconomicsData) GetEarnedAmountBetWinner() uint64 {
	if x != nil {
		return x.EarnedAmountBetWinner
	}
	return 0
}

func (x *BasicInfo_EconomicsData) GetDistributedAmount() uint64 {
	if x != nil {
		return x.DistributedAmount
	}
	return 0
}

func (x *BasicInfo_EconomicsData) GetBurnedAmount() uint64 {
	if x != nil {
		return x.BurnedAmount
	}
	return 0
}

var File_quottery_proto protoreflect.FileDescriptor

var file_quottery_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x71, 0x75, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x05, 0x0a, 0x07,
	0x42, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x71, 0x75, 0x62, 0x69,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x32, 0x0a, 0x07,
	0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x2e, 0x4f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x52, 0x07, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x73,
	0x12, 0x2c, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x65, 0x74, 0x49, 0x6e,
	0x66, 0x6f, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x2c,
	0x0a, 0x12, 0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x62, 0x65, 0x74, 0x5f, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x6d, 0x69, 0x6e, 0x69,
	0x6d, 0x75, 0x6d, 0x42, 0x65, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x1b,
	0x6d, 0x61, 0x78, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x62, 0x65, 0x74, 0x5f, 0x73, 0x6c, 0x6f, 0x74,
	0x5f, 0x70, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x17, 0x6d, 0x61, 0x78, 0x69, 0x6d, 0x75, 0x6d, 0x42, 0x65, 0x74, 0x53, 0x6c, 0x6f,
	0x74, 0x50, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x09, 0x6f, 0x70,
	0x65, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35,
	0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x1a, 0x40, 0x0a, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x3f, 0x0a, 0x06, 0x4f, 0x72, 0x61, 0x63, 0x6c,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0d, 0x66, 0x65, 0x65, 0x50, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x1a, 0x42, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x6f, 0x72, 0x61, 0x63, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x77, 0x6f, 0x6e, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x77, 0x6f, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x32, 0x0a, 0x0a,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x42, 0x65, 0x74, 0x73, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x5f, 0x62, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0d, 0x52, 0x0c, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x42, 0x65, 0x74, 0x49, 0x64, 0x73,
	0x22, 0xec, 0x06, 0x0a, 0x09, 0x42, 0x61, 0x73, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2c,
	0x0a, 0x04, 0x66, 0x65, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x71,
	0x75, 0x62, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x49, 0x6e, 0x66,
	0x6f, 0x2e, 0x46, 0x65, 0x65, 0x73, 0x52, 0x04, 0x66, 0x65, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x17,
	0x6d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x5f, 0x62, 0x65, 0x74, 0x5f, 0x73, 0x6c, 0x6f, 0x74,
	0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6d,
	0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x42, 0x65, 0x74, 0x53, 0x6c, 0x6f, 0x74, 0x41, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x5f, 0x62, 0x65,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64,
	0x42, 0x65, 0x74, 0x73, 0x12, 0x49, 0x0a, 0x0f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x5f, 0x66, 0x6c,
	0x6f, 0x77, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e,
	0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x49, 0x6e,
	0x66, 0x6f, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x46, 0x6c, 0x6f, 0x77, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x0d, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x46, 0x6c, 0x6f, 0x77, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x48, 0x0a, 0x0e, 0x65, 0x63, 0x6f, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x45, 0x63, 0x6f,
	0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0d, 0x65, 0x63, 0x6f, 0x6e,
	0x6f, 0x6d, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x28, 0x0a, 0x10, 0x67, 0x61, 0x6d,
	0x65, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x67, 0x61, 0x6d, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x1a, 0x83, 0x01, 0x0a, 0x04, 0x46, 0x65, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0c,
	0x73, 0x6c, 0x6f, 0x74, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0a, 0x73, 0x6c, 0x6f, 0x74, 0x50, 0x65, 0x72, 0x44, 0x61, 0x79, 0x12, 0x23,
	0x0a, 0x0d, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x67, 0x61, 0x6d, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x68,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x75, 0x72, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x62, 0x75, 0x72, 0x6e, 0x1a, 0x80, 0x01, 0x0a, 0x0d, 0x4d, 0x6f,
	0x6e, 0x65, 0x79, 0x46, 0x6c, 0x6f, 0x77, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x62, 0x65, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x42, 0x65, 0x74, 0x12, 0x19,
	0x0a, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x62, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x6a, 0x6f, 0x69, 0x6e, 0x42, 0x65, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x69, 0x6e,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0b, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x42, 0x65, 0x74, 0x1a, 0x90, 0x02, 0x0a,
	0x0d, 0x45, 0x63, 0x6f, 0x6e, 0x6f, 0x6d, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3a,
	0x0a, 0x19, 0x65, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x17, 0x65, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x53,
	0x68, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x17, 0x70, 0x61,
	0x69, 0x64, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x68,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x15, 0x70, 0x61, 0x69,
	0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x68, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x12, 0x37, 0x0a, 0x18, 0x65, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x62, 0x65, 0x74, 0x5f, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x15, 0x65, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x41, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x42, 0x65, 0x74, 0x57, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x12, 0x64,
	0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x11, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x75,
	0x72, 0x6e, 0x65, 0x64, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0c, 0x62, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x31, 0x0a, 0x10, 0x42, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x65, 0x74, 0x74,
	0x6f, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x65, 0x74, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x62, 0x65, 0x74, 0x74, 0x6f, 0x72, 0x49,
	0x64, 0x73, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2d, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2f,
	0x76, 0x31, 0x2f, 0x71, 0x75, 0x62, 0x69, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_quottery_proto_rawDescOnce sync.Once
	file_quottery_proto_rawDescData = file_quottery_proto_rawDesc
)

func file_quottery_proto_rawDescGZIP() []byte {
	file_quottery_proto_rawDescOnce.Do(func() {
		file_quottery_proto_rawDescData = protoimpl.X.CompressGZIP(file_quottery_proto_rawDescData)
	})
	return file_quottery_proto_rawDescData
}

var file_quottery_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_quottery_proto_goTypes = []interface{}{
	(*BetInfo)(nil),                 // 0: qubic.v1.BetInfo
	(*ActiveBets)(nil),              // 1: qubic.v1.ActiveBets
	(*BasicInfo)(nil),               // 2: qubic.v1.BasicInfo
	(*BetOptionBettors)(nil),        // 3: qubic.v1.BetOptionBettors
	(*BetInfo_Option)(nil),          // 4: qubic.v1.BetInfo.Option
	(*BetInfo_Oracle)(nil),          // 5: qubic.v1.BetInfo.Oracle
	(*BetInfo_Vote)(nil),            // 6: qubic.v1.BetInfo.Vote
	(*BasicInfo_Fees)(nil),          // 7: qubic.v1.BasicInfo.Fees
	(*BasicInfo_MoneyFlowData)(nil), // 8: qubic.v1.BasicInfo.MoneyFlowData
	(*BasicInfo_EconomicsData)(nil), // 9: qubic.v1.BasicInfo.EconomicsData
	(*timestamppb.Timestamp)(nil),   // 10: google.protobuf.Timestamp
}
var file_quottery_proto_depIdxs = []int32{
	4,  // 0: qubic.v1.BetInfo.options:type_name -> qubic.v1.BetInfo.Option
	5,  // 1: qubic.v1.BetInfo.oracles:type_name -> qubic.v1.BetInfo.Oracle
	6,  // 2: qubic.v1.BetInfo.votes:type_name -> qubic.v1.BetInfo.Vote
	10, // 3: qubic.v1.BetInfo.open_time:type_name -> google.protobuf.Timestamp
	10, // 4: qubic.v1.BetInfo.close_time:type_name -> google.protobuf.Timestamp
	10, // 5: qubic.v1.BetInfo.end_time:type_name -> google.protobuf.Timestamp
	7,  // 6: qubic.v1.BasicInfo.fees:type_name -> qubic.v1.BasicInfo.Fees
	8,  // 7: qubic.v1.BasicInfo.money_flow_data:type_name -> qubic.v1.BasicInfo.MoneyFlowData
	9,  // 8: qubic.v1.BasicInfo.economics_data:type_name -> qubic.v1.BasicInfo.EconomicsData
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_quottery_proto_init() }
func file_quottery_proto_init() {
	if File_quottery_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_quottery_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetInfo); i {
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
		file_quottery_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActiveBets); i {
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
		file_quottery_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicInfo); i {
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
		file_quottery_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetOptionBettors); i {
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
		file_quottery_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetInfo_Option); i {
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
		file_quottery_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetInfo_Oracle); i {
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
		file_quottery_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetInfo_Vote); i {
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
		file_quottery_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicInfo_Fees); i {
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
		file_quottery_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicInfo_MoneyFlowData); i {
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
		file_quottery_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicInfo_EconomicsData); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_quottery_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_quottery_proto_goTypes,
		DependencyIndexes: file_quottery_proto_depIdxs,
		MessageInfos:      file_quottery_proto_msgTypes,
	}.Build()
	File_quottery_proto = out.File
	file_quottery_proto_rawDesc = nil
	file_quottery_proto_goTypes = nil
	file_quottery_proto_depIdxs = nil
}
