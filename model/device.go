package model

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
	Status        int    `gorm:"default:0"`
	Ssid          []Ssid `gorm:"foreignkey:DeviceRefer"`
	MultiWan      int    `gorm:"default:1"`
	Wan           []Wan  `gorm:"foreignkey:DeviceRefer"`
	Qos           []Qos  `gorm:"foreignkey:DeviceRefer"`
	ModelType     int    `gorm:"default:0"` //AR9341,AR9344,AR9531,MT7620A,GW500...
	CloudHost     string //Cloud
	CloudPort     int    `gorm:"default:37001"`
	CloudToken    string
	Md5           string `gorm:"default:'00000000000000000000000000000000'"`
	ProjectRefer  int    `gorm:"default:0"`
	Sync          int    `gorm:"default:0"`
}

//Table ssid
type Ssid struct {
	Port        int    `json:"port"`
	Name        string `json:"name"`
	Name5       string `json:"5gname"`
	Url         string `json:"url"`
	Password    string `json:"password"`
	DeviceRefer int
}

//Table wan
//type Wan struct {
//	Port          int
//	Mode          int
//	FixIp         string
//	FixMask       string
//	FixGateway    string
//	PPPoEAccount  string
//	PPPoEPassword string
//	PrimaryDns    string
//	SecondaryDns  string
//}

type Wan struct {
	Port          int    `json:"port"`
	Mode          string `json:"mode"` //"0"-dhcp,"1"-fix,"2"-pppoe
	FixIp         string `json:"fixIp"`
	FixMask       string `json:"fixMask"`
	FixGateway    string `json:"fixGateway"`
	PrimaryDns    string `json:"primaryDns"`
	SecondaryDns  string `json:"secondaryDns"`
	PPPoEAccount  string `json:"pppoeAccount"`
	PPPoEPassword string `json:"pppoePassword"`
	DeviceRefer   int
}

//Table wan_qos
type WanQos struct {
	Port     int `json:"port"`
	Up       int `json:"up"`
	Down     int `json:"down"`
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

type Md5 struct {
	Id          int `gorm:"AUTO_INCREMENT"`
	Md5         string
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
				Mode:          "1",
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
		Sync:       0,
	}

	if DB.Find(dev).RowsAffected > 0 {
		fmt.Println("Device alread init")
		return
	}
	DB.Create(dev)
}

// Callback of Device
func (dev *Device) AfterCreate(tx *gorm.DB) (err error) {
	md5 := &Md5{
		DeviceRefer: dev.Id,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)

	return nil
}

func (dev *Device) AfterUpdate(tx *gorm.DB) (err error) {
	md5 := &Md5{
		DeviceRefer: dev.Id,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)

	return nil
}

// Callback of Qos

func (qos *Qos) AfterUpdate(tx *gorm.DB) (err error) {
	md5 := &Md5{
		DeviceRefer: qos.DeviceRefer,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)

	return nil
}

// Callback of SSID
func (ssid *Ssid) AfterUpdate(tx *gorm.DB) (err error) {
	md5 := &Md5{
		DeviceRefer: ssid.DeviceRefer,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)
	return nil
}

// Callback of WAN
func (w *Wan) AfterUpdate(tx *gorm.DB) (err error) {
	md5 := &Md5{
		DeviceRefer: w.DeviceRefer,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)

	return nil
}

// Callback of WANQOS
func (wq *WanQos) AfterUpdate(tx *gorm.DB) (err error) {

	//Get device refer form Qos
	var qos Qos
	if tx.Debug().Model(&Qos{}).Where("id=?", wq.QosRefer).Find(&qos).RowsAffected < 1 {
		fmt.Println("Not found qos")
		return nil
	}

	md5 := &Md5{
		DeviceRefer: qos.DeviceRefer,
		Md5:         md5sum(),
	}

	if tx.Debug().Model(md5).Where("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).Update(md5).RowsAffected >= 1 {
		return nil
	}
	tx.Debug().Create(md5)

	return nil
}

