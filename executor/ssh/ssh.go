package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSH struct {
	Host       string
	User       string
	Port       int
	privateKey []byte
	Config     *ssh.ClientConfig
	Client     *ssh.Client
	Session    *ssh.Session
}

func getAuthMethodByPrivateKey(privateKey []byte) (ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func NewSSHByPrivateKey(host string, user string, port int, privateKey []byte) (*SSH, error) {
	s := &SSH{
		Host:       host,
		User:       user,
		Port:       port,
		privateKey: privateKey,
		Config:     nil,
		Client:     nil,
		Session:    nil,
	}
	am, err := getAuthMethodByPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	s.Config = &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{am},
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, port)
	s.Client, err = ssh.Dial("tcp", addr, s.Config)
	if err != nil {
		return nil, err
	}
	s.Session, err = s.Client.NewSession()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SSH) Destruct() {
	s.Session.Close()
	s.Client.Close()
}

func (s *SSH) SetMultipleCommand(commands []string) error {
	commands = append(commands, "exit")
	stdin, err := s.Session.StdinPipe()
	if err != nil {
		return err
	}
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
