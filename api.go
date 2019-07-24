package forest

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/robfig/cron"
	"net/http"
)

type JobAPi struct {
	node *JobNode
	echo *echo.Echo
}

func NewJobAPi(node *JobNode) (api *JobAPi) {

	api = &JobAPi{
		node: node,
	}
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin},
	}))
	e.POST("/job/add", api.AddJob)
	e.POST("/job/edit", api.editJob)
	e.POST("/job/delete", api.deleteJob)
	e.POST("/job/list", api.jobList)
	e.POST("/group/add", api.addGroup)
	e.POST("/group/list", api.groupList)
	e.POST("/node/list", api.nodeList)
	e.POST("/plan/list", api.planList)
	e.POST("/client/list", api.clientList)
	e.POST("/snapshot/list", api.snapshotList)
	e.POST("/snapshot/delete", api.snapshotDelete)
	go func() {
		e.Logger.Fatal(e.Start(node.apiAddress))
	}()
	api.echo = e
	return
}

// add a new job
func (api *JobAPi) AddJob(context echo.Context) (err error) {

	var (
		message string
	)
	jobConf := new(JobConf)
	if err = context.Bind(jobConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if jobConf.Name == "" {
		message = "任务名称不能为空"
		goto ERROR
	}
	if jobConf.Group == "" {
		message = "任务分组不能为空"
		goto ERROR
	}

	if jobConf.Cron == "" {
		message = "任务Cron表达式不能为空"
		goto ERROR
	}

	if _, err = cron.Parse(jobConf.Cron); err != nil {
		message = "非法的Cron表达式"
		goto ERROR
	}

	if jobConf.Target == "" {
		message = "任务Target不能为空"
		goto ERROR
	}

	if jobConf.Status == 0 {
		message = "任务状态不能为空"
		goto ERROR
	}

	if err = api.node.manager.AddJob(jobConf); err != nil {
		message = err.Error()
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConf, Message: "创建成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// edit a job
func (api *JobAPi) editJob(context echo.Context) (err error) {

	var (
		message string
	)
	jobConf := new(JobConf)
	if err = context.Bind(jobConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if jobConf.Id == "" {
		message = "此任务记录不存在"
		goto ERROR
	}
	if jobConf.Name == "" {
		message = "任务名称不能为空"
		goto ERROR
	}
	if jobConf.Group == "" {
		message = "任务分组不能为空"
		goto ERROR
	}

	if jobConf.Cron == "" {
		message = "任务Cron表达式不能为空"
		goto ERROR
	}

	if _, err = cron.Parse(jobConf.Cron); err != nil {
		message = "非法的Cron表达式"
		goto ERROR
	}

	if jobConf.Target == "" {
		message = "任务Target不能为空"
		goto ERROR
	}

	if jobConf.Status == 0 {
		message = "任务状态不能为空"
		goto ERROR
	}

	if err = api.node.manager.editJob(jobConf); err != nil {
		message = err.Error()
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConf, Message: "修改成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// job  list
func (api *JobAPi) jobList(context echo.Context) (err error) {

	var (
		jobConfs []*JobConf
	)

	if jobConfs, err = api.node.manager.jobList(); err != nil {
		return context.JSON(http.StatusOK, Result{Code: -1, Message: err.Error()})
	}
	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConfs, Message: "查询成成功"})

}

// delete a job
func (api *JobAPi) deleteJob(context echo.Context) (err error) {

	var (
		message string
	)
	jobConf := new(JobConf)
	if err = context.Bind(jobConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if jobConf.Id == "" {
		message = "此任务记录不存在"
		goto ERROR
	}

	if err = api.node.manager.deleteJob(jobConf); err != nil {
		message = err.Error()
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConf, Message: "删除成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// add a job group
func (api *JobAPi) addGroup(context echo.Context) (err error) {

	var (
		message string
	)
	groupConf := new(GroupConf)
	if err = context.Bind(groupConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if groupConf.Name == "" {
		message = "任务集群名称不能为空"
		goto ERROR
	}

	if groupConf.Remark == "" {
		message = "任务集群描述"
		goto ERROR
	}

	if err = api.node.manager.addGroup(groupConf); err != nil {
		message = err.Error()
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: groupConf, Message: "添加成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// job group list
func (api *JobAPi) groupList(context echo.Context) (err error) {

	var (
		groupConfs []*GroupConf
	)

	if groupConfs, err = api.node.manager.groupList(); err != nil {
		return context.JSON(http.StatusOK, Result{Code: -1, Message: err.Error()})
	}
	return context.JSON(http.StatusOK, Result{Code: 0, Data: groupConfs, Message: "查询成成功"})

}

// job node list
func (api *JobAPi) nodeList(context echo.Context) (err error) {

	var (
		nodes     []*Node
		leader    []byte
		nodeNames []string
	)

	if nodeNames, err = api.node.manager.nodeList(); err != nil {
		return context.JSON(http.StatusOK, Result{Code: -1, Message: err.Error()})
	}

	if leader, err = api.node.etcd.Get(JobNodeElectPath); err != nil {
		return context.JSON(http.StatusOK, Result{Code: -1, Message: err.Error()})
	}

	if len(nodeNames) == 0 {
		return context.JSON(http.StatusOK, Result{Code: 0, Data: nodes, Message: "查询成成功"})
	}

	nodes = make([]*Node, 0)

	for _, name := range nodeNames {

		if name == string(leader) {
			nodes = append(nodes, &Node{Name: name, State: NodeLeaderState})
		} else {
			nodes = append(nodes, &Node{Name: name, State: NodeFollowerState})
		}

	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: nodes, Message: "查询成成功"})

}

func (api *JobAPi) planList(context echo.Context) (err error) {

	var (
		plans []*SchedulePlan
	)
	schedulePlans := api.node.scheduler.schedulePlans
	if len(schedulePlans) == 0 {

		return context.JSON(http.StatusOK, Result{Code: 0, Data: plans})

	}

	plans = make([]*SchedulePlan, 0)

	for _, p := range schedulePlans {

		plans = append(plans, p)

	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: plans})
}

func (api *JobAPi) clientList(context echo.Context) (err error) {

	var (
		query     *QueryClientParam
		message   string
		group     *Group
		clients   []*JobClient
		groupPath string
	)

	query = new(QueryClientParam)
	if err = context.Bind(query); err != nil {
		message = "请选择任务集群"
		goto ERROR
	}

	if query.Group == "" {
		message = "请选择任务集群"
		goto ERROR
	}
	groupPath = fmt.Sprintf("%s%s", GroupConfPath, query.Group)
	if group = api.node.groupManager.groups[groupPath]; group == nil {
		message = "此任务集群不存在"
		goto ERROR
	}

	clients = make([]*JobClient, 0)

	for _, c := range group.clients {

		clients = append(clients, &JobClient{Name: c.name, Path: c.path, Group: query.Group})
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: clients, Message: "查询成成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// 任务快照
func (api *JobAPi) snapshotList(context echo.Context) (err error) {

	var (
		query     *QuerySnapshotParam
		message   string
		keys      [][]byte
		values    [][]byte
		snapshots []*JobSnapshot
		prefix    string
	)

	query = new(QuerySnapshotParam)
	if err = context.Bind(query); err != nil {
		message = "非法的请求参数"
		goto ERROR
	}

	prefix = JobSnapshotPath
	if query.Group != "" && query.Id != "" && query.Ip != "" {
		prefix = fmt.Sprintf(JobClientSnapshotPath, query.Group, query.Ip)
		prefix = fmt.Sprintf("%s/%s", prefix, query.Id)
	} else if query.Group != "" && query.Ip != "" {
		prefix = fmt.Sprintf(JobClientSnapshotPath, query.Group, query.Ip)
	} else if query.Group != "" && query.Ip == "" {
		prefix = fmt.Sprintf(JobSnapshotGroupPath, query.Group)
	}

	if keys, values, err = api.node.etcd.GetWithPrefixKeyLimit(prefix, 500); err != nil {
		message = err.Error()
		goto ERROR
	}

	snapshots = make([]*JobSnapshot, 0)
	if len(keys) == 0 {
		return context.JSON(http.StatusOK, Result{Code: 0, Data: snapshots, Message: "查询成成功"})
	}

	for _, value := range values {

		if len(value) == 0 {
			continue
		}
		var snapshot *JobSnapshot

		if snapshot, err = UParkJobSnapshot(value); err != nil {
			continue
		}

		snapshots = append(snapshots, snapshot)

	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: snapshots, Message: "查询成成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// 任务删除任务快照
func (api *JobAPi) snapshotDelete(context echo.Context) (err error) {

	var (
		query   *QuerySnapshotParam
		message string
		key     string
	)

	query = new(QuerySnapshotParam)
	if err = context.Bind(query); err != nil {
		message = "非法的请求参数"
		goto ERROR
	}

	if query.Group == "" || query.Id == "" || query.Ip == "" {
		message = "非法的请求参数"
		goto ERROR
	}

	key = fmt.Sprintf(JobClientSnapshotPath, query.Group, query.Ip)
	key = fmt.Sprintf("%s/%s", key, query.Id)
	if err = api.node.etcd.Delete(key); err != nil {
		message = err.Error()
		goto ERROR
	}
	return context.JSON(http.StatusOK, Result{Code: 0, Message: "删除成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}
