package group

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

func KickGroupMember(c *gin.Context) {
	params := api.KickGroupMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.KickGroupMemberReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}

	log.NewInfo(req.OperationID, "KickGroupMember args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.KickGroupMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	var memberListResp api.KickGroupMemberResp
	memberListResp.ErrMsg = RpcResp.ErrMsg
	memberListResp.ErrCode = RpcResp.ErrCode
	for _, v := range RpcResp.Id2ResultList {
		memberListResp.UserIDResultList = append(memberListResp.UserIDResultList, &api.UserIDResult{UserID: v.UserID, Result: v.Result})
	}
	if len(memberListResp.UserIDResultList) == 0 {
		memberListResp.UserIDResultList = []*api.UserIDResult{}
	}

	log.NewInfo(req.OperationID, "KickGroupMember api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

func GetGroupMembersInfo(c *gin.Context) {
	params := api.GetGroupMembersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupMembersInfoReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		//c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		api.SetErrCodeMsg(c, http.StatusInternalServerError)
		return
	}
	log.NewInfo(req.OperationID, "GetGroupMembersInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetGroupMembersInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupMembersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, Data: RpcResp.MemberList}
	log.NewInfo(req.OperationID, "GetGroupMembersInfo api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

func GetGroupMemberList(c *gin.Context) {
	params := api.GetGroupMemberListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupMemberListReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupMemberList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetGroupMemberList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed, ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupMemberListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, MemberList: RpcResp.MemberList, NextSeq: RpcResp.NextSeq}
	if len(memberListResp.MemberList) == 0 {
		memberListResp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	}

	log.NewInfo(req.OperationID, "GetGroupMemberList api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

func GetGroupAllMemberList(c *gin.Context) {
	params := api.GetGroupAllMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupAllMemberReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupAllMember args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetGroupAllMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupAllMember failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupAllMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, MemberList: RpcResp.MemberList}
	if len(memberListResp.MemberList) == 0 {
		memberListResp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	}

	jsm := &jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}

	if len(memberListResp.MemberList) > 0 {
		s, err := jsm.MarshalToString(memberListResp.MemberList[0])
		//{"GroupID":"7836e478bc43ce1d3b8889cac983f59b","UserID":"openIM001","roleLevel":1,"JoinTime":"0","NickName":"","FaceUrl":"https://oss.com.cn/head","AppMangerLevel":0,"JoinSource":0,"OperatorUserID":"","Ex":"xxx"}
		log.NewDebug(req.OperationID, "MarshalToString ", s, err)
		m := ProtoToMap(memberListResp.MemberList[0], false)
		log.NewDebug(req.OperationID, "mmm ", m)
		memberListResp.Test = m
	}

	log.NewInfo(req.OperationID, "GetGroupAllMember api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{}
	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}

func GetJoinedGroupList(c *gin.Context) {
	params := api.GetJoinedGroupListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetJoinedGroupListReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "GetJoinedGroupList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetJoinedGroupList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetJoinedGroupList failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	GroupListResp := api.GetJoinedGroupListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, GroupInfoList: RpcResp.GroupList}
	c.JSON(http.StatusOK, GroupListResp)
	log.NewInfo(req.OperationID, "GetJoinedGroupList api return ", GroupListResp)
}

func InviteUserToGroup(c *gin.Context) {
	params := api.InviteUserToGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.InviteUserToGroupReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "InviteUserToGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.InviteUserToGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.InviteUserToGroupResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	for _, v := range RpcResp.Id2ResultList {
		resp.UserIDResultList = append(resp.UserIDResultList, api.UserIDResult{UserID: v.UserID, Result: v.Result})
	}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "InviteUserToGroup api return ", resp)
}

func CreateGroup(c *gin.Context) {
	params := api.CreateGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CreateGroupReq{}
	utils.CopyStructFields(req, &params)
	for _, v := range params.MemberList {
		req.InitMemberList = append(req.InitMemberList, &rpc.GroupAddMemberInfo{UserID: v.UserID, RoleLevel: v.RoleLevel})
	}
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "CreateGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.CreateGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.CreateGroupResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	if RpcResp.ErrCode == 0 {
		utils.CopyStructFields(&resp.GroupInfo, RpcResp.GroupInfo)
	}
	log.NewInfo(req.OperationID, "CreateGroup api return ", resp)
	c.JSON(http.StatusOK, resp)
}

//  群主或管理员收到的
func GetGroupApplicationList(c *gin.Context) {
	params := api.GetGroupApplicationListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupApplicationListReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupApplicationList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.GetGroupApplicationList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupApplicationList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.GetGroupApplicationListResp{CommResp: api.CommResp{ErrCode: reply.ErrCode, ErrMsg: reply.ErrMsg}, GroupRequestList: reply.GroupRequestList}
	if len(resp.GroupRequestList) == 0 {
		resp.GroupRequestList = []*open_im_sdk.GroupRequest{}
	}
	log.NewInfo(req.OperationID, "GetGroupApplicationList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetGroupsInfo(c *gin.Context) {
	params := api.GetGroupInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupsInfoReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupsInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetGroupsInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupsInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.GetGroupInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, GroupInfoList: RpcResp.GroupInfoList}
	if len(resp.GroupInfoList) == 0 {
		resp.GroupInfoList = []*open_im_sdk.GroupInfo{}
	}
	log.NewInfo(req.OperationID, "GetGroupsInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

//process application
func ApplicationGroupResponse(c *gin.Context) {
	params := api.ApplicationGroupResponseReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GroupApplicationResponseReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "ApplicationGroupResponse args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.GroupApplicationResponse(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GroupApplicationResponse failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.ApplicationGroupResponseResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "ApplicationGroupResponse api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func JoinGroup(c *gin.Context) {
	params := api.JoinGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.JoinGroupReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.JoinGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "JoinGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}
	log.NewInfo(req.OperationID, "JoinGroup api return", RpcResp.String())
	c.JSON(http.StatusOK, resp)
}

func QuitGroup(c *gin.Context) {
	params := api.QuitGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.QuitGroupReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "QuitGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.QuitGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call quit group rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}
	log.NewInfo(req.OperationID, "QuitGroup api return", RpcResp.String())
	c.JSON(http.StatusOK, resp)
}

func SetGroupInfo(c *gin.Context) {
	params := api.SetGroupInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGroupInfoReq{GroupInfo: &open_im_sdk.GroupInfo{}}
	utils.CopyStructFields(req.GroupInfo, &params.Group)
	req.OperationID = params.OperationID
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "SetGroupInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.SetGroupInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetGroupInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.SetGroupInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "SetGroupInfo api return ", resp)
}

func TransferGroupOwner(c *gin.Context) {
	params := api.TransferGroupOwnerReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.TransferGroupOwnerReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "TransferGroupOwner args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.TransferGroupOwner(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "TransferGroupOwner failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.TransferGroupOwnerResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "TransferGroupOwner api return ", resp)
	c.JSON(http.StatusOK, resp)

}
