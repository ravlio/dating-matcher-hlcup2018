package account

import "testing"
import "github.com/ravlio/highloadcup2018/gojay"

func TestJSONFullDecode(t *testing.T) {
	j := []byte(`{"birth": 811778830, 
"country": "\u0418\u0441\u043f\u0430\u0442\u0440\u0438\u0441", 
"sex": "m", 
"likes": [{"ts": 1532257894, "id": 4844}, {"ts": 1538148125, "id": 5856}, {"ts": 1484501471, "id": 6750}, {"ts": 1489241266, "id": 3332}, {"ts": 1540814553, "id": 1740}], 
"status": "\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e", 
"id": 1, 
"phone": "8(921)1234567", 
"joined": 1306886400, 
"premium": {"start": 1529225845, "finish": 1531817845}, 
"email": "achidogyr@yandex.ru", 
"fname": "\u0412\u0430\u0434\u0438\u043c", 
"interests": ["\u041f\u043b\u044f\u0436\u043d\u044b\u0439 \u043e\u0442\u0434\u044b\u0445", "\u0424\u043e\u0442\u043e\u0433\u0440\u0430\u0444\u0438\u044f", "\u0422\u0430\u043d\u0446\u0435\u0432\u0430\u043b\u044c\u043d\u0430\u044f", "\u041c\u0430\u0442\u0440\u0438\u0446\u0430", "\u0410\u043f\u0435\u043b\u044c\u0441\u0438\u043d\u043e\u0432\u044b\u0439 \u0441\u043e\u043a"], 
"sname": "\u0424\u0430\u043b\u0435\u043d\u043b\u0430\u043d", 
"city": "\u0420\u043e\u0441\u0435\u043b\u043e\u043d\u0430"}`)

	acc := &Account{}
	err := gojay.UnmarshalJSONObject(j, acc)
	if err != nil {
		t.Error(err)
	}
}
