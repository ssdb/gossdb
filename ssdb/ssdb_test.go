package ssdb

import "testing"
	
var ip = "ssdb"
var port = 16379

func TestConnect(t *testing.T) {

	db, err := Connect("zup", port)
	if err == nil {
        t.Error("connect to bad host did not return err")
	}
	if db != nil {
        t.Error("connect to bad host returned non nil db")
	}

	db, err = Connect("ssdb", 0)
	if err == nil {
        t.Error("connect to bad port did not return err")
	}
	if db != nil {
        t.Error("connect to bad port returned non nil db")
	}

	db, err = Connect(ip, port)
	if err != nil {
        t.Error("Failed to connect")
	}
	defer db.Close()

	db, err = Connect(ip, port)
	if err != nil {
        t.Error("Close:second connect raised an error:%v:", err)
	}

}


func TestClose(t *testing.T) {

	db, err := Connect(ip, port)
	if err != nil {
        t.Error("Close:connect returned an err:%v:", err)
	}
	if db == nil {
        t.Error("Close:connect returned a nil db")
	}

	db.Close()
	if err != nil {
        t.Error("Close:returned an error:%v:", err)
	}

	val, err := db.Set("a", "xxx")
	if val == true {
        t.Error("Close:Set:val returned true after db closed")
	}	
	if val != nil {
        t.Error("Close:Set:val returned non-nil after db closed")
	}	
	if err == nil {
        t.Error("Close:Set returned no error after db closed")
	}

}

func TestInitPool(t *testing.T) {

	cpm, err := InitPool(ip, port, 11)
	if err != nil {
		t.Error("failed to init pool")
	}
	if cpm == nil {
		t.Error("failed to init pool")
	}

}

func TestSet(t *testing.T) {

	db, err := Connect(ip, port)
	if err != nil {
        t.Error("Failed to connect")
	}
	//defer db.Close()

	val, err := db.Set("a", "xxx")
	if val != true {
        t.Error("Set val returned false")
	}	
	if err != nil {
        t.Error("Set err returned not nil err")
	}

	// add negative testts
	db.Close()

	val, err = db.Set("a", "xxx")
	if val == true {
        t.Error("Set val on closed db returned true")
	}	
	if err == nil {
        t.Error("Set err on closed db returned not nil err")
	}


}

func TestGet(t *testing.T) {

	db, err := Connect(ip, port)
	if err != nil {
        t.Error("Failed to connect")
	}
	defer db.Close()

	val, err := db.Set("a", "xxxdodoehodh eodhoe dohe j")
	if val != true {
        t.Error("Set val returned false")
	}	
	if err != nil {
        t.Error("Set err returned not nil err")
	}

	val, err = db.Get("a")
	if val == nil {
        t.Error("Get returned nil")
	}	
	if err != nil {
        t.Error("Get returned err")
	}
	if val != "xxxdodoehodh eodhoe dohe j" {
        t.Error("Get did not return the right value")
	}
	if val == "xxx" {
        t.Error("Get did not return return xxx")
	}	

}

func TestDel(t *testing.T) {

	db, err := Connect(ip, port)
	if err != nil {
        t.Error("Failed to connect")
	}
	defer db.Close()

	val, err := db.Set("a", "xxx")
	if val != true {
        t.Error("Set val returned false")
	}	
	if err != nil {
        t.Error("Set err returned not nil err")
	}

	val, err = db.Get("a")
	if val != "xxx" {
        t.Error("Get did not return xxx")
	}
	if val == nil {
        t.Error("Get returned nil")
	}	
	if err != nil {
        t.Error("Get returned err")
	}

	val, err = db.Del("a")
	if val != true {
        t.Error("Del returned false")
	}	
	if err != nil {
        t.Error("Del returned err")
	}
	
	val, err = db.Get("a")
	if val == "xxx" {
        t.Error("Get returned xxx after Del")
	}
	if val != nil {
        t.Error("Get returned non-nil")
	}	
	if err != nil {
        t.Error("Get returned err")
	}

	val, err = db.Del("a")
	if val != true {
        t.Error("Del returned non-nil:%v:", val)
	}	
	if err != nil {
        t.Error("Get returned err")
	}
}