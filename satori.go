package satori

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
)

type Satorier interface {
	Set()
	Get()
	Put()
	Delete()
	SetVertex()
	GetVertex()
	DeleteVertex()
	SetUser()
	GetUser()
	PutUser()
	DeleteUser()
	DFS()
	GetAllWith()
	GetOneWith()
	DeleteAllWith()
	DeleteOneWith()
	PutAllWith()
	PutOneWith()
	Inject()
	Heartbeat()
}

type Satori struct {
	host     string
	port     string
	username string
	token    string
}

type FieldEntry struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Payload struct {
	Command        string       `json:"command"`
	Username       string       `json:"username"`
	Token          string       `json:"token"`
	Key            string       `json:"key"`
	Expires        bool         `json:"expires"`
	ExpirationTime int64        `json:"expiration_time"`
	ObjectType     string       `json:"type"`
	Vertices       []string     `json:"vertices"`
	Data           any          `json:"data"`
	SetUsername    string       `json:"set_username"`
	GetUsername    string       `json:"get_username"`
	DeleteUsername string       `json:"delete_username"`
	PutUsername    string       `json:"put_username"`
	ReplaceField   string       `json:"replace_field"`
	ReplaceValue   any          `json:"replace_value"`
	Role           string       `json:"role"`
	Node           string       `json:"node"`
	Relation       string       `json:"relation"`
	EncryptionKey  string       `json:"encryption_key"`
	Vertex         any          `json:"vertex"`
	FieldArray     []FieldEntry `json:"field_array"`
	Ref            string       `json:"ref"`
	Code           string       `json:"code"`
	Args           any          `json:"args"`
	Array          string       `json:"array"`
	Value          any          `json:"value"`
}

func (s *Satori) getSocket() (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (satori *Satori) Heartbeat() string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "HEARTBEAT", Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Set(key string, expires bool, expiration_time int64, objType string, vertices []string, data any) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}
	defer conn.Close()

	p := Payload{Command: "SET", Username: satori.username, Token: satori.token, Key: key, Expires: expires, ExpirationTime: expiration_time, ObjectType: objType, Vertices: vertices, Data: data}
	b, err := json.Marshal(p)
	if err != nil {
		return "Error serializing Payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	conn.Close()
	return string(res[:])
}

func (satori *Satori) Get(key string, encryptionKey string) string {
	conn, err := satori.getSocket()

	if err != nil {
		fmt.Println(err)
		return "Error on Dial"
	}

	p := Payload{Command: "GET", Username: satori.username, Token: satori.token, Key: key, EncryptionKey: encryptionKey}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)

	for {
		_, err := conn.Read(res)
		if err != nil && err != io.EOF {
			return "Error reading over socket"
		} else {
			break
		}

	}
	conn.Close()
	return string(res[:])
}

func (satori *Satori) Put(key string, replaceField string, replaceValue any, encryptionKey string) string {
	conn, err := satori.getSocket()

	if err != nil {
		fmt.Println(err)
		return "Error on Dial"
	}

	p := Payload{Command: "PUT", Username: satori.username, Token: satori.token, Key: key, ReplaceField: replaceField, ReplaceValue: replaceValue, EncryptionKey: encryptionKey}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)

	for {
		_, err := conn.Read(res)
		if err != nil && err != io.EOF {
			return "Error reading over socket"
		} else {
			break
		}

	}
	return string(res[:])
}

