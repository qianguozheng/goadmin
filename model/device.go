package model

import (
	"fmt"
	"time"
)

//Store device info
type Device struct {
	Mac           string `gorm: "primary_key;not null;unique"`
	Id            int    `gorm: "AUTO_INCREMENT"`
	Name          string `gorm: "size:255"`
	Online        bool
	Heartbeat     int64
	CreatedAt     time.Time
	Mode          int `gorm:"default:1"` //Router mode AP：0， Route:1
	Version       string
	LanIp         string //Lan
	LanMask       string
	DhcpEnable    int //Dhcp
	DhcpStartIp   string
	DhcpEndIp     string
	DhcpLeaseTime int
	MultiSsid     bool   //Don't impement now
	RfType        string //2.4G
	RfMode        int
	RfFreq        int
	RfPower       int
	RfType5       string //5.8G
	RfMode5       int
	RfFreq5       int
	RfPower5      int
	Status        int
	Ssid          []Ssid `gorm:"foreignkey:DeviceRefer"`
	MultiWan      int    `gorm:"default:1"`
	Wan           []Wan  `gorm:"foreignkey:DeviceRefer"`
	Qos           []Qos  `gorm:"foreignkey:DeviceRefer"`
	ModelType     int    `gorm:"default:0"` //AR9341,AR9344,AR9531,MT7620A,GW500...
	CloudHost     string //Cloud
	CloudPort     int    `gorm:"default:37001"`
	CloudToken    string
	Md5           string `gorm:"default:'00000000000000000000000000000000'"`
}

//Table ssid
type Ssid struct {
	Port        int
	Name        string
	Name5       string
	Url         string
	Password    string
	DeviceRefer int
}

//Table wan
type Wan struct {
	Port          int
	Mode          int
	FixIp         string
	FixMask       string
	FixGateway    string
	PPPoEAccount  string
	PPPoEPassword string
	PrimaryDns    string
	SecondaryDns  string
	DeviceRefer   int
}

//Table wan_qos
type WanQos struct {
	Port     int
	Up       int
	Down     int
	QosRefer int
}

//Table qos
type Qos struct {
	Id          int `gorm:"AUTO_INCREMENT"`
	UpRate      int
	DownRate    int
	TcpLimit    int
	UdpLimit    int
	WanQos      []WanQos `gorm:"foreignkey:WanQos"`
	DeviceRefer int
}

//Device
func InitDevice() {
	dev := &Device{
		Mac:           "aabbccddeeff",
		Name:          "Test Device",
		Online:        true,
		Heartbeat:     1231242,
		Version:       "3.16",
		LanIp:         "192.168.48.1",
		LanMask:       "255.255.255.0",
		DhcpEnable:    1,
		DhcpStartIp:   "192.168.48.2",
		DhcpEndIp:     "192.168.48.254",
		DhcpLeaseTime: 3600,
		MultiSsid:     true,
		RfType:        "2.4",
		RfMode:        1,
		RfFreq:        3,
		RfPower:       23,
		RfType5:       "5",
		RfMode5:       5,
		RfFreq5:       2,
		RfPower5:      27,
		Ssid: []Ssid{
			Ssid{
				Port:        0,
				Name:        "ssid1",
				Password:    "",
				Url:         "http://a.c",
				DeviceRefer: 1,
			},
		},
		Wan: []Wan{
			Wan{
				Port:          1,
				Mode:          1,
				FixIp:         "hellofixip",
				FixMask:       "worldfixmask",
				FixGateway:    "gwate",
				PPPoEAccount:  "pppoe",
				PPPoEPassword: "pppoepass",
				PrimaryDns:    "114.114.114.114",
				SecondaryDns:  "115.115.115.115",
				DeviceRefer:   2,
			},
		},
		Qos: []Qos{
			Qos{
				UpRate:      200,
				DownRate:    4096,
				TcpLimit:    200,
				UdpLimit:    100,
				DeviceRefer: 1,
			},
		},
		ModelType:  3,
		CloudHost:  "112.74.112.103",
		CloudToken: "helloworld",
	}

	if DB.Find(dev).RowsAffected > 0 {
		fmt.Println("Device alread init")
		return
	}
	DB.Create(dev)
}

func GetDevices() []Device {
	var devs []Device
	DB.Debug().Find(&devs)
	return devs
}

//Wan
func GetWanByDeviceId(id int) []Wan {
	var wan []Wan
	DB.Where("device_refer=?", id).Find(&wan)
	return wan
}
func QueryWan(wan Wan) Wan {
	var w Wan
	DB.Debug().Where("port=? and device_refer=?", wan.Port, wan.DeviceRefer).Find(&w)
	return w
}

func InitWan() {
	wan := &Wan{
		Port:          1,
		Mode:          1,
		FixIp:         "hellofixip",
		FixMask:       "worldfixmask",
		FixGateway:    "gwate",
		PPPoEAccount:  "pppoe",
		PPPoEPassword: "pppoepass",
		PrimaryDns:    "114.114.114.114",
		SecondaryDns:  "115.115.115.115",
		DeviceRefer:   2,
	}
	if DB.Find(wan).RowsAffected > 0 {
		fmt.Println("already insert wan params")
		return
	}
	DB.Create(wan)
}

