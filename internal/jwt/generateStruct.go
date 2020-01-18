package jwt

type User struct {
	Sub int64 //用户标识  唯一id
	Prv string // laravel "App\\Modes\\Sc\\StudentUser"
}
