package service

import (
	"sync"
	"time"

	"gost-panel/internal/model"
	"gost-panel/internal/repository"
	"gost-panel/internal/utils"
	"gost-panel/pkg/logger"

	"gorm.io/gorm"
)

// NodeHealthService 节点健康检测服务
// 使用 Gost API 进行健康检查
type NodeHealthService struct {
	nodeRepo    *repository.NodeRepository
	ruleRepo    *repository.RuleRepository
	tunnelRepo  *repository.TunnelRepository
	ruleSvc     *RuleService
	tunnelSvc   *TunnelService
	resumeEvery time.Duration
	lastResume  sync.Map // map[uint]time.Time
	resumeLocks sync.Map // map[uint]*sync.Mutex
	ticker      *time.Ticker
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// NewNodeHealthService 创建节点健康检测服务
func NewNodeHealthService(db *gorm.DB) *NodeHealthService {
	return &NodeHealthService{
		nodeRepo:    repository.NewNodeRepository(db),
		ruleRepo:    repository.NewRuleRepository(db),
		tunnelRepo:  repository.NewTunnelRepository(db),
		ruleSvc:     NewRuleService(db),
		tunnelSvc:   NewTunnelService(db),
		resumeEvery: 15 * time.Second,
		stopChan:    make(chan struct{}),
	}
}

// Start 启动定时健康检测（每 5 秒）
func (s *NodeHealthService) Start() {
	s.ticker = time.NewTicker(5 * time.Second)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		logger.Info("节点健康检测服务已启动")

		// 立即执行一次
		s.checkAll()

		for {
			select {
			case <-s.ticker.C:
				s.checkAll()
			case <-s.stopChan:
				logger.Info("节点健康检测服务已停止")
				return
			}
		}
	}()
}

// Stop 停止健康检测
func (s *NodeHealthService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
	s.wg.Wait()
}

// checkAll 检测所有资源
func (s *NodeHealthService) checkAll() {
	s.checkNodes()
}

// checkNodes 检测所有节点
// 使用 Gost API 的 /config 接口进行健康检查
func (s *NodeHealthService) checkNodes() {
	nodes, _, err := s.nodeRepo.List(nil)
	if err != nil {
		logger.Errorf("获取节点列表失败: %v", err)
		return
	}

	for _, node := range nodes {
		go func(n model.GostNode) {
			status := s.checkNodeHealth(n)
			wasOnline := n.Status == model.NodeStatusOnline

			// 状态变更处理
			if status != n.Status {
				logger.Infof("节点 %s 状态变更: %s -> %s", n.Name, n.Status, status)
				if err := s.nodeRepo.UpdateStatus(n.ID, status); err != nil {
					logger.Errorf("更新节点 %s 状态失败: %v", n.Name, err)
				}
			}

			if status == model.NodeStatusOnline {
				logger.Debugf("节点 %s 在线", n.Name)
				s.scheduleResumeForNode(n.ID, n.Name, !wasOnline)
			} else {
				// 停止其关联的所有规则和隧道
				_ = s.ruleRepo.StopByNodeID(n.ID)

				// 查找并停止受影响的隧道关联的规则
				if tunnels, err := s.tunnelRepo.FindByNodeID(n.ID); err == nil && len(tunnels) > 0 {
					var tunnelIDs []uint
					for _, t := range tunnels {
						tunnelIDs = append(tunnelIDs, t.ID)
					}
					_ = s.ruleRepo.StopByTunnelIDs(tunnelIDs)
				}

				_ = s.tunnelRepo.StopByNodeID(n.ID)
				logger.Debugf("节点 %s 离线, status=%s, old=%s", n.Name, status, n.Status)
			}

			_ = s.nodeRepo.UpdateLastCheck(n.ID)
		}(node)
	}
}

