package helper

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"auto-initail-server/utils"

	"golang.org/x/crypto/ssh"
)

type InfoSSH struct {
	pKey string
	conf *utils.Config
}

type Prosses struct {
	client *ssh.Client
	stdout *bytes.Buffer
}

func NewInfoSSH(conf *utils.Config) *InfoSSH {
	return &InfoSSH{
		conf: conf,
	}
}

func (i *InfoSSH) RunQueue(pathKey string) error {
	key, err := os.ReadFile(pathKey)
	if err != nil {
		return fmt.Errorf("unable to read public key: %v", err)
	}
	i.pKey = string(key)

	log.Println("initialing ...")
	for _, host := range *i.conf.Yamls {
		if err := i.TaskManager(host.IP, host.User, host.Password, host.Port, host.NewPort); err != nil {
			log.Printf("Error :: %s:%d :: %s", host.IP, host.Port, err.Error())
		}
	}
	log.Println("finish\nenjoy it...")

	return nil
}

func (i *InfoSSH) TaskManager(ip, username, pssword string, port, newPort int) error {

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(pssword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}
	defer client.Close()

	proc := &Prosses{
		client: client,
	}

	if err := proc.copySSHKey(i.pKey); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%s :: copy-ssh-key :: succeed", ip)
	}

	if err := proc.changeSSHPort(newPort); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: change ssh port to %d :: succeed", ip, newPort)
	}

	if err := proc.changeSSHPasswordAuthentication(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: change PasswordAuthentication to `no` :: succeed", ip)
	}

	if err := proc.addSSHAllowUsers(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: added ssh AllowUsers root debian :: succeed", ip)
	}

	if err := proc.addDefaultUser(); err != nil {
		log.Println("warning user `debian existed`")
	} else {
		log.Printf("%s :: added user debian :: succeed", ip)
	}

	if err := proc.copySSHKeyToDefaultUser(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: copy ssh key to  user debian :: succeed", ip)
	}

	if err := proc.installUFW(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: install ufw :: succeed", ip)
	}

	if err := proc.allowSSHPort(newPort); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: allow port %d :: succeed", ip, newPort)
	}

	if err := proc.enableUFW(); err != nil {
		log.Println(err)
	} else {
		log.Printf("%s :: enable ufw :: succeed", ip)
	}

	return nil
}

func (p *Prosses) installUFW() error {
	return p.RunCommand("apt update && apt install ufw -y")
}

func (p *Prosses) allowSSHPort(port int) error {
	command := fmt.Sprintf("ufw allow %d/tcp", port)
	return p.RunCommand(command)
}

func (p *Prosses) enableUFW() error {
	return p.RunCommand("ufw --force enable")
}

func (p *Prosses) copySSHKey(key string) error {
	command := fmt.Sprintf("if [ -f ~/.ssh/authorized_keys ]; then cat ~/.ssh/authorized_keys;fi")
	p.RunCommand(command)
	if strings.Contains(p.stdout.String(), key) {
		return nil
	} else {
		command := fmt.Sprintf("if [ -d ~/.ssh ]; then echo '%s' >> ~/.ssh/authorized_keys; else mkdir ~/.ssh && echo '%s' >>  ~/.ssh/authorized_keys; fi", key, key)
		return p.RunCommand(command)
	}

}

func (p *Prosses) changeSSHPort(NPort int) error {
	command := fmt.Sprintf("sed -i 's/^#\\?Port 22/Port %d/' /etc/ssh/sshd_config && systemctl restart sshd", NPort)
	return p.RunCommand(command)
}

func (p *Prosses) changeSSHPasswordAuthentication() error {
	command := fmt.Sprintf("sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config && systemctl restart sshd")
	return p.RunCommand(command)
}

func (p *Prosses) addSSHAllowUsers() error {
	p.RunCommand("cat /etc/ssh/sshd_config")
	if strings.Contains(p.stdout.String(), "AllowUsers") {
		return nil
	} else {
		command := fmt.Sprintf("sed -i 's/PermitRootLogin yes/PermitRootLogin yes\\nAllowUsers root debian/' /etc/ssh/sshd_config && systemctl restart sshd")
		return p.RunCommand(command)
	}
}

func (p *Prosses) addDefaultUser() error {
	return p.RunCommand("adduser debian")
}

func (p *Prosses) copySSHKeyToDefaultUser() error {
	return p.RunCommand("cp -r .ssh /home/debian && chown -R debian:debian /home/debian/.ssh ")
}

func (p *Prosses) RunCommand(cmd string) error {
	session, err := p.client.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("err:: %v", err)
	}
	p.stdout = &b
	return nil
}