func GetMd5ByDeviceId(devId int) Md5 {
	var md5 Md5
	if DB.Debug().Model(&Md5{}).Where("device_refer=?", devId).Find(&md5).RowsAffected > 0 {
		//fmt.Println("Found md5")
	}
	return md5
}

func AddDevice(dev Device) error {
	if DB.Debug().Where("mac=?", dev.Mac).Find(&dev).RowsAffected > 0 {
		fmt.Println("Device alread init")
		return errors.New("Device already exist")
	}

	DB.Debug().Create(&dev)
	fmt.Println("###Id", dev.Id)

	initSsid(dev.Id) //SSID
	initWan(dev.Id)  //Wan
	AddQos(dev.Id)   //Qos, WanQos
	//Init SSID, Wan, Qos,
	return nil
}

func GetDevices() []Device {
	var devs []Device
	DB.Debug().Find(&devs)
	return devs
}
func DeleteDeviceById(id int) {
	DB.Debug().Where("id=?", id).Delete(&Device{})
}

//Wan
func GetWanByDeviceId(id int) []Wan {
	var wan []Wan
	DB.Where("device_refer=?", id).Find(&wan)
	return wan
}

func UpdateWan(wan Wan) {
	if DB.Debug().Model(&Wan{}).Where("port=? and device_refer=?", wan.Port, wan.DeviceRefer).Update(&wan).RowsAffected <= 0 {
		DB.Debug().Create(&wan)
	}
}

func DeleteWanByDeviceId(id int) {
	DB.Debug().Model(&Wan{}).Where("device_refer=?", id).Delete(&Wan{})
}

func GetWanByDeviceIdPort(wan Wan) Wan {
	var w Wan
	DB.Debug().Where("port=? and device_refer=?", wan.Port, wan.DeviceRefer).Find(&w)
	return w
}

func AddWan(wan Wan) {
	if DB.Debug().Where("device_refer=? and port=?", wan.DeviceRefer, wan.Port).Find(&wan).RowsAffected > 0 {
		fmt.Println("already insert wan params")
		return
	}
	DB.Debug().Create(&wan)
}

func initWan(id int) {
	wan := &Wan{
		Port:          1,
		Mode:          "0",
		FixIp:         "",
		FixMask:       "",
		FixGateway:    "",
		PPPoEAccount:  "",
		PPPoEPassword: "",
		PrimaryDns:    "114.114.114.114",
		SecondaryDns:  "115.115.115.115",
		DeviceRefer:   id,
	}
	if DB.Debug().Where("device_refer=?", wan.DeviceRefer).Find(wan).RowsAffected > 0 {
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
		Port:        port,
		DeviceRefer: id,
	}
	DB.Debug().Find(&ssid)
	return ssid
}

func DeleteSsidByDeviceId(id int) {
	DB.Model(&Ssid{}).Where("device_refer=?", id).Delete(&Ssid{})
}

func AddSsid(ssid []Ssid) {
	DeleteSsidByDeviceId(ssid[0].DeviceRefer)
	for _, v := range ssid {
		DB.Debug().Create(&v)
	}

}
func UpdateSsid(ssid Ssid) {
	if DB.Debug().Model(&Ssid{}).Where("port=? and device_refer=?", ssid.Port, ssid.DeviceRefer).Update(&ssid).RowsAffected < 1 {
		DB.Debug().Create(&ssid)
	}

}
func GetDeviceById(id int) Device {
	var dev Device
	DB.Debug().Where("id=?", id).Find(&dev)
	return dev
}

func GetDeviceByMac(mac string) Device {
	var dev Device
	DB.Debug().Where("mac=?", mac).Find(&dev)
	return dev
}

func UpdateDevice(dev Device) {
	DB.Debug().Model(&Device{}).Update(&dev)
}
func GetDeviceID(mac string) int {
	var device Device
	DB.Where("mac=?", mac).Find(&device)
	return device.Id
}

