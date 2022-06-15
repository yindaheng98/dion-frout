package aliyun

import (
	"github.com/pion/ion/proto/ion"
	"github.com/yindaheng98/dion/config"
	pb "github.com/yindaheng98/dion/proto"
)

type StupidAlgorithm struct {
}

const UserDirect = "direct"
const UserPath = "path"

func (s StupidAlgorithm) UpdateSFUStatus(current []*pb.SFUStatus, reports []*pb.QualityReport) (expected []*pb.SFUStatus) {
	expected = current
	for _, s := range expected {
		for _, c := range s.Clients {
			// 以user为标记
			if c.User == UserDirect { // 如果用户需要直连
				makeDirect(s, c.Session) // 就给用户直连
			} else if c.User == UserPath { // 如果用户需要构造路径
				makePath(expected, s.SFU.Nid, c.Session) // 就给用户构造路径
			}
		}
	}
	return
}

// makeDirect 用于构造直连
func makeDirect(s *pb.SFUStatus, session string) {
	addTrack(s, config.ServiceNameStupid, config.ServiceStupid, session)
}

const ServiceNameBeijing = "beijing"
const ServiceNameQingdao = "qingdao"
const ServiceNameNanjing = "nanjing"

var order = []string{
	ServiceNameBeijing,
	ServiceNameQingdao,
	ServiceNameNanjing,
}

// makePath 用于构造路径
func makePath(ss []*pb.SFUStatus, to, session string) {
	i := 0
	for i = range order {
		if to == order[i] { // 先看看这是路径上的第几个
			break
		}
	}
	for _, s := range ss { // 遍历修改所有节点以形成路径
		if s.SFU.Nid == order[0] { // 路径上的第一个要从stupid里取视频
			addTrack(s, config.ServiceNameStupid, config.ServiceStupid, session)
		} else {
			for j := 1; j < i; j++ {
				if s.SFU.Nid == order[j] { // 路径上的后一个从前一个里取视频
					addTrack(s, order[j-1], config.ServiceSXU, session)
				}
			}
		}
	}
}

// addTrack 将某个track添加到ForwardTracks里
func addTrack(s *pb.SFUStatus, nid, service, session string) {
	for _, t := range s.ForwardTracks {
		if t.Src.Service == service && t.Src.Nid == nid {
			t.RemoteSessionId = session
			t.LocalSessionId = session
			return
		}
	}
	s.ForwardTracks = append(s.ForwardTracks, &pb.ForwardTrack{
		Src: &ion.Node{
			Dc:      "dc1",
			Nid:     nid,
			Service: service,
		},
		RemoteSessionId: session,
		LocalSessionId:  session,
	})
}
