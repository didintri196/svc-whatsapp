package requests

type MDevicesRequest struct {
	MUserId string `gorm:"column:M_User_id"`
	Jid     string `gorm:"column:jid"`
	Server  string `gorm:"column:server"`
	Phone   string `gorm:"column:phone"`
}