func (s *NodeHealthService) resumeResourcesForNode(nodeID uint, nodeName string) {
	tunnelsToResume, err := s.tunnelRepo.FindAutoResumeByNodeID(nodeID)
	if err != nil {
		logger.Errorf("查询节点 %s 自动恢复隧道失败: %v", nodeName, err)
		return
	}
	for _, tunnel := range tunnelsToResume {
		if err := s.tunnelSvc.Start(tunnel.ID, 0, "system", "", "auto-resume"); err != nil {
			logger.Warnf("自动恢复隧道失败: tunnel_id=%d name=%s err=%v", tunnel.ID, tunnel.Name, err)
		}
	}

	relatedTunnels, err := s.tunnelRepo.FindByNodeID(nodeID)
	if err != nil {
		logger.Errorf("查询节点 %s 关联隧道失败: %v", nodeName, err)
		return
	}
	tunnelIDs := make([]uint, 0, len(relatedTunnels))
	for _, tunnel := range relatedTunnels {
		tunnelIDs = append(tunnelIDs, tunnel.ID)
	}

	rulesByNode, err := s.ruleRepo.FindAutoResumeByNodeID(nodeID)
	if err != nil {
		logger.Errorf("查询节点 %s 自动恢复规则失败: %v", nodeName, err)
		return
	}
	rulesByTunnel, err := s.ruleRepo.FindAutoResumeByTunnelIDs(tunnelIDs)
	if err != nil {
		logger.Errorf("查询节点 %s 隧道自动恢复规则失败: %v", nodeName, err)
		return
	}

	if len(tunnelsToResume) == 0 && len(rulesByNode) == 0 && len(rulesByTunnel) == 0 {
		return
	}

	logger.Infof("节点 %s 执行自动恢复: tunnels=%d rules(node)=%d rules(tunnel)=%d", nodeName, len(tunnelsToResume), len(rulesByNode), len(rulesByTunnel))

	ruleIDSet := make(map[uint]struct{})
	for _, rule := range rulesByNode {
		ruleIDSet[rule.ID] = struct{}{}
	}
	for _, rule := range rulesByTunnel {
		ruleIDSet[rule.ID] = struct{}{}
	}

	for ruleID := range ruleIDSet {
		if err := s.ruleSvc.Start(ruleID, 0, "system", "", "auto-resume"); err != nil {
			logger.Warnf("自动恢复规则失败: rule_id=%d err=%v", ruleID, err)
		}
	}
}

func (s *NodeHealthService) scheduleResumeForNode(nodeID uint, nodeName string, force bool) {
	if !s.shouldResumeNow(nodeID, force) {
		return
	}

	lock := s.getResumeLock(nodeID)
	lock.Lock()
	defer lock.Unlock()

	if !s.shouldResumeNow(nodeID, force) {
		return
	}
	s.lastResume.Store(nodeID, time.Now())
	s.resumeResourcesForNode(nodeID, nodeName)
}

func (s *NodeHealthService) shouldResumeNow(nodeID uint, force bool) bool {
	if force {
		return true
	}
	lastRaw, ok := s.lastResume.Load(nodeID)
	if !ok {
		return true
	}
	last, ok := lastRaw.(time.Time)
	if !ok {
		return true
	}
	return time.Since(last) >= s.resumeEvery
}

func (s *NodeHealthService) getResumeLock(nodeID uint) *sync.Mutex {
	if lock, ok := s.resumeLocks.Load(nodeID); ok {
		if m, ok := lock.(*sync.Mutex); ok {
			return m
		}
	}
	m := &sync.Mutex{}
	actual, _ := s.resumeLocks.LoadOrStore(nodeID, m)
	return actual.(*sync.Mutex)
}

// checkNodeHealth 检查单个节点的健康状态
// 通过调用 Gost API 的 /config 接口来判断节点是否可用
func (s *NodeHealthService) checkNodeHealth(node model.GostNode) model.NodeStatus {
	// 检查地址是否有效
	if node.Address == "" || node.Port == 0 {
		return model.NodeStatusOffline
	}

	// 验证 Gost API 是否可用
	client := utils.GetGostClient(&node)

	if err := client.HealthCheck(); err != nil {
		logger.Debugf("节点 %d (%s) API 检查失败: %v", node.ID, node.Name, err)
		return model.NodeStatusOffline
	}

	return model.NodeStatusOnline
}