func (satori *Satori) Delete(key string) string {
	conn, err := satori.getSocket()

	if err != nil {
		fmt.Println(err)
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE", Username: satori.username, Token: satori.token, Key: key}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)

	for {
		_, err := conn.Read(res)
		if err != nil && err != io.EOF {
			return "Error reading over socket"
		} else {
			break
		}

	}
	return string(res[:])
}

func (satori *Satori) DFS(node string, relation string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}
	p := Payload{Command: "DFS", Node: node, Key: node, Relation: relation, Username: satori.username, Token: satori.token}
	if relation == "" {
		p = Payload{Command: "DFS", Node: node, Key: node, Username: satori.username, Token: satori.token}
	}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) SetVertex(key string, vertex []string, encryptionKey string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "SET_VERTEX", Key: key, Vertex: vertex, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) GetVertex(key string, encryptionKey string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_VERTEX", Key: key, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) DeleteVertex(key string, vertex string, encryptionKey string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_VERTEX", Key: key, Vertex: vertex, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) Encrypt(key string, encryptionKey string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "ENCRYPT", Key: key, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) Decrypt(key string, encryptionKey string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DECRYPT", Key: key, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	return str
}

func (satori *Satori) SetUser(username string, role string) string {

	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "SET_USER", SetUsername: username, Role: role, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)
	str := string(res[:])

	if str != "ERROR" {
		if username == satori.username {
			satori.token = str
		}
	}

	return str
}

func (satori *Satori) GetUser(username string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_USER", GetUsername: username, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) PutUser(username string, role string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "PUT_USER", PutUsername: username, Username: satori.username, Token: satori.token, Role: role}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteUser(username string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_USER", DeleteUsername: username, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteAuth() string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_AUTH", Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) GetAll(objType string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_ALL", ObjectType: objType, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteAll(objType string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_ALL", ObjectType: objType, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) GetAllWith(fieldArray []FieldEntry) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_ALL_WITH", FieldArray: fieldArray, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) GetOneWith(fieldArray []FieldEntry) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_ONE_WITH", FieldArray: fieldArray, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteAllWith(fieldArray []FieldEntry) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_ALL_WITH", FieldArray: fieldArray, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteOneWith(fieldArray []FieldEntry) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_ONE_WITH", FieldArray: fieldArray, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) PutAllWith(fieldArray []FieldEntry, replaceField string, replaceValue any) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "PUT_ALL_WITH", FieldArray: fieldArray, ReplaceField: replaceField, ReplaceValue: replaceValue, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) PutOneWith(fieldArray []FieldEntry, replaceField string, replaceValue any) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "PUT_ONE_WITH", FieldArray: fieldArray, ReplaceField: replaceField, ReplaceValue: replaceValue, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) SetRef(key string, ref string, encryptionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "SET_REF", Key: key, Ref: ref, EncryptionKey: encryptionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) GetRefs(key string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_REFS", Key: key, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Pop(key string, array string, encrpytionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "POP", Key: key, Array: array, EncryptionKey: encrpytionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Push(key string, array string, value any, encrpytionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "PUSH", Key: key, Array: array, EncryptionKey: encrpytionKey, Value: value, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Remove(key string, array string, value any, encrpytionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "REMOVE", Key: key, Array: array, EncryptionKey: encrpytionKey, Value: value, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Splice(key string, array string, encrpytionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "SPLICE", Key: key, Array: array, EncryptionKey: encrpytionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteRefs(key string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_REFS", Key: key, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) DeleteRef(key string, ref string, encrpytionKey string) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "DELETE_REF", Key: key, Ref: ref, EncryptionKey: encrpytionKey, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) PutAll(objectType string, replaceField string, replaceValue any) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "PUT_ALL", ObjectType: objectType, ReplaceField: replaceField, ReplaceValue: replaceValue, Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Inject(path string, args any) string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "INJECT", Code: path, Username: satori.username, Token: satori.token, Args: args}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}

func (satori *Satori) Get_Stats() string {
	conn, err := satori.getSocket()
	if err != nil {
		return "Error on Dial"
	}

	p := Payload{Command: "GET_STATS", Username: satori.username, Token: satori.token}
	b, err := json.Marshal(p)

	if err != nil {
		return "Error serializing payload"
	}

	conn.Write(b)
	res := make([]byte, 3072)
	conn.Read(res)

	return string(res[:])
}
