package bingo

import (
	"github.com/boltdb/bolt"
	"errors"
	"bytes"
	"encoding/binary"
)

// 使用Bolt 持久化存储键值对
// 封装几个函数，对kv数据进行增删改查
// buckets


// 一个bolt的数据库连接,同样，一次只允许存在一个Bolt连接
var Bolt *bolt.DB

// 打开一个数据库，返回链接
func OpenBolt(boltName ... string) *bolt.DB {
	if len(boltName) > 1 { // 传入了多个参数,返回空指针
		return nil
	}
	if len(boltName) == 0 { // 没传参数
		if Bolt != nil {
			return Bolt
		}
		db, _ := bolt.Open(Env.Get("KVSTORAGE_DB_NAME"), 0600, nil) // 根据配置文件，创建db
		return db
	} else {
		db, _ := bolt.Open(boltName[0], 0600, nil)
		return db
	}
}

// 向数据库中设置键值对，只能对env文件中指定的数据文件，指定的buckets进行操作
// 最多可以传3个参数，第一个是键，第二个是值，第三个是bucket名,只允许 int 或者 string类型的值
func KVStorageSet(key string, value ... string) error {
    // 只支持string
	switch len(value) {
	case 1:
		// 默认buckets
		BoltInit()
		Bolt.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(Env.Get("KVSTORAGE_BUCKET")))
			if b == nil { // 不存在这个bucket，新建
				b, err := tx.CreateBucket([]byte(Env.Get("KVSTORAGE_BUCKET")))
				//serializeData,err := Serialize(value[0])
			    err = b.Put([]byte(key),[]byte(value[0]))
			    return err
			}else{   // 存在，直接存入数据
				//serializeData,err := Serialize(value[0])
				err := b.Put([]byte(key),[]byte(value[0]))
				return err
			}
		})
		break
	case 2:
		// 新建buckets
		BoltInit()
		Bolt.Update(func(tx *bolt.Tx) error {
			// 第三个参数必须是string
			b,err := tx.CreateBucketIfNotExists([]byte(value[1]))
			//serializeData,err := Serialize(value[0])
			err = b.Put([]byte(key),[]byte(value[0]))
			return err
		})
		break
	default:
		return errors.New("the KVStorageSet function need 2 or 3 arguments ")
	}
	// 把value转为bytes数组
	return nil
}

// 第一个是key，第二个是bucket名字，默认是
func KVStorageGet(keyAndBucket ... string) (interface{},error)  {
	if len(keyAndBucket)!=1 && len(keyAndBucket)!=2 {
		return nil,errors.New("the KVStorageGet function need 1 or 2 arguments")
	}
	BoltInit()
	var err error
	var res interface{}
	if len(keyAndBucket)==1 {  // 从默认的bucket中取数据
		Bolt.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(Env.Get("KVSTORAGE_BUCKET")))
			if b==nil{
				err = errors.New("the bucket "+Env.Get("KVSTORAGE_BUCKET")+" is not exists")
			}else{
				resByte := b.Get([]byte(keyAndBucket[0]))
				res = string(resByte[:])
			}
			return err
		})
	}
	return res,err   // 返回结果
}

// 初始化Bolt
func BoltInit()  {
	if Bolt==nil{
		Bolt = OpenBolt()
	}
}

// 整形转字节
func intToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}