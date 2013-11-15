package systray

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

func (p *_SystraySvr) Run() error {
	out, err := exec.Command(p.clientPath).Output()
	if err != nil {
		return err
	}
	if out != nil {
		return errors.New(string(out))
	}
	return p.serve()
}

func (p *_SystraySvr) Stop() error {
	cmd := map[string]string{
		"action": "exit",
	}
	return p.send(cmd)
}

func (p *_SystraySvr) Show(file string) error {
	cmd := map[string]string{
		"action": "show",
		"path": filepath.Join(p.iconPath, file),
	}
	return p.send(cmd)
}

func (p *_SystraySvr) OnClick(fun func()) {
	p.fclicked = fun
}

func _NewSystraySvr(iconPath string, clientPath string, port int) *_SystraySvr {
	return &_SystraySvr{iconPath, clientPath, port, make(map[net.Conn]bool), func(){}, sync.Mutex{}}
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

			for {
				n := uint32(0)
				err := binary.Read(conn, binary.LittleEndian, &n)
				if err != nil {
					println(err.Error())
					break
				}

				buf := new(bytes.Buffer)
				_, err = io.CopyN(buf, conn, int64(n))
				if err != nil {
					println(err.Error())
					break
				}

				data := buf.Bytes()
				kvs := map[string]string{}
				err = json.Unmarshal(data, &kvs)
				if err != nil {
					println(err.Error())
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

func (p *_SystraySvr) send(cmd map[string]string) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	ok := 0
	for conn, _ := range p.conns {
		_, ret := conn.Write(data)
		if ret != nil {
			err = ret
		} else {
			ok += 1
		}
	}
	if ok == 0 {
		if err != nil {
			return err
		} else {
			return errors.New("no conns")
		}
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
	fclicked func()
	lock sync.Mutex
}