//Ssid
func GetSsidByDeviceId(id int) []Ssid {
	var ssid []Ssid
	DB.Debug().Where("device_refer=?", id).Find(&ssid)
	return ssid
}
func GetSsidByDeviceIdPort(id, port int) Ssid {
	ssid := Ssid{
		Port:port,
		DeviceRefer:id,
	}
	DB.Debug().Find(&ssid)
	return ssid
}

func deleteSsid(id int) {
	DB.Where("device_refer=?", id).Delete(&Ssid{})
}

func AddSsid(ssid []Ssid) {
	deleteSsid(ssid[0].DeviceRefer)
	DB.Debug().Create(&ssid)
}
func UpdateSsid(ssid Ssid) {
	if DB.Debug().Model(&Ssid{}).Where("port=?", ssid.Port).Update(&ssid).RowsAffected < 1 {
		DB.Debug().Create(&ssid)
	}

}
func GetDeviceById(id int) Device {
	var dev Device
	DB.Debug().Where("id=?", id).Find(&dev)
	return dev
}

func UpdateDeviceById(dev Device) {
	DB.Debug().Model(&Device{}).Update(&dev)
}
func GetDeviceID(mac string) int {
	var device Device
	DB.Where("mac=?", mac).Find(&device)
	return device.Id
}

func InitSsid() {
	ssid := &Ssid{
		Port:        0,
		Name:        "ssid1",
		Name5:       "ssid5",
		Password:    "",
		Url:         "http://a.c",
		DeviceRefer: 2,
	}

	if DB.Find(ssid).RowsAffected > 0 {
		fmt.Println("already insert ssid")
		return
	}
	DB.Create(ssid)
}

///////////////WanQos Operation////////////////////////
func AddWanQos(qosrefer, port int) {
	def := &WanQos{
		Port:     port,
		Up:       10,
		Down:     100,
		QosRefer: qosrefer,
	}
	if DB.Find(def).RowsAffected > 0 {
		fmt.Println("already insert wanqos")
		return
	}
	DB.Debug().Create(def)

}

func UpdateWanQoss(wanQos []WanQos) {
	for _, v := range wanQos {
		UpdateWanQos(v)
	}
}
func UpdateWanQos(wanQos WanQos) {
	if DB.Debug().Model(&WanQos{}).Where("port=? and qos_refer=?", wanQos.Port, wanQos.QosRefer).Update(&wanQos).RowsAffected != 1{
		DB.Debug().Create(&wanQos)
	}
}

func QueryWanQos(refer int) []WanQos {
	var wanqoss []WanQos
	DB.Debug().Where("qos_refer=?", refer).Find(&wanqoss)
	return wanqoss
}

func DeleteWanQos(qosrefer, port int) {
	DB.Debug().Model(&WanQos{}).Where("qos_refer=? and port=?", qosrefer, port).Delete(WanQos{})
}

func InitWanQos() {

	//Add
	AddWanQos(1, 0)
	AddWanQos(1, 1)
	AddWanQos(1, 2)
	AddWanQos(1, 3)
	AddWanQos(1, 4)

	//	//Retrive
	//	x := QueryWanQos(12)

	//	for k, v := range x {
	//		fmt.Println("k,v", k, v)
	//	}

	//	del := WanQos{
	//		Port:     3,
	//		QosRefer: 12,
	//	}
	//	//Update
	//	del.Port = 2
	//	del.Down = 40
	//	del.Up = 5
	//	UpdateWanQos(del)

	//	x = QueryWanQos(12)

	//	for k, v := range x {
	//		fmt.Println(k, v)
	//	}
	//	//Delete
	//	DeleteWanQos(12, 2)

	//	x = QueryWanQos(12)

	//	for k, v := range x {
	//		fmt.Println(k, v)
	//	}

}

///////////////Qos Operation//////////////////////////
func AddQos(refer int) {
	qos := &Qos{
		UpRate:      200,
		DownRate:    4096,
		TcpLimit:    200,
		UdpLimit:    100,
		DeviceRefer: refer,
	}
	if DB.Find(qos).RowsAffected > 0 {
		return
	}

	DB.Debug().Create(qos)

}
func UpdateQos(qos Qos) {
	DB.Debug().Model(&Qos{}).Where("device_refer=?", qos.DeviceRefer).Update(&qos)
}
func QueryQos(refer int) Qos {
	var qos Qos
	DB.Debug().Find(&qos, "device_refer=?", refer)
	return qos
}
func InitQos() {
	AddQos(1)
	//	qos := Qos{
	//		DeviceRefer: 13,
	//		UpRate:      90,
	//		DownRate:    399,
	//	}

	//	q := QueryQos(13)
	//	fmt.Println(q)

	//	UpdateQos(qos)

	//	q = QueryQos(13)
	//	fmt.Println(q)
}

//Store QOS
//Store Trust Domain
//Store Trust Ips
//Store CommCfg
//var DB *gorm.DB

//func InitDB() *gorm.DB {
//	db, err := gorm.Open("sqlite3", "testgorm.db")
//	if err != nil {
//		panic("failed to connect database")
//	}
//	//Migrate the schema
//	db.AutoMigrate(&Device{}, &Qos{}, &WanQos{}, &Wan{}, &Ssid{})

//	return db
//}

func InitAllDeviceConfig() {
	InitDevice()
	InitQos()
	InitSsid()
	InitWan()
	InitWanQos()
}
