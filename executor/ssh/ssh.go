package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/QWERKael/utility-go/path"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type SSH struct {
	Host       string
	User       string
	Port       int
	Config     *ssh.ClientConfig
	Client     *ssh.Client
	Session    *ssh.Session
	SFTPClient *sftp.Client
}

func getAuthMethodByPrivateKey(privateKey []byte) (ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析私钥失败：%s", err.Error()))
	}
	return ssh.PublicKeys(signer), nil
}

func NewSSHByPassword(host string, user string, port int, password string) (*SSH, error) {
	am := ssh.Password(password)
	return NewSSHByAuthMethod(host, user, port, am)
}

func NewSSHByPrivateKey(host string, user string, port int, privateKey []byte) (*SSH, error) {
	am, err := getAuthMethodByPrivateKey(privateKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("根据私钥创建AuthMethod失败：%s", err.Error()))
	}
	return NewSSHByAuthMethod(host, user, port, am)
}

func NewSSHByAuthMethod(host string, user string, port int, am ssh.AuthMethod) (*SSH, error) {
	s := &SSH{
		Host:       host,
		User:       user,
		Port:       port,
		Config:     nil,
		Client:     nil,
		Session:    nil,
		SFTPClient: nil,
	}
	s.Config = &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{am},
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, port)
	var err error
	s.Client, err = ssh.Dial("tcp", addr, s.Config)
	if err != nil {
		return nil, err
	}
	s.Session, err = s.Client.NewSession()
	if err != nil {
		return nil, err
	}
	s.SFTPClient, err = sftp.NewClient(s.Client)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SSH) Destruct() {
	_ = s.Session.Close()
	_ = s.SFTPClient.Close()
	_ = s.Client.Close()
}

//func (s *SSH) SetMultipleCommand(commands []string) error {
//	commands = append(commands, "exit")
//	stdin, err := s.Session.StdinPipe()
//	if err != nil {
//		return err
//	}
//	for _, cmd := range commands {
//		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (s *SSH) CheckRemotePath(remotePath string) (path.FileType, error) {
	if remotePath == "" {
		return path.Unknown, errors.New("path is nil")
	}
	fi, err := s.SFTPClient.Lstat(remotePath)
	if err != nil {
		if os.IsNotExist(err) {
			return path.NotExist, nil
		} else {
			return path.Unknown, err
		}
	}
	if fi.IsDir() {
		return path.Dir, nil
	} else {
		return path.File, nil
	}
}

func (s *SSH) WriteRemoteFile(remotePath string, data []byte) error {
	if ft, _ := s.CheckRemotePath(remotePath); ft != path.NotExist {
		return errors.New("指定的文件路径不为空")
	}
	f, err := s.SFTPClient.Create(remotePath)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}

func (s *SSH) ReadRemoteFile(remotePath string) ([]byte, error) {
	if ft, _ := s.CheckRemotePath(remotePath); ft != path.File {
		return nil, errors.New("指定的文件路径不为文件")
	}
	f, err := s.SFTPClient.OpenFile(remotePath, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	var data bytes.Buffer
	_, err = f.WriteTo(&data)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	return data.Bytes(), nil
}