func initSsid(id int) {
	ssid := &Ssid{
		Port:        0,
		Name:        "hiweeds",
		Name5:       "hiweeds-5G",
		Password:    "",
		Url:         "http://hiweeds.net",
		DeviceRefer: id,
	}

	if DB.Debug().Where("device_refer=?", ssid.DeviceRefer).Find(ssid).RowsAffected > 0 {
		fmt.Println("already insert ssid")
		return
	}
	DB.Create(ssid)
}

///////////////WanQos Operation////////////////////////
func AddWanQosConfig(wanQos WanQos) {
	if DB.Debug().Where("qos_refer=? and port=?", wanQos.QosRefer, wanQos.Port).Find(&wanQos).RowsAffected > 0 {
		fmt.Println("already insert wanqos")
		return
	}
	DB.Debug().Create(&wanQos)
}
func AddWanQos(qosrefer, port int) {
	def := &WanQos{
		Port:     port,
		Up:       8192,
		Down:     102400,
		QosRefer: qosrefer,
	}
	if DB.Debug().Where("qos_refer=?", def.QosRefer).Find(def).RowsAffected > 0 {
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
	if DB.Debug().Model(&WanQos{}).Where("port=? and qos_refer=?", wanQos.Port, wanQos.QosRefer).Update(&wanQos).RowsAffected != 1 {
		DB.Debug().Create(&wanQos)
	}
}

func GetWanQosByQosId(refer int) []WanQos {
	var wanqoss []WanQos
	DB.Debug().Where("qos_refer=?", refer).Find(&wanqoss)
	return wanqoss
}

func DeleteWanQos(qosrefer, port int) {
	DB.Debug().Model(&WanQos{}).Where("qos_refer=? and port=?", qosrefer, port).Delete(&WanQos{})
}

func DeleteWanQosByQosId(qosid int) {
	DB.Debug().Model(&WanQos{}).Where("qos_refer=?", qosid).Delete(&WanQos{})
}

func InitWanQos() {

	//Add
	AddWanQos(1, 0)
	AddWanQos(1, 1)
	AddWanQos(1, 2)
	AddWanQos(1, 3)
	AddWanQos(1, 4)
}

///////////////Qos Operation//////////////////////////
func AddQos(refer int) {
	qos := &Qos{
		UpRate:      128,
		DownRate:    1536,
		TcpLimit:    200,
		UdpLimit:    100,
		DeviceRefer: refer,
	}
	if DB.Debug().Where("device_refer=?", qos.DeviceRefer).Find(qos).RowsAffected > 0 {
		fmt.Println("already insert qos")
		return
	}
	DB.Debug().Create(qos)

	fmt.Println("wan qos id:", qos.Id, "device id:", refer)
	AddWanQos(qos.Id, 0)
}
func UpdateQos(qos Qos) {
	DB.Debug().Model(&Qos{}).Where("device_refer=?", qos.DeviceRefer).Update(&qos)
}
func GetQosByDeviceId(refer int) Qos {
	var qos Qos
	DB.Debug().Find(&qos, "device_refer=?", refer)
	return qos
}
func InitQos() {
	AddQos(1)
}

func InitAllDeviceConfig() {
	InitDevice()
	InitQos()
	initSsid(1)
	initWan(1)
	InitWanQos()
}

//return []Device by PageNo and PageSize
func ListPageNoDevice(pageNo, pageSize int) []Device {
	var dev []Device

	DB.Debug().Order("id asc").Find(&dev).Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&dev)

	fmt.Println("offset len dev:", len(dev))
	for k, v := range dev {
		fmt.Println("k=", k, ", id=", v.Id)
	}

	return dev
}

//return Device number of total
func GetTotalDeviceNum() int {
	var count int
	DB.Debug().Table("devices").Count(&count)
	return count
}

func md5sum() string {
	h := md5.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
	h.Write(b)
	strCipher := h.Sum(nil)
	return hex.EncodeToString(strCipher)
}
