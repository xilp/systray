package systray

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

func (p *_SystraySvr) Run() error {
	_, err := os.Stat(p.clientPath)
	if err != nil {
		return err
	}

	go func() {
		if len(p.clientPath) == 0 {
			return
		}
		_, err := exec.Command(p.clientPath, "-xilp_systray_port", strconv.Itoa(p.port)).Output()
		if err != nil {
			panic(err)
		}
	}()
	return p.serve()
}

func (p *_SystraySvr) Stop() error {
	cmd := map[string]string{"action": "exit"}
	return p.send(cmd)
}

func (p *_SystraySvr) Show(file string, hint string) error {
	path, err := filepath.Abs(filepath.Join(p.iconPath, file))
	if err != nil {
		return err
	}
	cmd := map[string]string{"action": "show", "path": path, "hint": hint}
	return p.send(cmd)
}

func (p *_SystraySvr) OnClick(fun func()) {
	p.fclicked = fun
}

func _NewSystraySvr(iconPath string, clientPath string, port int) *_SystraySvr {
	return &_SystraySvr{iconPath, clientPath, port, make(map[net.Conn]bool), nil, func(){}, sync.Mutex{}}
}

func (p *_SystraySvr) serve() error {
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(p.port))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			p.lock.Lock()
			p.conns[conn] = true
			p.lock.Unlock()

			p.resend(conn)

			for {
				n := uint32(0)
				err := binary.Read(conn, binary.LittleEndian, &n)
				if err != nil {
					break
				}

				buf := new(bytes.Buffer)
				_, err = io.CopyN(buf, conn, int64(n))
				if err != nil {
					break
				}

				data := buf.Bytes()
				kvs := map[string]string{}
				err = json.Unmarshal(data, &kvs)
				if err != nil {
					continue
				}
				p.received(kvs)
			}

			p.lock.Lock()
			delete(p.conns, conn)
			p.lock.Unlock()
			conn.Close()
		}(conn)
	}
}

func (p *_SystraySvr) resend(conn net.Conn) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	_, err := conn.Write(p.lastest)
	return err
}

func (p *_SystraySvr) send(cmd map[string]string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	data = buf.Bytes()
	p.lastest = data

	ok := 0
	for conn, _ := range p.conns {
		_, ret := conn.Write(data)
		if ret != nil {
			err = ret
		} else {
			ok += 1
		}
	}
	if ok == 0 && err != nil {
		return err
	}
	return nil
}

func (p *_SystraySvr) received(cmd map[string]string) {
	action := cmd["action"]
	if len(action) == 0 {
		return
	}
	switch action {
	case "clicked":
		p.fclicked()
	}
}

type _SystraySvr struct {
	iconPath string
	clientPath string
	port int
	conns map[net.Conn]bool
	lastest []byte
	fclicked func()
	lock sync.Mutex
}
