package proxy

import (
	logger "log"
	"strings"
	"time"

	"github.com/zu1k/proxypool/pkg/tool"
)

type Base struct {
	Name        string    `yaml:"name" json:"name" gorm:"index"`
	Server      string    `yaml:"server" json:"server" gorm:"index"`
	Port        int       `yaml:"port" json:"port" gorm:"index"`
	Type        string    `yaml:"type" json:"type" gorm:"index"`
	UDP         bool      `yaml:"udp,omitempty" json:"udp,omitempty"`
	Country     string    `yaml:"country,omitempty" json:"country,omitempty" gorm:"index"`
	Useable     bool      `yaml:"useable,omitempty" json:"useable,omitempty" gorm:"index"`
	Delay       uint16    `yaml:"delay" json:"delay" gorm:"index"`
	From        string    `yaml:"from_source" json:"from_source" gorm:"column:from_source"`
	TestTime    time.Time `yaml:"test_time" json:"test_time" gorm:"column:test_time"`
	Host        string    `yaml:"host" json:"host" gorm:"cloumn:host"`
	LinkOrigin  string    `yaml:"link_origin" json:"link_origin" grom:"cloumn:link_origin"`
	LocalSpeed  float32   `yaml:"local_speed" json:"local_speed" grom:"cloumn:local_speed"`
	RemoteSpeed float32   `yaml:"remote_speed" json:"remote_speed" grom:"cloumn:remote_speed"`
}

var log *logger.Logger

func init() {
	log = tool.Logger

}
func (b *Base) TypeName() string {
	if b.Type == "" {
		return "unknown"
	}
	return b.Type
}

func (b *Base) SetLink(link string) {
	b.LinkOrigin = link
}
func (b *Base) SetFrom(from string) {
	b.From = from
}
func (b *Base) SetDelay(delay uint16) {
	b.Delay = delay
}

func (b *Base) SetName(name string) {
	b.Name = name
}

func (b *Base) SetIP(ip string) {
	b.Server = ip
}

// SetHost 添加域名
func (b *Base) SetHost(host string) {
	b.Host = host
}
func (b *Base) BaseInfo() *Base {
	return b
}

func (b *Base) Clone() Base {
	c := *b
	return c
}

func (b *Base) SetUseable(useable bool) {
	b.Useable = useable
}

func (b *Base) SetCountry(country string) {
	b.Country = country
}

type Proxy interface {
	String() string
	ToClash() string
	ToSurge() string
	Link() string
	Identifier() string
	SetName(name string)
	SetIP(ip string)
	SetHost(host string)
	TypeName() string
	BaseInfo() *Base
	Clone() Proxy
	SetUseable(useable bool)
	SetCountry(country string)
	SetDelay(delay uint16)
}

func ParseProxyFromLink(link string) (p Proxy, err error) {
	if strings.HasPrefix(link, "ssr://") {
		p, err = ParseSSRLink(link)
	} else if strings.HasPrefix(link, "vmess://") {
		p, err = ParseVmessLink(link)
	} else if strings.HasPrefix(link, "ss://") {
		p, err = ParseSSLink(link)
	} else if strings.HasPrefix(link, "trojan://") {
		p, err = ParseTrojanLink(link)
	}
	if err != nil || p == nil {
		return nil, err //errors.New("link parse failed")
	}
	ip, country, err := geoIp.Find(p.BaseInfo().Server)
	if err != nil {
		country = "🏁 ZZ"
	}
	p.SetCountry(country)
	p.BaseInfo().SetLink(link)
	// trojan依赖域名？
	if p.TypeName() != "trojan" {
		p.SetHost(p.BaseInfo().Server)
		p.SetIP(ip)
	}
	return
}
