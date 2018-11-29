package etcd

type MemberInfo struct {
	IP         string
	ClientPort string
}

var ClusterMembers []MemberInfo
